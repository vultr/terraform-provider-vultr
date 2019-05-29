package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrBackup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrBackupRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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

	backups, err := client.Backup.List(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting applications: %v", err)
	}

	backupList := []govultr.Backup{}

	f := buildVultrDataSourceFilter(filters.(*schema.Set))

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

	if len(backupList) > 1 {
		return fmt.Errorf("your search returned too many results : %d. Please refine your search to be more specific", len(backupList))
	}

	if len(backupList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(backupList[0].BackupID)
	d.Set("date_created", backupList[0].DateCreated)
	d.Set("description", backupList[0].Description)
	d.Set("size", backupList[0].Size)
	d.Set("status", backupList[0].Status)
	return nil
}
