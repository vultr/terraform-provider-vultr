package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
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
					Schema: nodePoolSchema(false),
				},
			},
			"kube_config": {
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
		k8s, meta, err := client.Kubernetes.ListClusters(ctx, options)
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

	kubeConfig, err := client.Kubernetes.GetKubeConfig(ctx, k8List[0].ID)
	if err != nil {
		return diag.Errorf("error getting kubeconfig")
	}

	d.SetId(k8List[0].ID)
	d.Set("label", k8List[0].Label)
	d.Set("date_created", k8List[0].DateCreated)
	d.Set("cluster_subnet", k8List[0].ClusterSubnet)
	d.Set("service_subnet", k8List[0].ServiceSubnet)
	d.Set("ip", k8List[0].IP)
	d.Set("endpoint", k8List[0].Endpoint)
	d.Set("version", k8List[0].Version)
	d.Set("region", k8List[0].Region)
	d.Set("status", k8List[0].Status)
	d.Set("kube_config", kubeConfig.KubeConfig)

	if err := d.Set("node_pools", flattenNodePools(k8List[0].NodePools)); err != nil {
		return diag.FromErr(err)
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
			"nodes":         instances,
		}

		nodePools = append(nodePools, pool)
	}

	return nodePools
}
