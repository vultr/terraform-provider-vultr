package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrStartupScript() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrStartupScriptRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"script": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrStartupScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var scriptList []govultr.StartupScript
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		scripts, meta, err := client.StartupScript.List(ctx, options)

		if err != nil {
			return diag.Errorf("error getting startup scripts: %v", err)
		}

		for _, script := range scripts {
			sm, err := structToMap(script)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				scriptList = append(scriptList, script)
			}
		}
		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}
	if len(scriptList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(scriptList) < 1 {
		return diag.Errorf("no results were found")
	}

	// The script field is not returned in the list call but only in the get.
	script, err := client.StartupScript.Get(ctx, scriptList[0].ID)
	if err != nil {
		return diag.Errorf("error retrieving script : %s", scriptList[0])
	}

	d.SetId(script.ID)
	d.Set("name", script.Name)
	d.Set("date_created", script.DateCreated)
	d.Set("date_modified", script.DateModified)
	d.Set("type", script.Type)
	d.Set("script", script.Script)
	return nil
}
