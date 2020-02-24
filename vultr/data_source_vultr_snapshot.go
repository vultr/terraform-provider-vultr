package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrSnapshotRead,
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
			"os_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrSnapshotRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	snapshots, err := client.Snapshot.List(context.Background())

	if err != nil {
		return fmt.Errorf("error getting snapshots: %v", err)
	}

	snapshotList := []govultr.Snapshot{}

	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, ssh := range snapshots {
		sm, err := structToMap(ssh)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			snapshotList = append(snapshotList, ssh)
		}
	}

	if len(snapshotList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(snapshotList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(snapshotList[0].SnapshotID)
	d.Set("date_created", snapshotList[0].DateCreated)
	d.Set("description", snapshotList[0].Description)
	d.Set("size", snapshotList[0].Size)
	d.Set("status", snapshotList[0].Status)
	d.Set("os_id", snapshotList[0].OsID)
	d.Set("app_id", snapshotList[0].AppID)
	return nil
}
