package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrKubernetes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrKubernetesRead,
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
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"kube_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	var k8List []govultr.Cluster
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		k8s, meta, err := client.Kubernetes.ListClusters(context.Background(), options)
		if err != nil {
			return fmt.Errorf("error getting kubernetes")
		}

		for _, k8 := range k8s {
			sm, err := structToMap(k8)

			if err != nil {
				return err
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
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(k8List) < 1 {
		return errors.New("no results were found")
	}

	kubeConfig, err := client.Kubernetes.GetKubeConfig(context.Background(), k8List[0].ID)
	if err != nil {
		return errors.New("error getting kubeconfig")
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
	d.Set("node_pools", k8List[0].NodePools)
	d.Set("kube_config", kubeConfig)

	return nil
}
