package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrDatabaseConnectionPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDatabaseConnectionPoolCreate,
		ReadContext:   resourceVultrDatabaseConnectionPoolRead,
		UpdateContext: resourceVultrDatabaseConnectionPoolUpdate,
		DeleteContext: resourceVultrDatabaseConnectionPoolDelete,
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
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceVultrDatabaseConnectionPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	req := &govultr.DatabaseConnectionPoolCreateReq{
		Name:     d.Get("name").(string),
		Database: d.Get("database").(string),
		Username: d.Get("username").(string),
		Mode:     d.Get("mode").(string),
		Size:     d.Get("size").(int),
	}

	log.Printf("[INFO] Creating database connection pool")
	databaseConnectionPool, _, err := client.Database.CreateConnectionPool(ctx, databaseID, req)
	if err != nil {
		return diag.Errorf("error creating database connection pool: %v", err)
	}

	d.SetId(databaseConnectionPool.Name)

	return resourceVultrDatabaseConnectionPoolRead(ctx, d, meta)
}

func resourceVultrDatabaseConnectionPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	databaseConnectionPool, _, err := client.Database.GetConnectionPool(ctx, databaseID, d.Id())
	if err != nil {
		return diag.Errorf("error getting database connection pool (%s): %v", d.Id(), err)
	}

	if err := d.Set("name", databaseConnectionPool.Name); err != nil {
		return diag.Errorf("unable to set resource database connection pool `name` read value: %v", err)
	}

	if err := d.Set("database", databaseConnectionPool.Database); err != nil {
		return diag.Errorf("unable to set resource database connection pool `database` read value: %v", err)
	}

	if err := d.Set("username", databaseConnectionPool.Username); err != nil {
		return diag.Errorf("unable to set resource database connection pool `username` read value: %v", err)
	}

	if err := d.Set("mode", databaseConnectionPool.Mode); err != nil {
		return diag.Errorf("unable to set resource database connection pool `mode` read value: %v", err)
	}

	if err := d.Set("size", databaseConnectionPool.Size); err != nil {
		return diag.Errorf("unable to set resource database connection pool `size` read value: %v", err)
	}

	return nil
}

func resourceVultrDatabaseConnectionPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	req := &govultr.DatabaseConnectionPoolUpdateReq{}

	if d.HasChange("database") {
		log.Printf("[INFO] Updating Pool Database")
		_, newVal := d.GetChange("database")
		database := newVal.(string)
		req.Database = database
	}

	if d.HasChange("username") {
		log.Printf("[INFO] Updating Pool User")
		_, newVal := d.GetChange("username")
		username := newVal.(string)
		req.Username = username
	}

	if d.HasChange("mode") {
		log.Printf("[INFO] Updating Pool Mode")
		_, newVal := d.GetChange("mode")
		mode := newVal.(string)
		req.Mode = mode
	}

	if d.HasChange("size") {
		log.Printf("[INFO] Updating Pool Size")
		_, newVal := d.GetChange("size")
		size := newVal.(int)
		req.Size = size
	}

	// Only update if we've actually passed some new values
	if req != (&govultr.DatabaseConnectionPoolUpdateReq{}) {
		if _, _, err := client.Database.UpdateConnectionPool(ctx, databaseID, d.Id(), req); err != nil {
			return diag.Errorf("error updating database connection pool %s : %s", d.Id(), err.Error())
		}
	}

	return resourceVultrDatabaseConnectionPoolRead(ctx, d, meta)
}

func resourceVultrDatabaseConnectionPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting database connection pool (%s)", d.Id())

	databaseID := d.Get("database_id").(string)

	if err := client.Database.DeleteConnectionPool(ctx, databaseID, d.Id()); err != nil {
		return diag.Errorf("error destroying database connection pool %s : %v", d.Id(), err)
	}

	return nil
}
