package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"date_created": {
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

func resourceVultrSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	req := &govultr.SnapshotReq{
		InstanceID:  d.Get("instance_id").(string),
		Description: d.Get("description").(string),
	}

	snapshot, err := client.Snapshot.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("error creating snapshot: %v", err)
	}

	d.SetId(snapshot.ID)

	if _, err = waitForSnapshot(d, "complete", []string{"pending"}, "status", meta); err != nil {
		return fmt.Errorf(
			"error while waiting for Snapshot %s to be completed: %s", d.Id(), err)
	}

	log.Printf("[INFO] Snapshot ID: %s", d.Id())

	return resourceVultrSnapshotRead(d, meta)
}

func resourceVultrSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	snapshot, err := client.Snapshot.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting snapshots: %v", err)
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

	log.Printf("[INFO] Deleting snapshot: %s", d.Id())
	if err := client.Snapshot.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error destroying snapshot (%s): %v", d.Id(), err)
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
			return nil, "", fmt.Errorf("error retrieving Snapshot %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] The SnapShot Status is %s", snap.Status)
		return snap, snap.Status, nil
	}
}
