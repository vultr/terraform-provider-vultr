package vultr

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrOS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrOSRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"family": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrOSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	osList := []govultr.OS{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		os, meta, err := client.OS.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting os list: %v", err)
		}

		for _, o := range os {
			sm, err := structToMap(o)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				osList = append(osList, o)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(osList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(osList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(strconv.Itoa(osList[0].ID))
	if err := d.Set("name", osList[0].Name); err != nil {
		return diag.Errorf("unable to set os `name` read value: %v", err)
	}
	if err := d.Set("arch", osList[0].Arch); err != nil {
		return diag.Errorf("unable to set os `arch` read value: %v", err)
	}
	if err := d.Set("family", osList[0].Family); err != nil {
		return diag.Errorf("unable to set os `family` read value: %v", err)
	}
	return nil
}
