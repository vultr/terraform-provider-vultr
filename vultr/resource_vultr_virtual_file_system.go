package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrVirtualFileSystemStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrVirtualFileSystemStorageCreate,
		ReadContext:   resourceVultrVirtualFileSystemStorageRead,
		UpdateContext: resourceVultrVirtualFileSystemStorageUpdate,
		DeleteContext: resourceVultrVirtualFileSystemStorageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: IgnoreCase,
			},
			"size_gb": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Default:  nil,
			},
			"attachments": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:     schema.TypeString,
							Required: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "nvme",
			},
			// computed fields
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
		},
	}
}

func resourceVultrVirtualFileSystemStorageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	req := govultr.VirtualFileSystemStorageReq{
		Region: d.Get("region").(string),
		Label:  d.Get("label").(string),
		StorageSize: govultr.VirtualFileSystemStorageSize{
			SizeGB: d.Get("size_gb").(int),
		},
	}

	storage, _, err := client.VirtualFileSystemStorage.Create(ctx, &req)
	if err != nil {
		return diag.Errorf("error creating virtual file system storage: %v", err)
	}

	d.SetId(storage.ID)
	log.Printf("[INFO] Virtual File System Storage ID: %s", d.Id())

	if _, err = waitForVirtualFileSystemStorageAvailable(ctx, d, "active", []string{"pending"}, "status", meta); err != nil { //nolint:lll
		return diag.Errorf("error while waiting for virtual file system storage %s to be completed: %s", d.Id(), err)
	}

	if attachments, ok := d.GetOk("attachments"); ok {
		att := attachments.([]interface{})
		for i := range att {
			elem := att[i].(map[string]interface{})
			log.Printf("[INFO] Attaching virtual file system storage %s to instance %s", d.Id(), elem["target"].(string))
			if _, _, err := client.VirtualFileSystemStorage.Attach(ctx, d.Id(), elem["target"].(string)); err != nil {
				return diag.Errorf("error attaching virtual file storage %s to instance %s", d.Id(), elem["target"].(string))
			}
		}
	}

	return resourceVultrVirtualFileSystemStorageRead(ctx, d, meta)
}

func resourceVultrVirtualFileSystemStorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	storage, _, err := client.VirtualFileSystemStorage.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Subscription ID Not Found.") {
			tflog.Warn(ctx, fmt.Sprintf("removing virtual file system storage (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting virtual file system storage: %v", err)
	}

	if err := d.Set("region", storage.Region); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `region` read value: %v", err)
	}
	if err := d.Set("size_gb", storage.StorageSize.SizeGB); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `size_gb` read value: %v", err)
	}
	if err := d.Set("label", storage.Label); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `label` read value: %v", err)
	}
	if err := d.Set("tags", storage.Tags); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `tags` read value: %v", err)
	}
	if err := d.Set("date_created", storage.DateCreated); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `date_created` read value: %v", err)
	}
	if err := d.Set("status", storage.Status); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `status` read value: %v", err)
	}
	if err := d.Set("size_gb", storage.StorageSize.SizeGB); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `size_gb` read value: %v", err)
	}
	if err := d.Set("disk_type", storage.DiskType); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `disk_type` read value: %v", err)
	}
	if err := d.Set("cost", storage.Billing.Monthly); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `cost` read value: %v", err)
	}
	if err := d.Set("charges", storage.Billing.Charges); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `charges` read value: %v", err)
	}

	attachments, _, err := client.VirtualFileSystemStorage.AttachmentList(ctx, d.Id())
	if err != nil {
		return diag.Errorf("unable to retrieve attachments for virtual file system storage %s", d.Id())
	}

	var attElems []map[string]interface{}
	for i := range attachments {
		attElems = append(attElems, map[string]interface{}{
			"target": attachments[i].TargetID,
			"state":  attachments[i].State,
			"mount":  attachments[i].MountTag,
		})
	}

	if err := d.Set("attachments", attElems); err != nil {
		return diag.Errorf("unable to set resource virtual_file_system_storage `attachments` read value: %v", err)
	}

	return nil
}

func resourceVultrVirtualFileSystemStorageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	updateReq := &govultr.VirtualFileSystemStorageUpdateReq{}
	if d.HasChange("label") {
		updateReq.Label = d.Get("label").(string)
	}

	if d.HasChange("size_gb") {
		updateReq.StorageSize.SizeGB = d.Get("size_gb").(int)
	}

	if _, _, err := client.VirtualFileSystemStorage.Update(ctx, d.Id(), updateReq); err != nil {
		return diag.Errorf("error updating virtual file system storage : %v", err)
	}

	if d.HasChange("attachments") {
		attOld, attNew := d.GetChange("attachments")
		elemsOld := attOld.(*schema.Set).List()
		elemsNew := attNew.(*schema.Set).List()

		var idOld, idNew, idDetach, idAttach []string
		for i := range elemsOld {
			idOld = append(idOld, elemsOld[i].(map[string]interface{})["target"].(string))
		}

		for i := range elemsNew {
			idNew = append(idNew, elemsNew[i].(map[string]interface{})["target"].(string))
		}

		idDetach = append(idDetach, diffSlice(idNew, idOld)...)
		idAttach = append(idAttach, diffSlice(idOld, idNew)...)

		for i := range idDetach {
			log.Printf(`[INFO] Detaching virtual file system storage %s from instance %s`, d.Id(), idDetach[i])
			if err := client.VirtualFileSystemStorage.Detach(ctx, d.Id(), idDetach[i]); err != nil {
				return diag.Errorf("error detaching instance %s from virtual file system storage %s", idDetach[i], d.Id())
			}
		}

		for i := range idAttach {
			log.Printf(`[INFO] Attaching virtual file system storage %s to instance %s`, d.Id(), idAttach[i])
			if _, _, err := client.VirtualFileSystemStorage.Attach(ctx, d.Id(), idAttach[i]); err != nil {
				return diag.Errorf("error attaching instance %s to virtual file system storage %s", idAttach[i], d.Id())
			}
		}
	}

	return resourceVultrVirtualFileSystemStorageRead(ctx, d, meta)
}

func resourceVultrVirtualFileSystemStorageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting virtual file system storage: %s", d.Id())
	if err := client.VirtualFileSystemStorage.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting virtual file system storage (%s): %v", d.Id(), err)
	}

	return nil
}

func waitForVirtualFileSystemStorageAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll
	log.Printf(
		"[INFO] Waiting for virtual file system storage (%s) to have %s of %s",
		d.Id(),
		attribute,
		target,
	)

	stateConf := &retry.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newVirtualFileSystemStorageStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newVirtualFileSystemStorageStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) retry.StateRefreshFunc { //nolint:lll
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Checking new virtual file system storage")
		storage, _, err := client.VirtualFileSystemStorage.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving virtual file system storage %s : %s", d.Id(), err)
		}

		if attr == "status" {
			log.Printf("[INFO] The virtual file system storage status is %s", storage.Status)
			return storage, storage.Status, nil
		} else {
			return nil, "", nil
		}
	}
}
