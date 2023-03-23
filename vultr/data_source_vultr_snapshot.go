package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrSnapshotRead,
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
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var snapshotList []govultr.Snapshot
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		snapshots, meta,_, err := client.Snapshot.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting snapshots: %v", err)
		}

		for _, ssh := range snapshots {
			sm, err := structToMap(ssh)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				snapshotList = append(snapshotList, ssh)
			}
		}
		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(snapshotList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(snapshotList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(snapshotList[0].ID)
	if err := d.Set("date_created", snapshotList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set snapshot `date_created` read value: %v", err)
	}
	if err := d.Set("description", snapshotList[0].Description); err != nil {
		return diag.Errorf("unable to set snapshot `description` read value: %v", err)
	}
	if err := d.Set("size", snapshotList[0].Size); err != nil {
		return diag.Errorf("unable to set snapshot `size` read value: %v", err)
	}
	if err := d.Set("status", snapshotList[0].Status); err != nil {
		return diag.Errorf("unable to set snapshot `status` read value: %v", err)
	}
	if err := d.Set("os_id", snapshotList[0].OsID); err != nil {
		return diag.Errorf("unable to set snapshot `os_id` read value: %v", err)
	}
	if err := d.Set("app_id", snapshotList[0].AppID); err != nil {
		return diag.Errorf("unable to set snapshot `app_id` read value: %v", err)
	}
	return nil
}
