package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrStartupScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrStartupScriptRead,
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

func dataSourceVultrStartupScriptRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	scripts, err := client.StartupScript.List(context.Background())

	if err != nil {
		return fmt.Errorf("error getting startup scripts: %v", err)
	}

	scriptList := []govultr.StartupScript{}

	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, script := range scripts {
		sm, err := structToMap(script)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			scriptList = append(scriptList, script)
		}
	}

	if len(scriptList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(scriptList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(scriptList[0].ScriptID)
	d.Set("name", scriptList[0].Name)
	d.Set("date_created", scriptList[0].DateCreated)
	d.Set("date_modified", scriptList[0].DateModified)
	d.Set("type", scriptList[0].Type)
	d.Set("script", scriptList[0].Script)
	return nil
}
