package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrDatabaseTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDatabaseTopicCreate,
		ReadContext:   resourceVultrDatabaseTopicRead,
		UpdateContext: resourceVultrDatabaseTopicUpdate,
		DeleteContext: resourceVultrDatabaseTopicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"database_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"partitions": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"replication": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"retention_hours": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"retention_bytes": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceVultrDatabaseTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	req := &govultr.DatabaseTopicCreateReq{
		Name:           d.Get("name").(string),
		Partitions:     d.Get("partitions").(int),
		Replication:    d.Get("replication").(int),
		RetentionHours: d.Get("retention_hours").(int),
		RetentionBytes: d.Get("retention_bytes").(int),
	}

	log.Printf("[INFO] Creating database topic")
	databaseTopic, _, err := client.Database.CreateTopic(ctx, databaseID, req)
	if err != nil {
		return diag.Errorf("error creating database topic: %v", err)
	}

	d.SetId(databaseTopic.Name)

	return resourceVultrDatabaseTopicRead(ctx, d, meta)
}

func resourceVultrDatabaseTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	databaseTopic, _, err := client.Database.GetTopic(ctx, databaseID, d.Id())
	if err != nil {
		return diag.Errorf("error getting database topic (%s): %v", d.Id(), err)
	}

	if err := d.Set("name", databaseTopic.Name); err != nil {
		return diag.Errorf("unable to set resource database topic `name` read value: %v", err)
	}

	if err := d.Set("partitions", databaseTopic.Partitions); err != nil {
		return diag.Errorf("unable to set resource database topic `partitions` read value: %v", err)
	}

	if err := d.Set("replication", databaseTopic.Replication); err != nil {
		return diag.Errorf("unable to set resource database topic `replication` read value: %v", err)
	}

	if err := d.Set("retention_hours", databaseTopic.RetentionHours); err != nil {
		return diag.Errorf("unable to set resource database topic `retention_hours` read value: %v", err)
	}

	if err := d.Set("retention_bytes", databaseTopic.RetentionBytes); err != nil {
		return diag.Errorf("unable to set resource database topic `retention_bytes` read value: %v", err)
	}

	return nil
}

func resourceVultrDatabaseTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	req := &govultr.DatabaseTopicUpdateReq{}
	log.Printf("[INFO] Updating database topic (%s)", d.Id())

	if d.HasChange("partitions") {
		log.Print("[INFO] Updating `partitions`")
		req.Partitions = d.Get("partitions").(int)
	}

	if d.HasChange("replication") {
		log.Print("[INFO] Updating `replication`")
		req.Replication = d.Get("replication").(int)
	}

	if d.HasChange("retention_hours") {
		log.Print("[INFO] Updating `retention_hours`")
		req.RetentionHours = d.Get("retention_hours").(int)
	}

	if d.HasChange("retention_bytes") {
		log.Print("[INFO] Updating `retention_bytes`")
		req.RetentionBytes = d.Get("retention_bytes").(int)
	}

	if _, _, err := client.Database.UpdateTopic(ctx, databaseID, d.Id(), req); err != nil {
		return diag.Errorf("error updating database topic %s : %s", d.Id(), err.Error())
	}

	return resourceVultrDatabaseTopicRead(ctx, d, meta)
}

func resourceVultrDatabaseTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting database topic (%s)", d.Id())

	databaseID := d.Get("database_id").(string)

	if err := client.Database.DeleteTopic(ctx, databaseID, d.Id()); err != nil {
		return diag.Errorf("error destroying database topic %s : %v", d.Id(), err)
	}

	return nil
}
