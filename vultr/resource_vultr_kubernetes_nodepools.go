package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
			State: func(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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
		Schema: nodePoolSchema(true),
	}
}

func resourceVultrKubernetesNodePoolsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)

	req := &govultr.NodePoolReq{
		NodeQuantity: d.Get("node_quantity").(int),
		Label:        d.Get("label").(string),
		Plan:         d.Get("plan").(string),
		Tag:          d.Get("tag").(string),
		AutoScaler:   govultr.BoolToBoolPtr(d.Get("auto_scaler").(bool)),
		MinNodes:     d.Get("min_nodes").(int),
		MaxNodes:     d.Get("max_nodes").(int),
		UserData:     d.Get("user_data").(string),
	}

	if labelsVal, labelsOK := d.GetOk("labels"); labelsOK {
		labels := make(map[string]string)
		for k, v := range labelsVal.(map[string]interface{}) {
			labels[k] = v.(string)
		}

		req.Labels = labels
	}

	if taintsVal, taintsOK := d.GetOk("taints"); taintsOK {
		var taints []govultr.Taint
		taintVals := taintsVal.(*schema.Set).List()
		for i := range taintVals {
			taint := taintVals[i].(map[string]interface{})
			taints = append(taints, govultr.Taint{
				Key:    taint["key"].(string),
				Value:  taint["value"].(string),
				Effect: taint["effect"].(string),
			})
		}

		req.Taints = taints
	}

	nodePool, _, err := client.Kubernetes.CreateNodePool(ctx, clusterID, req)
	if err != nil {
		return diag.Errorf("error creating node pool: %v", err)
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

	return resourceVultrKubernetesNodePoolsRead(ctx, d, meta)
}

func resourceVultrKubernetesNodePoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)

	nodePool, _, err := client.Kubernetes.GetNodePool(ctx, clusterID, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Unauthorized") {
			return diag.Errorf("API authorization error: %v", err)
		}
		if strings.Contains(err.Error(), "Invalid NodePool ID") {
			log.Printf("[WARN] Kubernetes NodePool (%v) not found", d.Id())
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
	if err := d.Set("labels", nodePool.Labels); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `labels` read value: %v", err)
	}

	var taints []map[string]interface{}
	for i := range nodePool.Taints {
		taints = append(taints, map[string]interface{}{
			"key":    nodePool.Taints[i].Key,
			"value":  nodePool.Taints[i].Value,
			"effect": nodePool.Taints[i].Effect,
		})
	}

	if err := d.Set("taints", taints); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `taints` read value: %v", err)
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

	return nil
}

func resourceVultrKubernetesNodePoolsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)
	nodePoolID := d.Id()

	// Handle labels changes using dedicated label endpoints
	if d.HasChange("labels") {
		if err := syncNodePoolLabels(ctx, client, clusterID, nodePoolID, d); err != nil {
			return diag.Errorf("error syncing VKE node pool labels %v : %v", nodePoolID, err)
		}
	}

	// Handle taints changes using dedicated taint endpoints
	if d.HasChange("taints") {
		if err := syncNodePoolTaints(ctx, client, clusterID, nodePoolID, d); err != nil {
			return diag.Errorf("error syncing VKE node pool taints %v : %v", nodePoolID, err)
		}
	}

	// Handle other node pool updates (quantity, tag, autoscaler, user_data)
	if d.HasChange("node_quantity") || d.HasChange("tag") || d.HasChange("auto_scaler") ||
		d.HasChange("min_nodes") || d.HasChange("max_nodes") || d.HasChange("user_data") {

		req := &govultr.NodePoolReqUpdate{
			NodeQuantity: d.Get("node_quantity").(int),
			Tag:          govultr.StringToStringPtr(d.Get("tag").(string)),
			AutoScaler:   govultr.BoolToBoolPtr(d.Get("auto_scaler").(bool)),
			MinNodes:     d.Get("min_nodes").(int),
			MaxNodes:     d.Get("max_nodes").(int),
		}

		if d.HasChange("user_data") {
			req.UserData = govultr.StringToStringPtr(d.Get("user_data").(string))
		}

		if _, _, err := client.Kubernetes.UpdateNodePool(ctx, clusterID, nodePoolID, req); err != nil {
			return diag.Errorf("error updating VKE node pool %v : %v", nodePoolID, err)
		}
	}

	return resourceVultrKubernetesNodePoolsRead(ctx, d, meta)
}

// syncNodePoolLabels synchronizes node pool labels using the dedicated label endpoints
func syncNodePoolLabels(ctx context.Context, client *govultr.Client, clusterID, nodePoolID string, d *schema.ResourceData) error {
	// Get current labels from API
	currentLabels, _, err := client.Kubernetes.ListNodePoolLabels(ctx, clusterID, nodePoolID)
	if err != nil {
		return fmt.Errorf("error listing current labels: %v", err)
	}

	// Build map of current labels by key for easy lookup
	currentLabelsByKey := make(map[string]govultr.NodePoolLabel)
	for _, label := range currentLabels {
		currentLabelsByKey[label.Key] = label
	}

	// Get desired labels from Terraform state
	desiredLabels := make(map[string]string)
	if labelsVal, ok := d.GetOk("labels"); ok {
		for k, v := range labelsVal.(map[string]interface{}) {
			desiredLabels[k] = v.(string)
		}
	}

	// Delete labels that are no longer desired or have changed values
	for key, currentLabel := range currentLabelsByKey {
		desiredValue, exists := desiredLabels[key]
		if !exists || desiredValue != currentLabel.Value {
			log.Printf("[INFO] Deleting label %s (ID: %s) from node pool %s", key, currentLabel.ID, nodePoolID)
			if err := client.Kubernetes.DeleteNodePoolLabel(ctx, clusterID, nodePoolID, currentLabel.ID); err != nil {
				return fmt.Errorf("error deleting label %s: %v", key, err)
			}
		}
	}

	// Create labels that are new or have changed values
	for key, desiredValue := range desiredLabels {
		currentLabel, exists := currentLabelsByKey[key]
		if !exists || desiredValue != currentLabel.Value {
			log.Printf("[INFO] Creating label %s=%s on node pool %s", key, desiredValue, nodePoolID)
			req := &govultr.NodePoolLabelReq{
				Key:   key,
				Value: desiredValue,
			}
			if _, _, err := client.Kubernetes.CreateNodePoolLabel(ctx, clusterID, nodePoolID, req); err != nil {
				return fmt.Errorf("error creating label %s: %v", key, err)
			}
		}
	}

	return nil
}

// syncNodePoolTaints synchronizes node pool taints using the dedicated taint endpoints
func syncNodePoolTaints(ctx context.Context, client *govultr.Client, clusterID, nodePoolID string, d *schema.ResourceData) error {
	// Get current taints from API
	currentTaints, _, err := client.Kubernetes.ListNodePoolTaints(ctx, clusterID, nodePoolID)
	if err != nil {
		return fmt.Errorf("error listing current taints: %v", err)
	}

	// Build map of current taints by key+effect for easy lookup
	// Using key+effect as the identifier since a node can have multiple taints with same key but different effects
	currentTaintsByKeyEffect := make(map[string]govultr.NodePoolTaint)
	for _, taint := range currentTaints {
		mapKey := taint.Key + ":" + taint.Effect
		currentTaintsByKeyEffect[mapKey] = taint
	}

	// Get desired taints from Terraform state
	type desiredTaint struct {
		Key    string
		Value  string
		Effect string
	}
	desiredTaints := make(map[string]desiredTaint)
	if taintsVal, ok := d.GetOk("taints"); ok {
		taintVals := taintsVal.(*schema.Set).List()
		for _, t := range taintVals {
			taint := t.(map[string]interface{})
			key := taint["key"].(string)
			effect := taint["effect"].(string)
			mapKey := key + ":" + effect
			desiredTaints[mapKey] = desiredTaint{
				Key:    key,
				Value:  taint["value"].(string),
				Effect: effect,
			}
		}
	}

	// Delete taints that are no longer desired or have changed values
	for mapKey, currentTaint := range currentTaintsByKeyEffect {
		desired, exists := desiredTaints[mapKey]
		if !exists || desired.Value != currentTaint.Value {
			log.Printf("[INFO] Deleting taint %s (ID: %s) from node pool %s", mapKey, currentTaint.ID, nodePoolID)
			if err := client.Kubernetes.DeleteNodePoolTaint(ctx, clusterID, nodePoolID, currentTaint.ID); err != nil {
				return fmt.Errorf("error deleting taint %s: %v", mapKey, err)
			}
		}
	}

	// Create taints that are new or have changed values
	for mapKey, desired := range desiredTaints {
		currentTaint, exists := currentTaintsByKeyEffect[mapKey]
		if !exists || desired.Value != currentTaint.Value {
			log.Printf("[INFO] Creating taint %s=%s:%s on node pool %s", desired.Key, desired.Value, desired.Effect, nodePoolID)
			req := &govultr.NodePoolTaintReq{
				Key:    desired.Key,
				Value:  desired.Value,
				Effect: desired.Effect,
			}
			if _, _, err := client.Kubernetes.CreateNodePoolTaint(ctx, clusterID, nodePoolID, req); err != nil {
				return fmt.Errorf("error creating taint %s: %v", mapKey, err)
			}
		}
	}

	return nil
}

func resourceVultrKubernetesNodePoolsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	if err := client.Kubernetes.DeleteNodePool(ctx, d.Get("cluster_id").(string), d.Id()); err != nil {
		return diag.Errorf("error deleting VKE nodepool %v : %v", d.Id(), err)
	}

	return nil
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
