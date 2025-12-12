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
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ha_controlplanes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"enable_firewall": {
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
			"firewall_group_id": {
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
		VPCID:           d.Get("vpc_id").(string),
		HAControlPlanes: d.Get("ha_controlplanes").(bool),
		EnableFirewall:  d.Get("enable_firewall").(bool),
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
	found := false
	for i := range vke.NodePools {
		if tfVKEDefault == vke.NodePools[i].Tag {
			if err := d.Set("node_pools", flattenNodePool(&vke.NodePools[i])); err != nil {
				return diag.Errorf("unable to set resource kubernetes `node_pools` read value: %v", err)
			}
			found = true
			break
		}
	}
	if !found {
		return diag.Errorf(`unable to set resource kubernetes default node pool with tag %s for %v.
	You must set the default tag on one node pool before importing.`,
			tfVKEDefault, d.Id())
	}

	if err := d.Set("region", vke.Region); err != nil {
		return diag.Errorf("unable to set resource kubernetes `region` read value: %v", err)
	}
	if err := d.Set("label", vke.Label); err != nil {
		return diag.Errorf("unable to set resource kubernetes `label` read value: %v", err)
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
	if err := d.Set("firewall_group_id", vke.FirewallGroupID); err != nil {
		return diag.Errorf("unable to set resource kubernetes `firewall_group_id` read value: %v", err)
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
			// Existing node pool update
			oldN := oldNP.([]interface{})[0].(map[string]interface{})
			n := newNP.([]interface{})[0].(map[string]interface{})
			nodePoolID := n["id"].(string)

			// Check if labels changed
			oldLabels := make(map[string]string)
			for k, v := range oldN["labels"].(map[string]interface{}) {
				oldLabels[k] = v.(string)
			}
			newLabels := make(map[string]string)
			for k, v := range n["labels"].(map[string]interface{}) {
				newLabels[k] = v.(string)
			}
			labelsChanged := !labelsEqual(oldLabels, newLabels)

			// Check if taints changed
			oldTaints := extractTaintsFromSet(oldN["taints"].(*schema.Set))
			newTaints := extractTaintsFromSet(n["taints"].(*schema.Set))
			taintsChanged := !taintsEqual(oldTaints, newTaints)

			// Sync labels using dedicated endpoints if they changed
			if labelsChanged {
				if err := syncClusterNodePoolLabels(ctx, client, d.Id(), nodePoolID, newLabels); err != nil {
					return diag.Errorf("error syncing VKE node pool labels %v : %v", nodePoolID, err)
				}
			}

			// Sync taints using dedicated endpoints if they changed
			if taintsChanged {
				if err := syncClusterNodePoolTaints(ctx, client, d.Id(), nodePoolID, newTaints); err != nil {
					return diag.Errorf("error syncing VKE node pool taints %v : %v", nodePoolID, err)
				}
			}

			// Update other node pool properties (not labels/taints)
			req := &govultr.NodePoolReqUpdate{
				NodeQuantity: n["node_quantity"].(int),
				AutoScaler:   govultr.BoolToBoolPtr(n["auto_scaler"].(bool)),
				MinNodes:     n["min_nodes"].(int),
				MaxNodes:     n["max_nodes"].(int),
				// Not updating tag for default node pool since it's needed to lookup in terraform
				UserData: govultr.StringToStringPtr(n["user_data"].(string)),
			}

			if _, _, err := client.Kubernetes.UpdateNodePool(ctx, d.Id(), nodePoolID, req); err != nil {
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

			labels := make(map[string]string)
			for k, v := range n["labels"].(map[string]interface{}) {
				labels[k] = v.(string)
			}

			var taints []govultr.Taint
			taintsList := n["taints"].(*schema.Set).List()
			for i := range taintsList {
				taintMap := taintsList[i].(map[string]interface{})
				taints = append(taints, govultr.Taint{
					Key:    taintMap["key"].(string),
					Value:  taintMap["value"].(string),
					Effect: taintMap["effect"].(string),
				})
			}

			req := &govultr.NodePoolReq{
				NodeQuantity: n["node_quantity"].(int),
				Tag:          tfVKEDefault,
				Plan:         n["plan"].(string),
				Label:        n["label"].(string),
				Labels:       labels,
				Taints:       taints,
				UserData:     n["user_data"].(string),
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

// labelsEqual compares two label maps for equality
func labelsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

// taintInfo represents a taint for comparison purposes
type taintInfo struct {
	Key    string
	Value  string
	Effect string
}

// extractTaintsFromSet extracts taints from a schema.Set into a slice of taintInfo
func extractTaintsFromSet(s *schema.Set) []taintInfo {
	var taints []taintInfo
	for _, t := range s.List() {
		taint := t.(map[string]interface{})
		taints = append(taints, taintInfo{
			Key:    taint["key"].(string),
			Value:  taint["value"].(string),
			Effect: taint["effect"].(string),
		})
	}
	return taints
}

// taintsEqual compares two taint slices for equality
func taintsEqual(a, b []taintInfo) bool {
	if len(a) != len(b) {
		return false
	}
	// Create maps for comparison using key+effect as identifier
	aMap := make(map[string]string)
	for _, t := range a {
		aMap[t.Key+":"+t.Effect] = t.Value
	}
	for _, t := range b {
		if v, ok := aMap[t.Key+":"+t.Effect]; !ok || v != t.Value {
			return false
		}
	}
	return true
}

// syncClusterNodePoolLabels synchronizes node pool labels using the dedicated label endpoints
func syncClusterNodePoolLabels(ctx context.Context, client *govultr.Client, clusterID, nodePoolID string, desiredLabels map[string]string) error {
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

// syncClusterNodePoolTaints synchronizes node pool taints using the dedicated taint endpoints
func syncClusterNodePoolTaints(ctx context.Context, client *govultr.Client, clusterID, nodePoolID string, desiredTaints []taintInfo) error {
	// Get current taints from API
	currentTaints, _, err := client.Kubernetes.ListNodePoolTaints(ctx, clusterID, nodePoolID)
	if err != nil {
		return fmt.Errorf("error listing current taints: %v", err)
	}

	// Build map of current taints by key+effect for easy lookup
	currentTaintsByKeyEffect := make(map[string]govultr.NodePoolTaint)
	for _, taint := range currentTaints {
		mapKey := taint.Key + ":" + taint.Effect
		currentTaintsByKeyEffect[mapKey] = taint
	}

	// Build map of desired taints by key+effect
	desiredTaintsByKeyEffect := make(map[string]taintInfo)
	for _, taint := range desiredTaints {
		mapKey := taint.Key + ":" + taint.Effect
		desiredTaintsByKeyEffect[mapKey] = taint
	}

	// Delete taints that are no longer desired or have changed values
	for mapKey, currentTaint := range currentTaintsByKeyEffect {
		desired, exists := desiredTaintsByKeyEffect[mapKey]
		if !exists || desired.Value != currentTaint.Value {
			log.Printf("[INFO] Deleting taint %s (ID: %s) from node pool %s", mapKey, currentTaint.ID, nodePoolID)
			if err := client.Kubernetes.DeleteNodePoolTaint(ctx, clusterID, nodePoolID, currentTaint.ID); err != nil {
				return fmt.Errorf("error deleting taint %s: %v", mapKey, err)
			}
		}
	}

	// Create taints that are new or have changed values
	for mapKey, desired := range desiredTaintsByKeyEffect {
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

		labels := make(map[string]string)
		for k, v := range r["labels"].(map[string]interface{}) {
			labels[k] = v.(string)
		}

		var taints []govultr.Taint
		taintsList := r["taints"].(*schema.Set).List()
		for i := range taintsList {
			taintMap := taintsList[i].(map[string]interface{})
			taints = append(taints, govultr.Taint{
				Key:    taintMap["key"].(string),
				Value:  taintMap["value"].(string),
				Effect: taintMap["effect"].(string),
			})
		}

		t := govultr.NodePoolReq{
			NodeQuantity: r["node_quantity"].(int),
			Label:        r["label"].(string),
			Plan:         r["plan"].(string),
			Tag:          tfVKEDefault,
			AutoScaler:   govultr.BoolToBoolPtr(r["auto_scaler"].(bool)),
			MinNodes:     r["min_nodes"].(int),
			MaxNodes:     r["max_nodes"].(int),
			Labels:       labels,
			Taints:       taints,
			UserData:     r["user_data"].(string),
		}

		npr = append(npr, t)
	}

	return npr
}

func waitForVKEAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll
	log.Printf(
		"[INFO] Waiting for kubernetes cluster (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
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

func newVKEStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) retry.StateRefreshFunc { //nolint:lll
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

	labels := make(map[string]interface{})
	for k, v := range np.Labels {
		labels[k] = v
	}

	var taints []map[string]interface{}
	for i := range np.Taints {
		taints = append(taints, map[string]interface{}{
			"key":    np.Taints[i].Key,
			"value":  np.Taints[i].Value,
			"effect": np.Taints[i].Effect,
		})
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
		"labels":        labels,
		"taints":        taints,
		"user_data":     np.UserData,
	}

	nodePools = append(nodePools, pool)

	return nodePools
}
