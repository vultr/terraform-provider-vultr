package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrBlockStorage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrBlockStorageRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attached_to_instance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrBlockStorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var blockList []govultr.BlockStorage
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		block, meta, _, err := client.BlockStorage.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting block storages: %v", err)
		}

		for _, b := range block {
			sm, err := structToMap(b)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				blockList = append(blockList, b)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}
	if len(blockList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(blockList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(blockList[0].ID)
	if err := d.Set("date_created", blockList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set block_storage `date_created` read value: %v", err)
	}
	if err := d.Set("cost", blockList[0].Cost); err != nil {
		return diag.Errorf("unable to set block_storage `cost` read value: %v", err)
	}
	if err := d.Set("status", blockList[0].Status); err != nil {
		return diag.Errorf("unable to set block_storage `status` read value: %v", err)
	}
	if err := d.Set("size_gb", blockList[0].SizeGB); err != nil {
		return diag.Errorf("unable to set block_storage `size_gb` read value: %v", err)
	}
	if err := d.Set("region", blockList[0].Region); err != nil {
		return diag.Errorf("unable to set block_storage `region` read value: %v", err)
	}
	if err := d.Set("attached_to_instance", blockList[0].AttachedToInstance); err != nil {
		return diag.Errorf("unable to set block_storage `attached_to_instance` read value: %v", err)
	}
	if err := d.Set("label", blockList[0].Label); err != nil {
		return diag.Errorf("unable to set block_storage `label` read value: %v", err)
	}
	if err := d.Set("mount_id", blockList[0].MountID); err != nil {
		return diag.Errorf("unable to set block_storage `mount_id` read value: %v", err)
	}
	if err := d.Set("block_type", blockList[0].BlockType); err != nil {
		return diag.Errorf("unable to set block_storage `block_type` read value: %v", err)
	}
	return nil
}
