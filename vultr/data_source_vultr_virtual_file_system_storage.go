package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrVirtualFileSystemStorage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrVirtualFileSystemStorageRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"attached_instances": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"charges": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVultrVirtualFileSystemStorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var storageList []govultr.VirtualFileSystemStorage
	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	options := &govultr.ListOptions{}
	for {
		storages, meta, _, err := client.VirtualFileSystemStorage.List(ctx, options)
		if err != nil {
			return diag.Errorf("error listing virtual file system storages: %v", err)
		}

		for i := range storages {
			sm, err := structToMap(storages[i])
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				storageList = append(storageList, storages[i])
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(storageList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(storageList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(storageList[0].ID)
	if err := d.Set("region", storageList[0].Region); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `region` read value: %v", err)
	}
	if err := d.Set("size_gb", storageList[0].StorageSize.SizeGB); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `size_gb` read value: %v", err)
	}
	if err := d.Set("label", storageList[0].Label); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `label` read value: %v", err)
	}
	if err := d.Set("tags", storageList[0].Tags); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `tags` read value: %v", err)
	}
	if err := d.Set("date_created", storageList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `date_created` read value: %v", err)
	}
	if err := d.Set("status", storageList[0].Status); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `status` read value: %v", err)
	}
	if err := d.Set("size_gb", storageList[0].StorageSize.SizeGB); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `size_gb` read value: %v", err)
	}
	if err := d.Set("disk_type", storageList[0].DiskType); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `disk_type` read value: %v", err)
	}
	if err := d.Set("cost", storageList[0].Billing.Monthly); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `cost` read value: %v", err)
	}
	if err := d.Set("charges", storageList[0].Billing.Charges); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `charges` read value: %v", err)
	}

	attachments, _, err := client.VirtualFileSystemStorage.AttachmentList(ctx, d.Id())
	if err != nil {
		return diag.Errorf("unable to retrieve attachments for virtual file system storage %s", d.Id())
	}

	var attInstIDs []string
	var attStates []map[string]interface{}
	if len(attachments) != 0 {
		for i := range attachments {
			attInstIDs = append(attInstIDs, attachments[i].TargetID)
			attStates = append(attStates, map[string]interface{}{
				"instance_id": attachments[i].TargetID,
				"state":       attachments[i].State,
				"mount":       attachments[i].MountTag,
			})
		}
	}

	if err := d.Set("attached_instances", attInstIDs); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `attached_instances` read value: %v", err)
	}
	if err := d.Set("attachments", attStates); err != nil {
		return diag.Errorf("unable to set data source virtual_file_system_storage `attachments` read value: %v", err)
	}

	return nil
}
