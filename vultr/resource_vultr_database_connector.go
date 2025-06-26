package vultr

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrDatabaseConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDatabaseConnectorCreate,
		ReadContext:   resourceVultrDatabaseConnectorRead,
		UpdateContext: resourceVultrDatabaseConnectorUpdate,
		DeleteContext: resourceVultrDatabaseConnectorDelete,
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
				ForceNew: true,
			},
			"class": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topics": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVultrDatabaseConnectorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)
	config := d.Get("config").(string)

	var configMap map[string]interface{}
	if config != "" {
		if err := json.Unmarshal([]byte(config), &configMap); err != nil {
			return diag.Errorf("error parsing JSON for field `config` for database connector create: %v", err)
		}
	}

	req := &govultr.DatabaseConnectorCreateReq{
		Name:   d.Get("name").(string),
		Class:  d.Get("partitions").(string),
		Topics: d.Get("replication").(string),
		Config: configMap,
	}

	log.Printf("[INFO] Creating database connector")
	databaseConnector, _, err := client.Database.CreateConnector(ctx, databaseID, req)
	if err != nil {
		return diag.Errorf("error creating database connector: %v", err)
	}

	d.SetId(databaseConnector.Name)

	return resourceVultrDatabaseConnectorRead(ctx, d, meta)
}

func resourceVultrDatabaseConnectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	databaseConnector, _, err := client.Database.GetConnector(ctx, databaseID, d.Id())
	if err != nil {
		return diag.Errorf("error getting database connector (%s): %v", d.Id(), err)
	}

	if err := d.Set("name", databaseConnector.Name); err != nil {
		return diag.Errorf("unable to set resource database connector `name` read value: %v", err)
	}

	if err := d.Set("class", databaseConnector.Class); err != nil {
		return diag.Errorf("unable to set resource database connector `class` read value: %v", err)
	}

	if err := d.Set("topics", databaseConnector.Topics); err != nil {
		return diag.Errorf("unable to set resource database connector `topics` read value: %v", err)
	}

	// handle config json to map

	return nil
}

func resourceVultrDatabaseConnectorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	req := &govultr.DatabaseConnectorUpdateReq{}
	log.Printf("[INFO] Updating database connector (%s)", d.Id())

	if d.HasChange("topics") {
		log.Print("[INFO] Updating `topics`")
		req.Topics = d.Get("topics").(string)
	}

	if d.HasChange("config") {
		log.Print("[INFO] Updating `config`")
		config := d.Get("config").(string)

		var configMap map[string]interface{}
		if config != "" {
			if err := json.Unmarshal([]byte(config), &configMap); err != nil {
				return diag.Errorf("error parsing JSON for field `config` for updating database connector %s : %s", d.Id(), err.Error())
			}
		}

		req.Config = configMap
	}

	if _, _, err := client.Database.UpdateConnector(ctx, databaseID, d.Id(), req); err != nil {
		return diag.Errorf("error updating database connector %s : %s", d.Id(), err.Error())
	}

	return resourceVultrDatabaseConnectorRead(ctx, d, meta)
}

func resourceVultrDatabaseConnectorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting database connector (%s)", d.Id())

	databaseID := d.Get("database_id").(string)

	if err := client.Database.DeleteConnector(ctx, databaseID, d.Id()); err != nil {
		return diag.Errorf("error destroying database connector %s : %v", d.Id(), err)
	}

	return nil
}
