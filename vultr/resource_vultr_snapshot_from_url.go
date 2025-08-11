package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vultr/govultr/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVultrSnapshotFromURL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrSnapshotFromURLCreate,
		ReadContext:   resourceVultrSnapshotRead,
		DeleteContext: resourceVultrSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_uefi": {
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceVultrSnapshotFromURLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	snapReq := &govultr.SnapshotURLReq{
		URL:  d.Get("url").(string),
		UEFI: govultr.BoolToBoolPtr(d.Get("use_uefi").(bool)),
	}

	snapshot, _, err := client.Snapshot.CreateFromURL(ctx, snapReq)
	if err != nil {
		return diag.Errorf("error creating snapshot: %v", err)
	}

	d.SetId(snapshot.ID)
	log.Printf("[INFO] Snapshot ID: %s", d.Id())

	return resourceVultrSnapshotRead(ctx, d, meta)
}
