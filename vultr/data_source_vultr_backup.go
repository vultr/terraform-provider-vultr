package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrBackup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrBackupRead,
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

func dataSourceVultrBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	backupList := []govultr.Backup{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		backups, meta, err := client.Backup.List(context.Background(), options)
		if err != nil {
			return fmt.Errorf("Error getting backups: %v", err)
		}

		for _, b := range backups {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(b)
			if err != nil {
				return err
			}

			if filterLoop(f, sm) {
				backupList = append(backupList, b)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(backupList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(backupList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(backupList[0].ID)
	d.Set("description", backupList[0].Description)
	d.Set("date_created", backupList[0].DateCreated)
	d.Set("size", backupList[0].Size)
	d.Set("status", backupList[0].Status)

	return nil
}
