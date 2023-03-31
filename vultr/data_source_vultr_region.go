package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrRegionRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"continent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"city": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVultrRegionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	regionList := []govultr.Region{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		regions, meta, _, err := client.Region.List(ctx, options)
		if err != nil {
			return diag.Errorf("Error getting regions: %v", err)
		}

		for _, a := range regions {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(a)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				regionList = append(regionList, a)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(regionList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(regionList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(regionList[0].ID)
	if err := d.Set("country", regionList[0].Country); err != nil {
		return diag.Errorf("unable to set region `country` read value: %v", err)
	}
	if err := d.Set("continent", regionList[0].Continent); err != nil {
		return diag.Errorf("unable to set region `continent` read value: %v", err)
	}
	if err := d.Set("city", regionList[0].City); err != nil {
		return diag.Errorf("unable to set region `city` read value: %v", err)
	}
	if err := d.Set("options", regionList[0].Options); err != nil {
		return diag.Errorf("unable to set region `options` read value: %v", err)
	}
	return nil
}
