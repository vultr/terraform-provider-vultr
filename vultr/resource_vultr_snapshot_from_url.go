package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
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

	url := d.Get("url").(string)

	snapshot, err := client.Snapshot.CreateFromURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error creating snapshot: %v", err)
	}

	d.SetId(snapshot.SnapshotID)
	log.Printf("[INFO] Snapshot ID: %s", d.Id())

	return resourceVultrSnapshotRead(d, meta)
}
