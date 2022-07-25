package vultr

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrObjectStorageClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrObjectStorageClustersRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deploy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrObjectStorageClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	clusterList := []govultr.ObjectStorageCluster{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		clusters, meta, err := client.ObjectStorage.ListCluster(ctx, options)
		if err != nil {
			return diag.Errorf("Error getting plans: %v", err)
		}

		for _, a := range clusters {
			// we need convert the  struct INTO a map allowing for easy manipulation of the data here
			sm, err := structToMap(a)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				clusterList = append(clusterList, a)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(clusterList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(clusterList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(strconv.Itoa(clusterList[0].ID))
	d.Set("region", clusterList[0].Region)
	d.Set("hostname", clusterList[0].Hostname)
	d.Set("deploy", clusterList[0].Deploy)

	return nil
}
