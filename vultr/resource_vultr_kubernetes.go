package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
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
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: IgnoreCase,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ha_controlplanes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"node_pools": {
				Type:     schema.TypeList,
				Optional: true,
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
				Description: "Base64 encoded KubeConfig",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"cluster_ca_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"client_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"client_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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
		Label:           d.Get("label").(string),
		Region:          d.Get("region").(string),
		Version:         d.Get("version").(string),
		HAControlPlanes: d.Get("ha_controlplanes").(bool),
		NodePools:       nodePoolReq,
	}

	cluster, _, err := client.Kubernetes.CreateCluster(ctx, req)
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

	vke, _, err := client.Kubernetes.GetCluster(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Unauthorized") {
			return diag.Errorf("API authorization error: %v", err)
		}
		if strings.Contains(err.Error(), "Invalid resource ID") {
			log.Printf("[WARN] Kubernetes Cluster (%v) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting cluster (%s): %v", d.Id(), err)
	}

	// Look for the node pool with the tag `tf-vke-default`
	for _, v := range vke.NodePools {
		if tfVKEDefault == v.Tag {
			if err := d.Set("node_pools", flattenNodePool(&v)); err != nil {
				return diag.Errorf("unable to set resource kubernetes `node_pool` read value: %v", err)
			}
			break
		}
	}

	if err := d.Set("date_created", vke.DateCreated); err != nil {
		return diag.Errorf("unable to set resource kubernetes `date_created` read value: %v", err)
	}
	if err := d.Set("cluster_subnet", vke.ClusterSubnet); err != nil {
		return diag.Errorf("unable to set resource kubernetes `cluster_subnet` read value: %v", err)
	}
	if err := d.Set("service_subnet", vke.ServiceSubnet); err != nil {
		return diag.Errorf("unable to set resource kubernetes `service_subnet` read value: %v", err)
	}
	if err := d.Set("ip", vke.IP); err != nil {
		return diag.Errorf("unable to set resource kubernetes `ip` read value: %v", err)
	}
	if err := d.Set("endpoint", vke.Endpoint); err != nil {
		return diag.Errorf("unable to set resource kubernetes `endpoint` read value: %v", err)
	}
	if err := d.Set("status", vke.Status); err != nil {
		return diag.Errorf("unable to set resource kubernetes `status` read value: %v", err)
	}

	config, _, err := client.Kubernetes.GetKubeConfig(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not get kubeconfig : %v", err)
	}

	ca, cert, key, err := getCertsFromKubeConfig(config.KubeConfig)
	if err != nil {
		return diag.Errorf("error getting certs from kubeconfig : %v", err)
	}

	if err := d.Set("kube_config", config.KubeConfig); err != nil {
		return diag.Errorf("unable to set resource kubernetes `kube_config` read value: %v", err)
	}
	if err := d.Set("cluster_ca_certificate", ca); err != nil {
		return diag.Errorf("unable to set kubernetes `cluster_ca_certificate` read value: %v", err)
	}
	if err := d.Set("client_certificate", cert); err != nil {
		return diag.Errorf("unable to set kubernetes `client_certificate` read value: %v", err)
	}
	if err := d.Set("client_key", key); err != nil {
		return diag.Errorf("unable to set kubernetes `client_key` read value: %v", err)
	}
	if err := d.Set("version", vke.Version); err != nil {
		return diag.Errorf("unable to set resource kubernetes `version` read value: %v", err)
	}
	if err := d.Set("ha_controlplanes", vke.HAControlPlanes); err != nil {
		return diag.Errorf("unable to set resource kubernetes `ha_controlplanes` read value: %v", err)
	}

	return nil
}

func resourceVultrKubernetesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	if d.HasChange("label") {
		req := &govultr.ClusterReqUpdate{}
		req.Label = d.Get("label").(string)

		if err := client.Kubernetes.UpdateCluster(ctx, d.Id(), req); err != nil {
			return diag.Errorf("error updating vke cluster (%v): %v", d.Id(), err)
		}
	}

	if d.HasChange("node_pools") {

		oldNP, newNP := d.GetChange("node_pools")

		if len(newNP.([]interface{})) != 0 && len(oldNP.([]interface{})) != 0 {

			n := newNP.([]interface{})[0].(map[string]interface{})

			req := &govultr.NodePoolReqUpdate{
				NodeQuantity: n["node_quantity"].(int),
				AutoScaler:   govultr.BoolToBoolPtr(n["auto_scaler"].(bool)),
				MinNodes:     n["min_nodes"].(int),
				MaxNodes:     n["max_nodes"].(int),
				// Not updating tag for default node pool since it's needed to lookup in terraform
			}

			if _, _, err := client.Kubernetes.UpdateNodePool(ctx, d.Id(), n["id"].(string), req); err != nil {
				return diag.Errorf("error updating VKE node pool %v : %v", d.Id(), err)
			}
		} else if len(newNP.([]interface{})) == 0 && len(oldNP.([]interface{})) != 0 {
			// if we have an old node pool state but don't have a new node pool state
			// we can safely assume this is a node pool removal

			n := oldNP.([]interface{})[0].(map[string]interface{})

			if err := client.Kubernetes.DeleteNodePool(ctx, d.Id(), n["id"].(string)); err != nil {
				return diag.Errorf("error deleting VKE node pool %v : %v", d.Id(), err)
			}
		} else if len(newNP.([]interface{})) != 0 && len(oldNP.([]interface{})) == 0 {
			// if we don't have an old node pool state but have a new node pool state
			// we can safely assume this is a new node pool creation

			n := newNP.([]interface{})[0].(map[string]interface{})

			req := &govultr.NodePoolReq{
				NodeQuantity: n["node_quantity"].(int),
				Tag:          tfVKEDefault,
				Plan:         n["plan"].(string),
				Label:        n["label"].(string),
			}

			if _, _, err := client.Kubernetes.CreateNodePool(ctx, d.Id(), req); err != nil {
				return diag.Errorf("error creating VKE node pool %v : %v", d.Id(), err)
			}
		}
	}

	// k8s version upgrade
	if d.HasChange("version") {
		upgradeReq := &govultr.ClusterUpgradeReq{
			UpgradeVersion: d.Get("version").(string),
		}

		if err := client.Kubernetes.Upgrade(ctx, d.Id(), upgradeReq); err != nil {
			return diag.Errorf("error upgrading VKE cluster %v : %v", d.Id(), err)
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
			AutoScaler:   govultr.BoolToBoolPtr(r["auto_scaler"].(bool)),
			MinNodes:     r["min_nodes"].(int),
			MaxNodes:     r["max_nodes"].(int),
		}

		npr = append(npr, t)
	}
	return npr
}

func waitForVKEAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for kubernetes cluster (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{ // nolint:all
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

func newVKEStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc { // nolint:all
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating kubernetes cluster")

		vke, _, err := client.Kubernetes.GetCluster(ctx, d.Id())
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
		"auto_scaler":   np.AutoScaler,
		"min_nodes":     np.MinNodes,
		"max_nodes":     np.MaxNodes,
	}

	nodePools = append(nodePools, pool)

	return nodePools
}
