package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrOrganizationRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
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

func dataSourceVultrOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var orgList []govultr.Organization
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		list, meta, _, err := client.Organization.ListOrganizations(ctx, options)
		if err != nil {
			return diag.Errorf("error getting organizations : %v", err)
		}

		for i := range list {
			sm, err := structToMap(list[i])

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				orgList = append(orgList, list[i])
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(orgList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(orgList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(orgList[0].ID)
	if err := d.Set("name", orgList[0].Name); err != nil {
		return diag.Errorf("unable to set data source organization `name` read value: %v", err)
	}
	if err := d.Set("type", orgList[0].Type); err != nil {
		return diag.Errorf("unable to set data source organization `type` read value: %v", err)
	}
	if err := d.Set("date_created", orgList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set data source organization `date_created` read value: %v", err)
	}

	return nil
}
