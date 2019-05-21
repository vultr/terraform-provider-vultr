package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
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
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost_per_month": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_gb": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"attached_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVultrBlockStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	regionID := d.Get("region_id").(int)
	size := d.Get("size_gb").(int)

	var label string
	l, ok := d.GetOk("label")
	if ok {
		label = l.(string)
	}

	bs, err := client.BlockStorage.Create(context.Background(), regionID, size, label)
	if err != nil {
		return fmt.Errorf("Error creating block storage: %v", err)
	}

	d.SetId(bs.BlockStorageID)
	log.Printf("[INFO] Block Storage ID: %s", d.Id())

	return resourceVultrBlockStorageRead(d, meta)
}

func resourceVultrBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	bses, err := client.BlockStorage.GetList(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting block storage: %v", err)
	}

	var bs *govultr.BlockStorage
	for i := range bses {
		if bses[i].BlockStorageID == d.Id() {
			bs = &bses[i]
			break
		}
	}

	if bs == nil {
		log.Printf("[WARN] Vultr block storage (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("date_created", bs.DateCreated)
	d.Set("cost_per_month", bs.CostPerMonth)
	d.Set("status", bs.Status)
	d.Set("size_gb", bs.SizeGB)
	d.Set("region_id", bs.RegionID)
	d.Set("attached_id", bs.VpsID)
	d.Set("label", bs.Label)

	return nil
}

func resourceVultrBlockStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	d.Partial(true)

	if d.HasChange("label") {
		log.Printf(`[INFO] Updating block storage label (%s)`, d.Id())
		_, newVal := d.GetChange("label")
		err := client.BlockStorage.SetLabel(context.Background(), d.Id(), newVal.(string))
		if err != nil {
			return fmt.Errorf("Error setting block storage label (%s): %v", d.Id(), err)
		}
		d.SetPartial("label")
	}

	if d.HasChange("size_gb") {
		log.Printf(`[INFO] Resizing block storage (%s)`, d.Id())
		_, newVal := d.GetChange("size_gb")
		err := client.BlockStorage.Resize(context.Background(), d.Id(), newVal.(int))
		if err != nil {
			return fmt.Errorf("Error resizing block storage (%s): %v", d.Id(), err)
		}
		d.SetPartial("size_gb")
	}

	if d.HasChange("attached_id") {
		old, newVal := d.GetChange("attached_id")
		if old.(string) != "" {
			log.Printf(`[INFO] Detaching block storage (%s)`, d.Id())
			err := client.BlockStorage.Detach(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("Error detaching block storage (%s): %v", d.Id(), err)
			}
		}
		if newVal.(string) != "" {
			log.Printf(`[INFO] Attaching block storage (%s)`, d.Id())
			err := client.BlockStorage.Attach(context.Background(), d.Id(), newVal.(string))
			if err != nil {
				return fmt.Errorf("Error attaching block storage (%s): %v", d.Id(), err)
			}
		}
		d.SetPartial("attached_id")
	}

	d.Partial(false)

	return resourceVultrBlockStorageRead(d, meta)
}

func resourceVultrBlockStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting block storage: %s", d.Id())
	err := client.BlockStorage.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting block storage (%s): %v", d.Id(), err)
	}

	return nil
}
