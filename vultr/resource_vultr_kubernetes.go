package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

var tfVKEDefault = "tf-vke-default"

func resourceVultrKubernetes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrKubernetesCreate,
		ReadContext:   resourceVultrKubernetesRead,
		UpdateContext: resourceVultrKubernetesUpdate,
		DeleteContext: resourceVultrKubernetesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"node_pools": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: nodePoolSchema(false),
				},
			},

			// Computed fields
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrKubernetesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	var nodePoolReq []govultr.NodePoolReq
	if np, npOk := d.GetOk("node_pools"); npOk {
		nodePoolReq = generateNodePool(np)
	} else {
		nodePoolReq = nil
	}

	req := &govultr.ClusterReq{
		Label:     d.Get("label").(string),
		Region:    d.Get("region").(string),
		Version:   d.Get("version").(string),
		NodePools: nodePoolReq,
	}

	cluster, err := client.Kubernetes.CreateCluster(ctx, req)
	if err != nil {
		return diag.Errorf("error creating kubernetes cluster: %v", err)
	}

	d.SetId(cluster.ID)

	//block until status is ready
	if _, err = waitForVKEAvailable(ctx, d, "active", []string{"pending"}, "status", meta); err != nil {
		return diag.Errorf(
			"error while waiting for kubernetes cluster %v to be completed: %v", cluster.ID, err)
	}

	return resourceVultrKubernetesRead(ctx, d, meta)
}

func resourceVultrKubernetesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vke, err := client.Kubernetes.GetCluster(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Kubernetes Cluster (%v) not found", d.Id())
		d.SetId("")
		return nil
	}

	found := false
	for _, v := range vke.NodePools {
		if tfVKEDefault == v.Tag {
			if err := d.Set("node_pools", flattenNodePool(&v)); err != nil {
				return diag.Errorf("error setting `node_pool`: %v", err)
			}
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("no default node_pool found : %v", err)
	}

	d.Set("date_created", vke.DateCreated)
	d.Set("cluster_subnet", vke.ClusterSubnet)
	d.Set("service_subnet", vke.ServiceSubnet)
	d.Set("ip", vke.IP)
	d.Set("endpoint", vke.Endpoint)
	d.Set("status", vke.Status)

	config, err := client.Kubernetes.GetKubeConfig(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not get kubeconfig : %v", err)
	}

	d.Set("kube_config", config.KubeConfig)

	return nil
}

func resourceVultrKubernetesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.ClusterReqUpdate{}

	if d.HasChange("label") {
		req.Label = d.Get("label").(string)
	}

	if err := client.Kubernetes.UpdateCluster(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating vke cluster (%v): %v", d.Id(), err)
	}

	if d.HasChange("node_pools") {

		oldNP, newNP := d.GetChange("node_pools")

		o := oldNP.([]interface{})[0].(map[string]interface{})
		n := newNP.([]interface{})[0].(map[string]interface{})

		if o["node_quantity"] != n["node_quantity"] {
			req := &govultr.NodePoolReqUpdate{NodeQuantity: n["node_quantity"].(int)}
			if _, err := client.Kubernetes.UpdateNodePool(ctx, d.Id(), o["id"].(string), req); err != nil {
				return diag.Errorf("error deleting VKE node pool %v : %v", o["id"], err)
			}
		}
	}

	return resourceVultrKubernetesRead(ctx, d, meta)
}

func resourceVultrKubernetesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Delete VKE : %v", d.Id())

	if err := client.Kubernetes.DeleteCluster(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting VKE %v : %v", d.Id(), err)
	}
	return nil
}

func generateNodePool(pools interface{}) []govultr.NodePoolReq {
	var npr []govultr.NodePoolReq
	pool := pools.([]interface{})
	for _, p := range pool {
		r := p.(map[string]interface{})

		t := govultr.NodePoolReq{
			NodeQuantity: r["node_quantity"].(int),
			Label:        r["label"].(string),
			Plan:         r["plan"].(string),
			Tag:          tfVKEDefault,
		}

		npr = append(npr, t)
	}
	return npr
}

func waitForVKEAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for kubernetes cluster (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newVKEStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     5 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newVKEStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating kubernetes cluster")

		vke, err := client.Kubernetes.GetCluster(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving kubernetes cluster %s ", d.Id())
		}

		if attr == "status" {
			log.Printf("[INFO] The kubernetes cluster Status is %v", vke.Status)
			return vke, vke.Status, nil
		}

		return nil, "", nil
	}
}

func flattenNodePool(np *govultr.NodePool) []map[string]interface{} {
	var nodePools []map[string]interface{}

	var instances []map[string]interface{}
	for _, v := range np.Nodes {
		n := map[string]interface{}{
			"id":           v.ID,
			"status":       v.Status,
			"date_created": v.DateCreated,
			"label":        v.Label,
		}
		instances = append(instances, n)
	}

	pool := map[string]interface{}{
		"label":         np.Label,
		"plan":          np.Plan,
		"node_quantity": np.NodeQuantity,
		"id":            np.ID,
		"date_created":  np.DateCreated,
		"date_updated":  np.DateUpdated,
		"status":        np.Status,
		"tag":           np.Tag,
		"nodes":         instances,
	}

	nodePools = append(nodePools, pool)

	return nodePools
}
