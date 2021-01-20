package vultr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrSnapshotCreate,
		ReadContext:   resourceVultrSnapshotRead,
		DeleteContext: resourceVultrSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceVultrSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.SnapshotReq{
		InstanceID:  d.Get("instance_id").(string),
		Description: d.Get("description").(string),
	}

	snapshot, err := client.Snapshot.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating snapshot: %v", err)
	}

	d.SetId(snapshot.ID)

	if _, err = waitForSnapshot(ctx, d, "complete", []string{"pending"}, "status", meta); err != nil {
		return diag.Errorf(
			"error while waiting for Snapshot %s to be completed: %s", d.Id(), err)
	}

	log.Printf("[INFO] Snapshot ID: %s", d.Id())

	return resourceVultrSnapshotRead(ctx, d, meta)
}

func resourceVultrSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	snapshot, err := client.Snapshot.Get(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting snapshots: %v", err)
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

func resourceVultrSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting snapshot: %s", d.Id())
	if err := client.Snapshot.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying snapshot (%s): %v", d.Id(), err)
	}

	return nil
}

func waitForSnapshot(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
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

	return stateConf.WaitForStateContext(ctx)
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
