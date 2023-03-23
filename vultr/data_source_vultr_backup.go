package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrBackup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrBackupRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func dataSourceVultrBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var backupList []map[string]interface{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		backups, meta,_, err := client.Backup.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting backups: %v", err)
		}

		for _, b := range backups {
			// We need convert the struct into a map. This allows us to easily manipulate the data here.
			sm, err := structToMap(b)
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				backupList = append(backupList, sm)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(backupList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(backupList[0]["description"].(string))
	if err := d.Set("backups", backupList); err != nil {
		return diag.Errorf("error setting `backups`: %#v", err)
	}

	return nil
}
