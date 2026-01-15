package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrObjectStorageBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrObjectStorageBucketCreate,
		ReadContext:   resourceVultrObjectStorageBucketRead,
		DeleteContext: resourceVultrObjectStorageBucketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"object_storage_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"enable_versioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"enable_lock": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrObjectStorageBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf(
		"[INFO] Creating object storage bucket %s in object storage %s",
		d.Get("name").(string),
		d.Get("object_storage_id").(string),
	)

	err := client.ObjectStorage.CreateBucket(ctx, d.Get("object_storage_id").(string), &govultr.ObjectStorageBucketReq{
		Name:             d.Get("name").(string),
		EnableVersioning: d.Get("enable_versioning").(bool),
		EnableLock:       d.Get("enable_lock").(bool),
	})
	if err != nil {
		return diag.Errorf("error while creating object storage bucket : %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", d.Get("object_storage_id").(string), d.Get("name").(string)))

	return resourceVultrObjectStorageBucketRead(ctx, d, meta)
}

func resourceVultrObjectStorageBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	buckets, _, err := client.ObjectStorage.ListBuckets(ctx, d.Get("object_storage_id").(string))
	if err != nil {
		return diag.Errorf("error getting object storage buckets: %v", err)
	}

	found := false
	for i := range buckets {
		if buckets[i].Name == d.Get("name").(string) {
			found = true
			if err := d.Set("name", buckets[i].Name); err != nil {
				return diag.Errorf("unable to set resource object_storage_bucket `name` read value: %v", err)
			}
			if err := d.Set("date_created", buckets[i].DateCreated); err != nil {
				return diag.Errorf("unable to set resource object_storage_bucket `date_created` read value: %v", err)
			}
		}
	}

	if !found {
		return diag.Errorf("object storage bucket (%v) not found in object storage (%v)", d.Get("name").(string), d.Id())
	}

	return nil
}

func resourceVultrObjectStorageBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf(
		"[INFO] Deleting object storage bucket %s from object storage %s",
		d.Get("name").(string),
		d.Get("object_storage_id").(string),
	)

	if err := client.ObjectStorage.DeleteBucket(
		ctx,
		d.Get("object_storage_id").(string),
		d.Get("name").(string),
	); err != nil {
		return diag.Errorf("error deleting object storage bucket %s : %v", d.Get("name").(string), err)
	}

	return nil
}
