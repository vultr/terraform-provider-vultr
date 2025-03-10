package vultr

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrObjectStorageTier() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrObjectStorageTierRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"price": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rate_limit_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rate_limit_operations": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrObjectStorageTierRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	tierList := []govultr.ObjectStorageTier{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	tiers, _, err := client.ObjectStorage.ListTiers(ctx)
	if err != nil {
		return diag.Errorf("Error getting object storage tier list : %v", err)
	}

	for i := range tiers {
		sm, err := structToMap(tiers[i])

		if err != nil {
			return diag.FromErr(err)
		}

		if filterLoop(f, sm) {
			tierList = append(tierList, tiers[i])
		}
	}

	if len(tierList) > 1 {
		return diag.Errorf(`your object storage tier search returned too many results. 
Please refine your search to be more specific`)
	}

	if len(tierList) < 1 {
		return diag.Errorf("no object storage tier results were found")
	}

	d.SetId(strconv.Itoa(tierList[0].ID))
	if err := d.Set("price", tierList[0].Price); err != nil {
		return diag.Errorf("unable to set object storage tier `price` read value: %v", err)
	}
	if err := d.Set("slug", tierList[0].Slug); err != nil {
		return diag.Errorf("unable to set object storage tier `slug` read value: %v", err)
	}
	if err := d.Set("rate_limit_bytes", tierList[0].RateLimitBytesSec); err != nil {
		return diag.Errorf("unable to set object storage tier `rate_limit_bytes` read value: %v", err)
	}
	if err := d.Set("rate_limit_operations", tierList[0].RateLimitOpsSec); err != nil {
		return diag.Errorf("unable to set object storage tier `rate_limit_operations` read value: %v", err)
	}

	var tierLocs []map[string]interface{}
	for j := range tierList[0].Locations {
		tierLocs = append(tierLocs, map[string]interface{}{
			"id":       tierList[0].Locations[j].ID,
			"region":   tierList[0].Locations[j].Region,
			"hostname": tierList[0].Locations[j].Hostname,
			"name":     tierList[0].Locations[j].Name,
		})
	}

	if err := d.Set("locations", tierLocs); err != nil {
		return diag.Errorf("unable to set object storage tier `locations` read value: %v", err)
	}

	return nil
}
