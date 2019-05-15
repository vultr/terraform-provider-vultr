package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func resourceVultrSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrSnapshotCreate,
		Read:   resourceVultrSnapshotRead,
		Delete: resourceVultrSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vps_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"date_created": {
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

func resourceVultrSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	vpsID := d.Get("vps_id").(string)

	var description string
	desc, ok := d.GetOk("description")
	if ok {
		description = desc.(string)
	}

	snapshot, err := client.Snapshot.Create(context.Background(), vpsID, description)
	if err != nil {
		return fmt.Errorf("Error creating snapshot: %v", err)
	}

	d.SetId(snapshot.SnapshotID)
	log.Printf("[INFO] Snapshot ID: %s", d.Id())

	return resourceVultrSnapshotRead(d, meta)
}

func resourceVultrSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	snapshots, err := client.Snapshot.GetList(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting snapshots: %v", err)
	}

	var snapshot *govultr.Snapshot
	for i := range snapshots {
		if snapshots[i].SnapshotID == d.Id() {
			snapshot = &snapshots[i]
			break
		}
	}

	if snapshot == nil {
		log.Printf("[WARN] Vultr snapshot (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("description", snapshot.Description)
	d.Set("date_created", snapshot.DateCreated)
	d.Set("size", snapshot.Size)
	d.Set("status", snapshot.Status)
	d.Set("os_id", snapshot.OsID)
	d.Set("app_id", snapshot.AppID)

	return nil
}

func resourceVultrSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Destroying snapshot: %s", d.Id())
	if err := client.Snapshot.Destroy(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("Error destroying snapshot (%s): %v", d.Id(), err)
	}

	return nil
}
