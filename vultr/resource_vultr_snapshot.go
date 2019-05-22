package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
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
	_, err = waitForSnapshot(d, "complete", []string{"pending"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error while waiting for Snapshot %s to be completed: %s", d.Id(), err)
	}

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

func waitForSnapshot(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for Snapshot (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newSnapStateRefresh(d, meta),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForState()
}

func newSnapStateRefresh(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating Snapshot")
		snap, err := client.Snapshot.Get(context.Background(), d.Id())

		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving Snapshot %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] The SnapShot Status is %s", snap.Status)
		return snap, snap.Status, nil
	}
}
