package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrIsoPublic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrIsoPublicRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"md5sum": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrIsoPublicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOK := d.GetOk("filter")

	if !filtersOK {
		return diag.Errorf("issue with filter: %v", filtersOK)
	}

	isoList := []govultr.PublicISO{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		iso, meta,_, err := client.ISO.ListPublic(ctx, options)
		if err != nil {
			return diag.Errorf("Error getting isos: %v", err)
		}

		for _, i := range iso {
			sm, err := structToMap(i)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				isoList = append(isoList, i)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(isoList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(isoList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(isoList[0].ID)
	if err := d.Set("description", isoList[0].Description); err != nil {
		return diag.Errorf("unable to set iso_public `description` read value: %v", err)
	}
	if err := d.Set("name", isoList[0].Name); err != nil {
		return diag.Errorf("unable to set iso_public `name` read value: %v", err)
	}
	if err := d.Set("md5sum", isoList[0].MD5Sum); err != nil {
		return diag.Errorf("unable to set iso_public `md5sum` read value: %v", err)
	}
	return nil
}
