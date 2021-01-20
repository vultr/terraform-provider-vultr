package vultr

import (
	"context"
	"fmt"
	"github.com/vultr/govultr/v2"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVultrSnapshotFromURL() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrSnapshotFromURLCreate,
		Read:   resourceVultrSnapshotRead,
		Delete: resourceVultrSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceVultrSnapshotFromURLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	snapReq := &govultr.SnapshotURLReq{
		URL: d.Get("url").(string),
	}

	snapshot, err := client.Snapshot.CreateFromURL(context.Background(), snapReq)
	if err != nil {
		return fmt.Errorf("error creating snapshot: %v", err)
	}

	d.SetId(snapshot.ID)
	log.Printf("[INFO] Snapshot ID: %s", d.Id())

	return resourceVultrSnapshotRead(d, meta)
}
