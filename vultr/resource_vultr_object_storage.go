package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVultrObjectStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrObjectStorageCreate,
		Read:   resourceVultrObjectStorageRead,
		Update: resourceVultrObjectStorageUpdate,
		Delete: resourceVultrObjectStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"s3_secret_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceVultrObjectStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	objStoreCluster := d.Get("cluster_id").(int)
	label := d.Get("label").(string)

	obj, err := client.ObjectStorage.Create(context.Background(), objStoreCluster, label)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	d.SetId(obj.ID)

	if _, err = waitForObjAvailable(d, "active", []string{"pending"}, "status", meta); err != nil {
		return fmt.Errorf("error while waiting for Object Storage %s to be in a active state : %s", d.Id(), err)
	}

	return resourceVultrObjectStorageRead(d, meta)
}

func resourceVultrObjectStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	obj, err := client.ObjectStorage.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting object storage account: %v", err)
	}

	d.Set("date_created", obj.DateCreated)
	d.Set("cluster_id", obj.ObjectStoreClusterID)
	d.Set("label", obj.Label)
	d.Set("region", obj.Region)
	d.Set("status", obj.Status)
	d.Set("s3_hostname", obj.S3Hostname)
	d.Set("s3_access_key", obj.S3AccessKey)
	d.Set("s3_secret_key", obj.S3SecretKey)

	return nil
}

func resourceVultrObjectStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	label := d.Get("label").(string)

	if err := client.ObjectStorage.Update(context.Background(), d.Id(), label); err != nil {
		return fmt.Errorf("error updating object storage %s label : %v", d.Id(), err)
	}

	return resourceVultrObjectStorageRead(d, meta)
}

func resourceVultrObjectStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Object storage subscription %s", d.Id())

	if err := client.ObjectStorage.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error deleting object storage subscription %s : %v", d.Id(), err)
	}

	return nil
}

func waitForObjAvailable(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for Object Storage (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newServerObjRefresh(d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForState()
}

func newServerObjRefresh(d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Creating Object Storage")

		obj, err := client.ObjectStorage.Get(context.Background(), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Object Store %s : %s", d.Id(), err)
		}

		log.Print(obj)
		if attr == "status" {
			log.Printf("[INFO] The Object Storage Status is %s", obj.Status)
			return obj, obj.Status, nil
		} else {
			return nil, "", nil
		}
	}
}
