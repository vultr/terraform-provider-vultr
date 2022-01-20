package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrPrivateNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrPrivateNetworkRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v4_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v4_subnet_mask": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrPrivateNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var networkList []govultr.Network
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		networks, meta, err := client.Network.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting networks: %v", err)
		}

		for _, n := range networks {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(n)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				networkList = append(networkList, n)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(networkList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(networkList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(networkList[0].NetworkID)
	d.Set("region", networkList[0].Region)
	d.Set("description", networkList[0].Description)
	d.Set("date_created", networkList[0].DateCreated)
	d.Set("v4_subnet", networkList[0].V4Subnet)
	d.Set("v4_subnet_mask", networkList[0].V4SubnetMask)

	return nil
}
