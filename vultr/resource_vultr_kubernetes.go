package vultr

import (
	"context"
	"fmt"
	"log"
	"slices"
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
		Schema:        resourceVultrKubernetesV1(),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceVultrKubernetesV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceVultrKubernetesStateUpgradeV0ToV1,
				Version: 0,
			},
		},
	}
}

func resourceVultrKubernetesStateUpgradeV0ToV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) { //nolint:lll
	if len(rawState) == 0 {
		return rawState, nil
	}

	migrateState := rawState["node_pools"].([]interface{})[0].(map[string]interface{})
	migrateState["cluster_id"] = rawState["id"]

	nps, err := resourceVultrKubernetesNodePoolsStateUpgradeV0ToV1(
		ctx,
		migrateState,
		meta,
	)
	if err != nil {
		log.Println("[ERROR] unable to migrate kubernetes node pool state")
		return rawState, err
	}

	rawState["node_pools"].([]interface{})[0] = nps

	return rawState, nil
}

func resourceVultrKubernetesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	nodePools := d.Get("node_pools").([]interface{})
	nodePool := nodePools[0].(map[string]interface{})

	nodePoolReq := []govultr.NodePoolReq{
		{
			NodeQuantity: nodePool["node_quantity"].(int),
			Label:        nodePool["label"].(string),
			Plan:         nodePool["plan"].(string),
			Tag:          tfVKEDefault,
			AutoScaler:   govultr.BoolToBoolPtr(nodePool["auto_scaler"].(bool)),
			MinNodes:     nodePool["min_nodes"].(int),
			MaxNodes:     nodePool["max_nodes"].(int),
			UserData:     nodePool["user_data"].(string),
		},
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

	labels := nodePool["labels"].(*schema.Set).List()
	for i := range labels {
		labelMap := labels[i].(map[string]interface{})
		_, _, err := client.Kubernetes.CreateNodePoolLabel(
			ctx,
			d.Id(),
			cluster.NodePools[0].ID,
			&govultr.NodePoolLabelReq{
				Key:   labelMap["key"].(string),
				Value: labelMap["value"].(string),
			},
		)

		if err != nil {
			return diag.Errorf(
				"error creating label %q for node pool %v on cluster %v : %v",
				labelMap["key"].(string),
				cluster.NodePools[0].Label,
				d.Id(),
				err,
			)
		}
	}

	taints := nodePool["taints"].(*schema.Set).List()
	for i := range taints {
		taintMap := taints[i].(map[string]interface{})
		_, _, err := client.Kubernetes.CreateNodePoolTaint(
			ctx,
			d.Id(),
			cluster.NodePools[0].ID,
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
				cluster.NodePools[0].Label,
				d.Id(),
				err,
			)
		}
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
			nodePoolData := flattenNodePool(&vke.NodePools[i])

			// add node pool labels
			labelData, _, err := client.Kubernetes.ListNodePoolLabels(ctx, d.Id(), vke.NodePools[i].ID)
			if err != nil {
				return diag.Errorf("error getting cluster (%v) node pool (%v) labels : %v", d.Id(), vke.NodePools[i].ID, err)
			}

			var labels []map[string]interface{}
			for j := range labelData {
				labels = append(labels, map[string]interface{}{
					"id":    labelData[j].ID,
					"key":   labelData[j].Key,
					"value": labelData[j].Value,
				})
			}

			nodePoolData[0]["labels"] = labels

			// add node pool taints
			taintData, _, err := client.Kubernetes.ListNodePoolTaints(ctx, d.Id(), vke.NodePools[i].ID)
			if err != nil {
				return diag.Errorf("error getting cluster (%v) node pool (%v) taints : %v", d.Id(), vke.NodePools[i].ID, err)
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

			nodePoolData[0]["taints"] = taints

			if err := d.Set("node_pools", nodePoolData); err != nil {
				return diag.Errorf("unable to set resource kubernetes `node_pools` read value: %v", err)
			}

			found = true
			break
		}
	}

	if !found {
		return diag.Errorf(`
unable to set resource kubernetes default node pool with tag %s for %v. 
You must set the default tag on one node pool before importing.`,
			tfVKEDefault,
			d.Id(),
		)
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

		oldNodePool := oldNP.([]interface{})[0]
		newNodePool := newNP.([]interface{})[0]

		oldNodePoolData := oldNodePool.(map[string]interface{})
		newNodePoolData := newNodePool.(map[string]interface{})

		req := &govultr.NodePoolReqUpdate{
			NodeQuantity: newNodePoolData["node_quantity"].(int),
		}

		if newNodePoolData["auto_scaler"] != oldNodePoolData["auto_scaler"] {
			req.AutoScaler = govultr.BoolToBoolPtr(newNodePoolData["auto_scaler"].(bool))
		}

		if newNodePoolData["min_nodes"] != oldNodePoolData["min_nodes"] {
			req.MinNodes = newNodePoolData["min_nodes"].(int)
		}

		if newNodePoolData["max_nodes"] != oldNodePoolData["max_nodes"] {
			req.MaxNodes = newNodePoolData["max_nodes"].(int)
		}

		if newNodePoolData["user_data"] != oldNodePoolData["user_data"] {
			req.UserData = govultr.StringToStringPtr(newNodePoolData["user_data"].(string))
		}

		if _, _, err := client.Kubernetes.UpdateNodePool(ctx, d.Id(), newNodePoolData["id"].(string), req); err != nil {
			return diag.Errorf("error updating vke node pool %v : %v", d.Id(), err)
		}

		if d.HasChange("node_pools.0.labels") {
			oldLabels := oldNodePoolData["labels"].(*schema.Set).List()
			newLabels := newNodePoolData["labels"].(*schema.Set).List()

			err := updateNodePoolOptions(ctx, client, d.Id(), newNodePoolData["id"].(string), "labels", oldLabels, newLabels)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if d.HasChange("node_pools.0.taints") {
			oldTaints := oldNodePoolData["taints"].(*schema.Set).List()
			newTaints := newNodePoolData["taints"].(*schema.Set).List()

			err := updateNodePoolOptions(ctx, client, d.Id(), newNodePoolData["id"].(string), "taints", oldTaints, newTaints)
			if err != nil {
				return diag.FromErr(err)
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
		"user_data":     np.UserData,
	}

	nodePools = append(nodePools, pool)

	return nodePools
}

func updateNodePoolOptions(ctx context.Context, client *govultr.Client, clusterID, nodePoolID, optionKind string, oldData, newData []interface{}) error { //nolint:lll
	type optionData struct {
		Create   bool
		Delete   bool
		OptionID string
		Key      string
		Value    string
		Effect   string
	}

	oldOptionData := []optionData{}
	optionRequests := []optionData{}
	for i := range oldData {
		oldOption := oldData[i].(map[string]interface{})

		oldData := optionData{
			OptionID: oldOption["id"].(string),
			Key:      oldOption["key"].(string),
		}

		oldOptionData = append(oldOptionData, oldData)
	}

	for i := range newData {
		newOption := newData[i].(map[string]interface{})

		newRequest := optionData{
			Key:    newOption["key"].(string),
			Value:  newOption["value"].(string),
			Create: true,
			Delete: false,
		}

		if optionKind == "taints" {
			newRequest.Effect = newOption["effect"].(string)
		}

		oldIndex := slices.IndexFunc(oldOptionData, func(o optionData) bool { return o.Key == newRequest.Key })

		if oldIndex >= 0 {
			// delete the old option in the process
			newRequest.Delete = true
			newRequest.OptionID = oldOptionData[oldIndex].OptionID
		}

		optionRequests = append(optionRequests, newRequest)
	}

	// mark delete options not in changed data
	for i := range oldOptionData {
		if !slices.ContainsFunc(optionRequests, func(o optionData) bool { return o.Key == oldOptionData[i].Key }) {
			optionRequests = append(optionRequests, optionData{
				OptionID: oldOptionData[i].OptionID,
				Delete:   true,
				Create:   false,
			})
		}
	}

	for i := range optionRequests {
		switch optionKind {
		case "labels":
			if optionRequests[i].Delete {
				err := client.Kubernetes.DeleteNodePoolLabel(
					ctx,
					clusterID,
					nodePoolID,
					optionRequests[i].OptionID,
				)
				if err != nil {
					return fmt.Errorf(
						"error deleting label %q from vke %q node pool %q during option update : %w",
						optionRequests[i].Key,
						clusterID,
						nodePoolID,
						err,
					)
				}
			}

			if optionRequests[i].Create {
				_, _, err := client.Kubernetes.CreateNodePoolLabel(
					ctx,
					clusterID,
					nodePoolID,
					&govultr.NodePoolLabelReq{
						Key:   optionRequests[i].Key,
						Value: optionRequests[i].Value,
					},
				)
				if err != nil {
					return fmt.Errorf(
						"error creating label %q on vke %q node pool %q during option update : %w",
						optionRequests[i].Key,
						clusterID,
						nodePoolID,
						err,
					)
				}
			}
		case "taints":
			if optionRequests[i].Delete {
				err := client.Kubernetes.DeleteNodePoolTaint(
					ctx,
					clusterID,
					nodePoolID,
					optionRequests[i].OptionID,
				)
				if err != nil {
					return fmt.Errorf(
						"error deleting taint %q from vke %q node pool %q during option update : %w",
						optionRequests[i].Key,
						clusterID,
						nodePoolID,
						err,
					)
				}
			}

			if optionRequests[i].Create {
				_, _, err := client.Kubernetes.CreateNodePoolTaint(
					ctx,
					clusterID,
					nodePoolID,
					&govultr.NodePoolTaintReq{
						Key:    optionRequests[i].Key,
						Value:  optionRequests[i].Value,
						Effect: optionRequests[i].Effect,
					},
				)
				if err != nil {
					return fmt.Errorf(
						"error creating label %q on vke %q node pool %q during update : %w",
						optionRequests[i].Key,
						clusterID,
						nodePoolID,
						err,
					)
				}
			}
		}
	}

	return nil
}
