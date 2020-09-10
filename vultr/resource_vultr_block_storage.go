package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrBlockStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrBlockStorageCreate,
		Read:   resourceVultrBlockStorageRead,
		Update: resourceVultrBlockStorageUpdate,
		Delete: resourceVultrBlockStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"size_gb": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrBlockStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	bsReq := &govultr.BlockStorageCreate{
		Region: d.Get("region").(string),
		SizeGB: d.Get("size_gb").(int),
		Label:  d.Get("label").(string),
	}

	bs, err := client.BlockStorage.Create(context.Background(), bsReq)
	if err != nil {
		return fmt.Errorf("error creating block storage: %v", err)
	}

	d.SetId(bs.ID)
	log.Printf("[INFO] Block Storage ID: %s", d.Id())

	if instanceID, ok := d.GetOkExists("attached_to_instance"); ok {
		log.Printf("[INFO] Attaching block storage (%s)", d.Id())

		// Wait for the BS state to become active for 30 seconds
		bsReady := false
		for i := 0; i <= 30; i++ {
			bState, err := client.BlockStorage.Get(context.Background(), bs.ID)
			if err != nil {
				return fmt.Errorf("error attaching: %s", err.Error())
			}
			if bState.Status == "active" {
				bsReady = true
				break
			}
			time.Sleep(1 * time.Second)
		}

		if !bsReady {
			return fmt.Errorf("blockstorage was not in ready state after 30 seconds")
		}

		attachReq := &govultr.BlockStorageAttach{
			InstanceID: instanceID.(string),
			Live:       d.Get("live").(bool),
		}

		if err := client.BlockStorage.Attach(context.Background(), d.Id(), attachReq); err != nil {
			return fmt.Errorf("error attaching block storage (%s): %v", d.Id(), err)
		}
	}

	return resourceVultrBlockStorageRead(d, meta)
}

func resourceVultrBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	bs, err := client.BlockStorage.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting block storage: %v", err)
	}

	d.Set("live", d.Get("live").(bool))
	d.Set("date_created", bs.DateCreated)
	d.Set("cost", bs.Cost)
	d.Set("status", bs.Status)
	d.Set("size_gb", bs.SizeGB)
	d.Set("region", bs.Region)
	d.Set("attached_to_instance", bs.AttachedToInstance)
	d.Set("label", bs.Label)

	return nil
}

func resourceVultrBlockStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	blockReq := &govultr.BlockStorageUpdate{}
	if d.HasChange("label") {
		blockReq.Label = d.Get("label").(string)
	}

	if d.HasChange("size_gb") {
		blockReq.SizeGB = d.Get("size_gb").(int)
	}

	if err := client.BlockStorage.Update(context.Background(), d.Id(), blockReq); err != nil {
		return fmt.Errorf("error getting block storage: %v", err)
	}

	if d.HasChange("attached_to_instance") {
		old, newVal := d.GetChange("attached_to_instance")
		live := d.Get("live").(bool)

		if old.(string) != "" {
			// The following check is necessary so we do not erroneously detach after a formerly attached server has been tainted and/or destroyed.
			bs, err := client.BlockStorage.Get(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("error getting block storage: %v", err)
			}

			if bs.AttachedToInstance != "" {
				log.Printf(`[INFO] Detaching block storage (%s)`, d.Id())

				blockReq := &govultr.BlockStorageDetach{Live: live}
				err := client.BlockStorage.Detach(context.Background(), d.Id(), blockReq)
				if err != nil {
					return fmt.Errorf("error detaching block storage (%s): %v", d.Id(), err)
				}
			}
		}

		if newVal.(string) != "" {
			log.Printf(`[INFO] Attaching block storage (%s)`, d.Id())
			blockReq := &govultr.BlockStorageAttach{
				InstanceID: newVal.(string),
				Live:       live,
			}
			if err := client.BlockStorage.Attach(context.Background(), d.Id(), blockReq); err != nil {
				return fmt.Errorf("error attaching block storage (%s): %v", d.Id(), err)
			}
		}
	}

	return resourceVultrBlockStorageRead(d, meta)
}

func resourceVultrBlockStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting block storage: %s", d.Id())
	if err := client.BlockStorage.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error deleting block storage (%s): %v", d.Id(), err)
	}

	return nil
}
