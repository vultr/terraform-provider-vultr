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

func resourceVultrKubernetesNodePoolsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	}

	// Handle Labels if provided
	if labels, ok := d.GetOk("labels"); ok {
		labelMap := make(map[string]string)
		for k, v := range labels.(map[string]interface{}) {
			labelMap[k] = v.(string)
		}
		req.Labels = labelMap
	}

	// Handle Taints if provided
	if taints, ok := d.GetOk("taints"); ok {
		taintsList := taints.([]interface{})
		reqTaints := make([]govultr.Taint, 0, len(taintsList))

		for _, t := range taintsList {
			taintMap := t.(map[string]interface{})
			reqTaints = append(reqTaints, govultr.Taint{
				Key:    taintMap["key"].(string),
				Value:  taintMap["value"].(string),
				Effect: taintMap["effect"].(string),
			})
		}
		req.Taints = reqTaints
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

func resourceVultrKubernetesNodePoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Set labels
	if err := d.Set("labels", nodePool.Labels); err != nil {
		return diag.Errorf("unable to set resource kubernetes_nodepools `labels` read value: %v", err)
	}

	// Set taints
	if len(nodePool.Taints) > 0 {
		taints := make([]map[string]interface{}, len(nodePool.Taints))
		for i, taint := range nodePool.Taints {
			taints[i] = map[string]interface{}{
				"key":    taint.Key,
				"value":  taint.Value,
				"effect": taint.Effect,
			}
		}
		if err := d.Set("taints", taints); err != nil {
			return diag.Errorf("unable to set resource kubernetes_nodepools `taints` read value: %v", err)
		}
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

func resourceVultrKubernetesNodePoolsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	clusterID := d.Get("cluster_id").(string)

	req := &govultr.NodePoolReqUpdate{
		NodeQuantity: d.Get("node_quantity").(int),
		Tag:          govultr.StringToStringPtr(d.Get("tag").(string)),
		AutoScaler:   govultr.BoolToBoolPtr(d.Get("auto_scaler").(bool)),
		MinNodes:     d.Get("min_nodes").(int),
		MaxNodes:     d.Get("max_nodes").(int),
	}

	// Handle Labels if provided or changed
	if d.HasChange("labels") {
		labelMap := make(map[string]string)
		if labels, ok := d.GetOk("labels"); ok {
			for k, v := range labels.(map[string]interface{}) {
				labelMap[k] = v.(string)
			}
		}
		req.Labels = labelMap
	}

	// Handle Taints if provided or changed
	if d.HasChange("taints") {
		taints := d.Get("taints").([]interface{})
		reqTaints := make([]govultr.Taint, 0, len(taints))

		for _, t := range taints {
			taintMap := t.(map[string]interface{})
			reqTaints = append(reqTaints, govultr.Taint{
				Key:    taintMap["key"].(string),
				Value:  taintMap["value"].(string),
				Effect: taintMap["effect"].(string),
			})
		}
		req.Taints = reqTaints
	}

	if _, _, err := client.Kubernetes.UpdateNodePool(ctx, clusterID, d.Id(), req); err != nil {
		return diag.Errorf("error updating VKE node pool %v : %v", d.Id(), err)
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
