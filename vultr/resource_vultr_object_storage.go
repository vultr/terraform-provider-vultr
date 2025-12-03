package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrObjectStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrObjectStorageCreate,
		ReadContext:   resourceVultrObjectStorageRead,
		UpdateContext: resourceVultrObjectStorageUpdate,
		DeleteContext: resourceVultrObjectStorageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"tier_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"bucket": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enable_versioning": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"enable_lock": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
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

func resourceVultrObjectStorageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	objReq := &govultr.ObjectStorageReq{
		ClusterID: d.Get("cluster_id").(int),
		TierID:    d.Get("tier_id").(int),
		Label:     d.Get("label").(string),
	}

	obj, _, err := client.ObjectStorage.Create(ctx, objReq)
	if err != nil {
		return diag.Errorf("error creating object storage: %v", err)
	}

	d.SetId(obj.ID)

	if _, err = waitForObjAvailable(ctx, d, "active", []string{"pending"}, "status", meta); err != nil {
		return diag.Errorf("error while waiting for Object Storage %s to be in a active state : %s", d.Id(), err)
	}

	if buckets, bucketsOK := d.GetOk("bucket"); bucketsOK {
		bucketList := buckets.([]interface{})
		for i := range bucketList {
			bucketObj := bucketList[i].(map[string]interface{})
			err := client.ObjectStorage.CreateBucket(ctx, d.Id(), &govultr.ObjectStorageBucketReq{
				Name:             bucketObj["name"].(string),
				EnableVersioning: bucketObj["enable_versioning"].(bool),
				EnableLock:       bucketObj["enable_lock"].(bool),
			})
			if err != nil {
				return diag.Errorf("error while creating object storage bucket : %s", err)
			}
		}
	}

	return resourceVultrObjectStorageRead(ctx, d, meta)
}

func resourceVultrObjectStorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	obj, _, err := client.ObjectStorage.Get(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting object storage account: %v", err)
	}

	if err := d.Set("date_created", obj.DateCreated); err != nil {
		return diag.Errorf("unable to set resource object_storage `date_created` read value: %v", err)
	}
	if err := d.Set("cluster_id", obj.ObjectStoreClusterID); err != nil {
		return diag.Errorf("unable to set resource object_storage `cluster_id` read value: %v", err)
	}
	if err := d.Set("label", obj.Label); err != nil {
		return diag.Errorf("unable to set resource object_storage `label` read value: %v", err)
	}
	if err := d.Set("location", obj.Location); err != nil {
		return diag.Errorf("unable to set resource object_storage `location` read value: %v", err)
	}
	if err := d.Set("region", obj.Region); err != nil {
		return diag.Errorf("unable to set resource object_storage `region` read value: %v", err)
	}
	if err := d.Set("status", obj.Status); err != nil {
		return diag.Errorf("unable to set resource object_storage `status` read value: %v", err)
	}
	if err := d.Set("s3_hostname", obj.S3Hostname); err != nil {
		return diag.Errorf("unable to set resource object_storage `s3_hostname` read value: %v", err)
	}
	if err := d.Set("s3_access_key", obj.S3AccessKey); err != nil {
		return diag.Errorf("unable to set resource object_storage `s3_access_key` read value: %v", err)
	}
	if err := d.Set("s3_secret_key", obj.S3SecretKey); err != nil {
		return diag.Errorf("unable to set resource object_storage `s3_secret_key` read value: %v", err)
	}

	return nil
}

func resourceVultrObjectStorageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	objReq := &govultr.ObjectStorageReq{Label: d.Get("label").(string)}

	if err := client.ObjectStorage.Update(ctx, d.Id(), objReq); err != nil {
		return diag.Errorf("error updating object storage %s label : %v", d.Id(), err)
	}

	if d.HasChange("bucket") {
		log.Printf("[INFO] Updating object storage buckets")

		oldBuckets, newBuckets := d.GetChange("bucket")
		oldBucketList := oldBuckets.([]interface{})
		newBucketList := newBuckets.([]interface{})

		var oldBucketNames []string
		var oldBucketObjs []map[string]interface{}
		for i := range oldBucketList {
			oldBucketObj := oldBucketList[i].(map[string]interface{})
			oldBucketObjs = append(oldBucketObjs, oldBucketObj)
			oldBucketNames = append(oldBucketNames, oldBucketObj["name"].(string))
		}

		var newBucketNames []string
		var newBucketObjs []map[string]interface{}
		for i := range newBucketList {
			newBucketObj := newBucketList[i].(map[string]interface{})
			newBucketObjs = append(newBucketObjs, newBucketObj)
			newBucketNames = append(newBucketNames, newBucketObj["name"].(string))
		}

		removeBucketNames := diffSlice(newBucketNames, oldBucketNames)
		createBucketNames := diffSlice(oldBucketNames, newBucketNames)

		// remove deleted buckets
		for i := range removeBucketNames {
			if err := client.ObjectStorage.DeleteBucket(ctx, d.Id(), removeBucketNames[i]); err != nil {
				return diag.Errorf("error deleting object storage bucket %q : %s", removeBucketNames[i], err)
			}
		}

		// handle create and updates
		for i := range newBucketObjs {

			// add new buckets
			for j := range createBucketNames {
				if newBucketObjs[i]["name"] == createBucketNames[j] {
					err := client.ObjectStorage.CreateBucket(ctx, d.Id(), &govultr.ObjectStorageBucketReq{
						Name:             newBucketObjs[i]["name"].(string),
						EnableVersioning: newBucketObjs[i]["enable_versioning"].(bool),
						EnableLock:       newBucketObjs[i]["enable_lock"].(bool),
					})
					if err != nil {
						return diag.Errorf("error adding object storage bucket %q : %s", createBucketNames[j], err)
					}
				}
			}

			// check against values of old state
			for k := range oldBucketObjs {
				if oldBucketObjs[k]["name"] == newBucketObjs[i]["name"] {
					recreate := false

					if oldBucketObjs[k]["enable_versioning"] != newBucketObjs[i]["enable_versioning"] {
						recreate = true
					}

					if oldBucketObjs[k]["enable_lock"] != newBucketObjs[i]["enable_lock"] {
						recreate = true
					}

					if recreate {
						if err := client.ObjectStorage.DeleteBucket(ctx, d.Id(), newBucketObjs[i]["name"].(string)); err != nil {
							return diag.Errorf(
								"error deleting (re-creating) object storage bucket %q : %s",
								newBucketObjs[i]["name"].(string),
								err,
							)
						}

						err := client.ObjectStorage.CreateBucket(ctx, d.Id(), &govultr.ObjectStorageBucketReq{
							Name:             newBucketObjs[i]["name"].(string),
							EnableVersioning: newBucketObjs[i]["enable_versioning"].(bool),
							EnableLock:       newBucketObjs[i]["enable_lock"].(bool),
						})
						if err != nil {
							return diag.Errorf(
								"error adding (re-creating) object storage bucket %q : %s",
								newBucketObjs[i]["name"].(string),
								err,
							)
						}
					}
				}
			}
		}
	}

	return resourceVultrObjectStorageRead(ctx, d, meta)
}

func resourceVultrObjectStorageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Object storage subscription %s", d.Id())

	if err := client.ObjectStorage.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting object storage subscription %s : %v", d.Id(), err)
	}

	return nil
}

func waitForObjAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll
	log.Printf(
		"[INFO] Waiting for Object Storage (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newServerObjRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newServerObjRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) retry.StateRefreshFunc { //nolint:lll
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Creating Object Storage")

		obj, _, err := client.ObjectStorage.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Object Store %s : %s", d.Id(), err)
		}

		log.Print(obj)
		if attr == "status" {
			log.Printf("[INFO] The Object Storage Status is %s", obj.Status)
			return obj, obj.Status, nil
		}
		return nil, "", nil
	}
}
