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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrBlockStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrBlockStorageCreate,
		ReadContext:   resourceVultrBlockStorageRead,
		UpdateContext: resourceVultrBlockStorageUpdate,
		DeleteContext: resourceVultrBlockStorageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"size_gb": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"region": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: IgnoreCase,
			},
			"attached_to_instance": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"live": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"storage_opt", "high_perf"}, false),
				ForceNew:     true,
			},
		},
	}
}

func resourceVultrBlockStorageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	bsReq := &govultr.BlockStorageCreate{
		Region:    d.Get("region").(string),
		SizeGB:    d.Get("size_gb").(int),
		Label:     d.Get("label").(string),
		BlockType: d.Get("block_type").(string),
	}

	bs, _, err := client.BlockStorage.Create(ctx, bsReq)
	if err != nil {
		return diag.Errorf("error creating block storage: %v", err)
	}

	d.SetId(bs.ID)
	log.Printf("[INFO] Block Storage ID: %s", d.Id())

	if _, err = waitForBlockAvailable(ctx, d, "active", []string{"pending"}, "status", meta); err != nil {
		return diag.Errorf("error while waiting for block %s to be completed: %s", d.Id(), err)
	}

	if instanceID, ok := d.GetOk("attached_to_instance"); ok {
		log.Printf("[INFO] Attaching block storage (%s)", d.Id())

		// Wait for the BS state to become active for 30 seconds
		bsReady := false
		for i := 0; i <= 30; i++ {
			bState, _, err := client.BlockStorage.Get(ctx, bs.ID)
			if err != nil {
				return diag.Errorf("error attaching: %s", err.Error())
			}
			if bState.Status == "active" {
				bsReady = true
				break
			}
			time.Sleep(1 * time.Second)
		}

		if !bsReady {
			return diag.Errorf("block storage was not in ready state after 30 seconds")
		}

		attachReq := &govultr.BlockStorageAttach{
			InstanceID: instanceID.(string),
			Live:       govultr.BoolToBoolPtr(d.Get("live").(bool)),
		}

		if err := client.BlockStorage.Attach(ctx, d.Id(), attachReq); err != nil {
			return diag.Errorf("error attaching block storage (%s): %v", d.Id(), err)
		}
	}

	return resourceVultrBlockStorageRead(ctx, d, meta)
}

func resourceVultrBlockStorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	bs, _, err := client.BlockStorage.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid block storage ID") {
			tflog.Warn(ctx, fmt.Sprintf("Removing block storage (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting block storage: %v", err)
	}

	if err := d.Set("live", d.Get("live").(bool)); err != nil {
		return diag.Errorf("unable to set resource block_storage `live` read value: %v", err)
	}
	if err := d.Set("date_created", bs.DateCreated); err != nil {
		return diag.Errorf("unable to set resource block_storage `date_created` read value: %v", err)
	}
	if err := d.Set("cost", bs.Cost); err != nil {
		return diag.Errorf("unable to set resource block_storage `cost` read value: %v", err)
	}
	if err := d.Set("status", bs.Status); err != nil {
		return diag.Errorf("unable to set resource block_storage `status` read value: %v", err)
	}
	if err := d.Set("size_gb", bs.SizeGB); err != nil {
		return diag.Errorf("unable to set resource block_storage `size_gb` read value: %v", err)
	}
	if err := d.Set("region", bs.Region); err != nil {
		return diag.Errorf("unable to set resource block_storage `region` read value: %v", err)
	}
	if err := d.Set("attached_to_instance", bs.AttachedToInstance); err != nil {
		return diag.Errorf("unable to set resource block_storage `attached_to_instance` read value: %v", err)
	}
	if err := d.Set("label", bs.Label); err != nil {
		return diag.Errorf("unable to set resource block_storage `label` read value: %v", err)
	}
	if err := d.Set("mount_id", bs.MountID); err != nil {
		return diag.Errorf("unable to set resource block_storage `mount_id` read value: %v", err)
	}
	if err := d.Set("block_type", bs.BlockType); err != nil {
		return diag.Errorf("unable to set resource block_storage `block_type` read value: %v", err)
	}

	return nil
}

func resourceVultrBlockStorageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	blockReq := &govultr.BlockStorageUpdate{}
	if d.HasChange("label") {
		blockReq.Label = d.Get("label").(string)
	}

	if d.HasChange("size_gb") {
		blockReq.SizeGB = d.Get("size_gb").(int)
	}

	if err := client.BlockStorage.Update(ctx, d.Id(), blockReq); err != nil {
		return diag.Errorf("error getting block storage: %v", err)
	}

	if d.HasChange("attached_to_instance") {
		old, newVal := d.GetChange("attached_to_instance")

		if old.(string) != "" {
			// The following check is necessary so we do not erroneously detach
			// after a formerly attached server has been tainted and/or
			// destroyed.
			bs, _, err := client.BlockStorage.Get(ctx, d.Id())
			if err != nil {
				return diag.Errorf("error getting block storage: %v", err)
			}

			if bs.AttachedToInstance != "" {
				log.Printf(`[INFO] Detaching block storage (%s)`, d.Id())

				blockReq := &govultr.BlockStorageDetach{Live: govultr.BoolToBoolPtr(d.Get("live").(bool))}
				err := client.BlockStorage.Detach(ctx, d.Id(), blockReq)
				if err != nil {
					return diag.Errorf("error detaching block storage (%s): %v", d.Id(), err)
				}
			}
		}

		if newVal.(string) != "" {
			log.Printf(`[INFO] Attaching block storage (%s)`, d.Id())
			blockReq := &govultr.BlockStorageAttach{
				InstanceID: newVal.(string),
				Live:       govultr.BoolToBoolPtr(d.Get("live").(bool)),
			}
			if err := client.BlockStorage.Attach(ctx, d.Id(), blockReq); err != nil {
				return diag.Errorf("error attaching block storage (%s): %v", d.Id(), err)
			}
		}
	}

	return resourceVultrBlockStorageRead(ctx, d, meta)
}

func resourceVultrBlockStorageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting block storage: %s", d.Id())
	if err := client.BlockStorage.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting block storage (%s): %v", d.Id(), err)
	}

	return nil
}

func waitForBlockAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll
	log.Printf(
		"[INFO] Waiting for Server (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newBlockStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newBlockStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) retry.StateRefreshFunc { //nolint:lll
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Creating Block")
		block, _, err := client.BlockStorage.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving block %s : %s", d.Id(), err)
		}

		if attr == "status" {
			log.Printf("[INFO] The Block Status is %s", block.Status)
			return block, block.Status, nil
		} else {
			return nil, "", nil
		}
	}
}
