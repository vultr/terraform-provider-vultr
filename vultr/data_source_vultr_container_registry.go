package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrContainerRegistry() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrContainerRegistryRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"root_user": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	crList := []govultr.ContainerRegistry{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{PerPage: 10}

	for {
		crs, meta, _, err := client.ContainerRegistry.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting container registries: %v", err)
		}

		for _, u := range crs {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(u)
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				crList = append(crList, u)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(crList) > 1 {
		return diag.Errorf("your search returned too many results : %d. Please refine your search to be more specific", len(crList))
	}
	if len(crList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(crList[0].ID)
	if err := d.Set("name", crList[0].Name); err != nil {
		return diag.Errorf("unable to set container registry `name` read value: %v", err)
	}
	if err := d.Set("urn", crList[0].URN); err != nil {
		return diag.Errorf("unable to set container registry `urn` read value: %v", err)
	}
	if err := d.Set("public", crList[0].Public); err != nil {
		return diag.Errorf("unable to set container registry `public` read value: %v", err)
	}
	if err := d.Set("storage", flattenCRStorage(&crList[0])); err != nil {
		return diag.Errorf("unable to set container registry `storage` read value: %v", err)
	}
	if err := d.Set("root_user", flattenCRRootUser(&crList[0])); err != nil {
		return diag.Errorf("unable to set container registry `root_user` read value: %v", err)
	}
	if err := d.Set("date_created", crList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set container registry `date_created` read value: %v", err)
	}
	return nil
}
