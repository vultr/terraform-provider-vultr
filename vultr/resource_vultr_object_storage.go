package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
			"object_storage_cluster_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"location": {
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_secret_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrObjectStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	objStoreCluster := d.Get("object_storage_cluster_id").(int)

	// optional param
	label := d.Get("label").(string)

	obj, err := client.ObjectStorage.Create(context.Background(), objStoreCluster, label)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	d.SetId(strconv.Itoa(obj.ID))

	_, err = waitForObjAvailable(d, "active", []string{"pending"}, "status", meta)
	if err != nil {
		return fmt.Errorf("Error while waiting for Object Storage %s to be in a active state : %s", d.Id(), err)
	}

	return resourceVultrObjectStorageRead(d, meta)
}

func resourceVultrObjectStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	obj, err := client.ObjectStorage.Get(context.Background(), id)
	if err != nil {
		return fmt.Errorf("error getting object storage account: %v", err)
	}

	d.Set("date_created", obj.DateCreated)
	d.Set("object_storage_cluster_id", obj.ObjectStoreClusterID)
	d.Set("label", obj.Label)
	d.Set("region_id", obj.RegionID)
	d.Set("location", obj.Location)
	d.Set("status", obj.Status)
	d.Set("s3_hostname", obj.S3Hostname)
	d.Set("s3_access_key", obj.S3AccessKey)
	d.Set("s3_secret_key", obj.S3SecretKey)
	return nil
}

func resourceVultrObjectStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	label := d.Get("label").(string)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = client.ObjectStorage.SetLabel(context.Background(), id, label)
	if err != nil {
		return fmt.Errorf("error updating object storage %d label : %v", id, err)
	}

	return resourceVultrObjectStorageRead(d, meta)
}

func resourceVultrObjectStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Object storage subscription %s", d.Id())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("error with subsc")
	}

	err = client.ObjectStorage.Delete(context.Background(), id)
	if err != nil {
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
		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, "", err
		}

		obj, err := client.ObjectStorage.Get(context.Background(), id)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving Object Store %s : %s", d.Id(), err)
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
