package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrKubernetes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrKubernetesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ha_controlplanes": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"firewall_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_scaler": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"min_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"taints": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"effect": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"date_created": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"label": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"kube_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_ca_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrKubernetesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var k8List []govultr.Cluster
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		k8s, meta, _, err := client.Kubernetes.ListClusters(ctx, options)
		if err != nil {
			return diag.Errorf("error getting kubernetes")
		}

		for _, k8 := range k8s {
			sm, err := structToMap(k8)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				k8List = append(k8List, k8)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(k8List) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(k8List) < 1 {
		return diag.Errorf("no results were found")
	}

	kubeConfig, _, err := client.Kubernetes.GetKubeConfig(ctx, k8List[0].ID)
	if err != nil {
		return diag.Errorf("error getting kubeconfig")
	}

	ca, cert, key, err := getCertsFromKubeConfig(kubeConfig.KubeConfig)
	if err != nil {
		return diag.Errorf("error getting certs from kubeconfig : %v", err)
	}

	d.SetId(k8List[0].ID)
	if err := d.Set("label", k8List[0].Label); err != nil {
		return diag.Errorf("unable to set kubernetes `label` read value: %v", err)
	}
	if err := d.Set("date_created", k8List[0].DateCreated); err != nil {
		return diag.Errorf("unable to set kubernetes `date_created` read value: %v", err)
	}
	if err := d.Set("cluster_subnet", k8List[0].ClusterSubnet); err != nil {
		return diag.Errorf("unable to set kubernetes `cluster_subnet` read value: %v", err)
	}
	if err := d.Set("service_subnet", k8List[0].ServiceSubnet); err != nil {
		return diag.Errorf("unable to set kubernetes `service_subnet` read value: %v", err)
	}
	if err := d.Set("ip", k8List[0].IP); err != nil {
		return diag.Errorf("unable to set kubernetes `ip` read value: %v", err)
	}
	if err := d.Set("endpoint", k8List[0].Endpoint); err != nil {
		return diag.Errorf("unable to set kubernetes `endpoint` read value: %v", err)
	}
	if err := d.Set("version", k8List[0].Version); err != nil {
		return diag.Errorf("unable to set kubernetes `version` read value: %v", err)
	}
	if err := d.Set("ha_controlplanes", k8List[0].HAControlPlanes); err != nil {
		return diag.Errorf("unable to set kubernetes `ha_controlplanes` read value: %v", err)
	}
	if err := d.Set("firewall_group_id", k8List[0].FirewallGroupID); err != nil {
		return diag.Errorf("unable to set kubernetes `firewall_group_id` read value: %v", err)
	}
	if err := d.Set("region", k8List[0].Region); err != nil {
		return diag.Errorf("unable to set kubernetes `region` read value: %v", err)
	}
	if err := d.Set("status", k8List[0].Status); err != nil {
		return diag.Errorf("unable to set kubernetes `status` read value: %v", err)
	}
	if err := d.Set("kube_config", kubeConfig.KubeConfig); err != nil {
		return diag.Errorf("unable to set kubernetes `kube_config` read value: %v", err)
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

	nodePools := flattenNodePools(k8List[0].NodePools)

	for i := range nodePools {
		labelData, _, err := client.Kubernetes.ListNodePoolLabels(ctx, k8List[0].ID, nodePools[i]["id"].(string))
		if err != nil {
			return diag.Errorf(
				"error getting data source cluster (%v) node pool (%v) labels : %v",
				k8List[0].ID,
				nodePools[i]["id"].(string),
				err,
			)
		}

		var labels []map[string]interface{}
		for j := range labelData {
			labels = append(labels, map[string]interface{}{
				"id":    labelData[j].ID,
				"key":   labelData[j].Key,
				"value": labelData[j].Value,
			})
		}

		nodePools[i]["labels"] = labels

		taintData, _, err := client.Kubernetes.ListNodePoolTaints(ctx, k8List[0].ID, nodePools[i]["id"].(string))
		if err != nil {
			return diag.Errorf(
				"error getting data source cluster (%v) node pool (%v) taints : %v",
				k8List[0].ID,
				nodePools[i]["id"].(string),
				err,
			)
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

		nodePools[i]["taints"] = taints
	}

	if err := d.Set("node_pools", nodePools); err != nil {
		return diag.Errorf("unable to set kubernetes `node_pools` read value: %v", err)
	}

	return nil
}

func flattenNodePools(np []govultr.NodePool) []map[string]interface{} {
	var nodePools []map[string]interface{}

	for _, n := range np {
		var instances []map[string]interface{}

		for _, v := range n.Nodes {
			a := map[string]interface{}{
				"id":           v.ID,
				"status":       v.Status,
				"date_created": v.DateCreated,
				"label":        v.Label,
			}
			instances = append(instances, a)
		}

		pool := map[string]interface{}{
			"label":         n.Label,
			"plan":          n.Plan,
			"node_quantity": n.NodeQuantity,
			"id":            n.ID,
			"date_created":  n.DateCreated,
			"date_updated":  n.DateUpdated,
			"status":        n.Status,
			"tag":           n.Tag,
			"auto_scaler":   n.AutoScaler,
			"min_nodes":     n.MinNodes,
			"max_nodes":     n.MaxNodes,
			"nodes":         instances,
			"user_data":     n.UserData,
		}

		nodePools = append(nodePools, pool)
	}

	return nodePools
}
