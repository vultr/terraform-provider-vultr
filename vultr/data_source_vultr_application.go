package vultr

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrApplicationRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"deploy_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"short_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	appList := []govultr.Application{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		apps, meta, err := client.Application.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting applications: %v", err)
		}

		for _, a := range apps {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(a)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				appList = append(appList, a)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}
	if len(appList) > 1 {
		return diag.Errorf("your search returned too many results : %d. Please refine your search to be more specific", len(appList))
	}

	if len(appList) < 1 {
		return diag.Errorf("no results were found")
	}
	d.SetId(strconv.Itoa(appList[0].ID))
	d.Set("deploy_name", appList[0].DeployName)
	d.Set("name", appList[0].Name)
	d.Set("short_name", appList[0].ShortName)
	d.Set("vendor", appList[0].Vendor)
	d.Set("image_id", appList[0].ImageID)
	d.Set("type", appList[0].Type)
	return nil
}
