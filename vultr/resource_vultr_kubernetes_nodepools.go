package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrKubernetesNodePools() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrKubernetesNodePoolsCreate,
		ReadContext:   resourceVultrKubernetesNodePoolsRead,
		UpdateContext: resourceVultrKubernetesNodePoolsUpdate,
		DeleteContext: resourceVultrKubernetesNodePoolsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
				ids := strings.SplitN(d.Id(), " ", 2)
				if len(ids) != 2 || ids[0] == "" || ids[1] == "" {
					err := fmt.Errorf("unexpected format of node pool import IDs (%s): expected 'clusterID nodePoolID'", d.Id())
					return nil, err
				}

				d.SetId(ids[1])
				if err := d.Set("cluster_id", ids[0]); err != nil {
					return nil, fmt.Errorf("unable to set cluster ID for import state function")
				}

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema:        resourceVultrKubernetesNodePoolsV1(true),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceVultrKubernetesNodePoolsV0(true).CoreConfigSchema().ImpliedType(),
				Upgrade: resourceVultrKubernetesNodePoolsStateUpgradeV0ToV1,
				Version: 0,
			},
		},
	}
}

func resourceVultrKubernetesNodePoolsStateUpgradeV0ToV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) { //nolint:lll
	if len(rawState) == 0 {
		return rawState, nil
	}

	client := meta.(*Client).govultrClient()

	stateLabels := map[string]interface{}{}
	switch rawState["labels"].(type) {
	case map[string]interface{}:
		stateLabels = rawState["labels"].(map[string]interface{})
	}

	stateTaints := []interface{}{}
	switch rawState["taints"].(type) {
	case []interface{}:
		stateTaints = rawState["taints"].([]interface{})
	}

	log.Println("[INFO] migrating kubernetes node pool labels from v0 to v1")
	refLabels, _, err := client.Kubernetes.ListNodePoolLabels(
		ctx,
		rawState["cluster_id"].(string),
		rawState["id"].(string),
	)
	if err != nil {
		log.Println("[ERROR] unable to retrieve updated kubernetes node pool labels from client")
		return rawState, err
	}

	newStateLabels := []map[string]interface{}{}
	for j := range refLabels {
		for key, val := range stateLabels {
			newStateLabel := map[string]interface{}{}
			if key == refLabels[j].Key {
				newStateLabel["key"] = key
				newStateLabel["value"] = val.(string)
				newStateLabel["id"] = refLabels[j].ID

				newStateLabels = append(newStateLabels, newStateLabel)
			}
		}
	}

	// replace the labels state
	delete(rawState, "labels")
	rawState["labels"] = newStateLabels

	log.Println("[INFO] migrating kubernetes node pool taints from v0 to v1")
	refTaints, _, err := client.Kubernetes.ListNodePoolTaints(
		ctx,
		rawState["cluster_id"].(string),
		rawState["id"].(string),
	)
	if err != nil {
		log.Println("[ERROR] unable to retrieve updated kubernetes node pool taints from client")
		return rawState, err
	}

	for j := range refTaints {
		for k := range stateTaints {
			stateTaintData := stateTaints[k].(map[string]interface{})
			if stateTaintData["key"] == refTaints[j].Key {
				stateTaintData["id"] = refTaints[j].ID
			}
		}
	}

	return rawState, nil
}

func resourceVultrKubernetesNodePoolsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)

	nodeQuantity, err := nodePoolCreateQuantity(
		d.Get("node_quantity").(int),
		d.Get("auto_scaler").(bool),
		d.Get("min_nodes").(int),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &govultr.NodePoolReq{
		NodeQuantity: nodeQuantity,
		Label:        d.Get("label").(string),
		Plan:         d.Get("plan").(string),
		Tag:          d.Get("tag").(string),
		AutoScaler:   govultr.BoolToBoolPtr(d.Get("auto_scaler").(bool)),
		MinNodes:     d.Get("min_nodes").(int),
		MaxNodes:     d.Get("max_nodes").(int),
		UserData:     d.Get("user_data").(string),
	}

	nodePool, _, err := client.Kubernetes.CreateNodePool(ctx, clusterID, req)
	if err != nil {
		return diag.Errorf(
			"error creating node pool: %v. "+
				"node_quantity must be a positive integer on create (initial size), even when auto_scaler is true",
			err,
		)
	}

	d.SetId(nodePool.ID)
	if err := d.Set("cluster_id", clusterID); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `cluster_id` create value: %v", err)
	}

	//block until status is ready
	if _, err = waitForNodePoolAvailable(ctx, d, "active", []string{"pending"}, "status", meta); err != nil {
		return diag.Errorf(
			"error while waiting for node pool %v to be completed: %v", d.Id(), err)
	}

	labels := d.Get("labels").(*schema.Set).List()
	for i := range labels {
		labelMap := labels[i].(map[string]interface{})
		_, _, err := client.Kubernetes.CreateNodePoolLabel(
			ctx,
			clusterID,
			d.Id(),
			&govultr.NodePoolLabelReq{
				Key:   labelMap["key"].(string),
				Value: labelMap["value"].(string),
			},
		)

		if err != nil {
			return diag.Errorf(
				"error creating label %q for node pool %v on cluster %v : %v",
				labelMap["key"].(string),
				d.Id(),
				clusterID,
				err,
			)
		}
	}

	taints := d.Get("taints").(*schema.Set).List()
	for i := range taints {
		taintMap := taints[i].(map[string]interface{})
		_, _, err := client.Kubernetes.CreateNodePoolTaint(
			ctx,
			clusterID,
			d.Id(),
			&govultr.NodePoolTaintReq{
				Key:    taintMap["key"].(string),
				Value:  taintMap["value"].(string),
				Effect: taintMap["effect"].(string),
			},
		)

		if err != nil {
			return diag.Errorf(
				"error creating taint %q for node pool %v on cluster %v : %v",
				taintMap["key"].(string),
				d.Id(),
				clusterID,
				err,
			)
		}
	}

	return resourceVultrKubernetesNodePoolsRead(ctx, d, meta)
}

func resourceVultrKubernetesNodePoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)

	nodePool, _, err := client.Kubernetes.GetNodePool(ctx, clusterID, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Unauthorized") {
			return diag.Errorf("api authorization error: %v", err)
		}

		missing, missErr := checkIsMissing(err, "Invalid NodePool ID")
		if missErr != nil {
			return diag.Errorf("error in api response %q : %v", err, missErr)
		}

		if missing {
			tflog.Warn(ctx, fmt.Sprintf("removing kubernetes node pool (%v) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting node pool: %v", err)
	}

	if err := d.Set("status", nodePool.Status); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `status` read value: %v", err)
	}
	if err := d.Set("label", nodePool.Label); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `label` read value: %v", err)
	}
	if err := d.Set("plan", nodePool.Plan); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `plan` read value: %v", err)
	}
	if err := d.Set("tag", nodePool.Tag); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `tag` read value: %v", err)
	}
	if err := d.Set("node_quantity", nodePool.NodeQuantity); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `node_quantity` read value: %v", err)
	}
	if err := d.Set("date_created", nodePool.DateCreated); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `date_created` read value: %v", err)
	}
	if err := d.Set("date_updated", nodePool.DateUpdated); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `date_updated` read value: %v", err)
	}
	if err := d.Set("auto_scaler", nodePool.AutoScaler); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `auto_scaler` read value: %v", err)
	}
	if err := d.Set("min_nodes", nodePool.MinNodes); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `min_nodes` read value: %v", err)
	}
	if err := d.Set("max_nodes", nodePool.MaxNodes); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `max_nodes` read value: %v", err)
	}

	var instances []map[string]interface{}
	for _, v := range nodePool.Nodes {
		n := map[string]interface{}{
			"id":           v.ID,
			"status":       v.Status,
			"date_created": v.DateCreated,
			"label":        v.Label,
		}
		instances = append(instances, n)
	}

	if err := d.Set("nodes", instances); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `nodes` read value: %v", err)
	}

	// add node pool labels
	labelData, _, err := client.Kubernetes.ListNodePoolLabels(ctx, clusterID, d.Id())
	if err != nil {
		return diag.Errorf("error getting cluster (%v) node pool (%v) labels : %v", clusterID, d.Id(), err)
	}

	var labels []map[string]interface{}
	for j := range labelData {
		labels = append(labels, map[string]interface{}{
			"id":    labelData[j].ID,
			"key":   labelData[j].Key,
			"value": labelData[j].Value,
		})
	}

	if err := d.Set("labels", labels); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `labels` read value: %v", err)
	}

	// add node pool taints
	taintData, _, err := client.Kubernetes.ListNodePoolTaints(ctx, clusterID, d.Id())
	if err != nil {
		return diag.Errorf("error getting cluster (%v) node pool (%v) taints : %v", clusterID, d.Id(), err)
	}

	var taints []map[string]interface{}
	for j := range taintData {
		taints = append(taints, map[string]interface{}{
			"id":     taintData[j].ID,
			"key":    taintData[j].Key,
			"value":  taintData[j].Value,
			"effect": taintData[j].Effect,
		})
	}

	if err := d.Set("taints", taints); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `taints` read value: %v", err)
	}

	return nil
}

func resourceVultrKubernetesNodePoolsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)

	req := &govultr.NodePoolReqUpdate{
		NodeQuantity: d.Get("node_quantity").(int),
		Tag:          govultr.StringToStringPtr(d.Get("tag").(string)),
	}

	if d.HasChange("auto_scaler") {
		req.AutoScaler = govultr.BoolToBoolPtr(d.Get("auto_scaler").(bool))
	}

	if d.HasChange("min_nodes") {
		req.MinNodes = d.Get("min_nodes").(int)
	}

	if d.HasChange("max_nodes") {
		req.MaxNodes = d.Get("max_nodes").(int)
	}

	if d.HasChange("user_data") {
		req.UserData = govultr.StringToStringPtr(d.Get("user_data").(string))
	}

	if _, _, err := client.Kubernetes.UpdateNodePool(ctx, clusterID, d.Id(), req); err != nil {
		return diag.Errorf("error updating VKE node pool %v : %v", d.Id(), err)
	}

	if d.HasChange("labels") {
		oldLabels, newLabels := d.GetChange("labels")
		oldLabelsData := oldLabels.(*schema.Set).List()
		newLabelsData := newLabels.(*schema.Set).List()

		err := updateNodePoolOptions(ctx, client, clusterID, d.Id(), "labels", oldLabelsData, newLabelsData)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("taints") {
		oldTaints, newTaints := d.GetChange("taints")
		oldTaintsData := oldTaints.(*schema.Set).List()
		newTaintsData := newTaints.(*schema.Set).List()

		err := updateNodePoolOptions(ctx, client, clusterID, d.Id(), "taints", oldTaintsData, newTaintsData)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVultrKubernetesNodePoolsRead(ctx, d, meta)
}

func resourceVultrKubernetesNodePoolsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	if err := client.Kubernetes.DeleteNodePool(ctx, d.Get("cluster_id").(string), d.Id()); err != nil {
		return diag.Errorf("error deleting VKE nodepool %v : %v", d.Id(), err)
	}

	return nil
}

// nodePoolCreateQuantity returns the initial node count for CreateNodePool.
// The Vultr API always requires a positive node_quantity on create. When auto_scaler
// is enabled, that value is the initial size and min_nodes/max_nodes bound the scaler.
// If node_quantity is 0 (for example lifecycle.ignore_changes on create) and
// auto_scaler is true, min_nodes is used when it is at least 1.
func nodePoolCreateQuantity(nodeQuantity int, autoScaler bool, minNodes int) (int, error) {
	if nodeQuantity >= 1 {
		return nodeQuantity, nil
	}

	if autoScaler && minNodes >= 1 {
		log.Printf(
			"[WARN] node_quantity was %d on create; using min_nodes (%d) as the initial node count "+
				"because auto_scaler is enabled",
			nodeQuantity,
			minNodes,
		)
		return minNodes, nil
	}

	return 0, fmt.Errorf(
		"node_quantity must be a positive integer on create (it is the initial node count). " +
			"When auto_scaler is true, set node_quantity to the initial size (usually equal to min_nodes). " +
			"lifecycle.ignore_changes = [node_quantity] can send 0 on create and is unsafe unless min_nodes >= 1",
	)
}

func waitForNodePoolAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll
	log.Printf(
		"[INFO] Waiting for node pool (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newNodePoolStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     5 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newNodePoolStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) retry.StateRefreshFunc { //nolint:lll
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Creating node pool")

		np, _, err := client.Kubernetes.GetNodePool(ctx, d.Get("cluster_id").(string), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving node pool %s ", d.Id())
		}

		if attr == "status" {
			log.Printf("[INFO] The node pool status is %v", np.Status)
			return np, np.Status, nil
		}

		return nil, "", nil
	}
}
