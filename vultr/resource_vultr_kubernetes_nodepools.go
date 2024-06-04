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
		Schema:        nodePoolSchema(true),
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

	if _, _, err := client.Kubernetes.UpdateNodePool(ctx, clusterID, d.Id(), req); err != nil {
		return diag.Errorf("error deleting VKE node pool %v : %v", d.Id(), err)
	}

	return resourceVultrKubernetesNodePoolsRead(ctx, d, meta)
}

func resourceVultrKubernetesNodePoolsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	if err := client.Kubernetes.DeleteNodePool(ctx, d.Get("cluster_id").(string), d.Id()); err != nil {
		return diag.Errorf("error deleting VKE nodepool %v : %v", d.Id(), err)
	}

	return nil
}

func waitForNodePoolAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
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

func newNodePoolStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) retry.StateRefreshFunc {
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
