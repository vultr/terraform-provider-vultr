package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrDatabaseDB() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDatabaseDBCreate,
		ReadContext:   resourceVultrDatabaseDBRead,
		DeleteContext: resourceVultrDatabaseDBDelete,
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
		},
	}
}

func resourceVultrDatabaseDBCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	req := &govultr.DatabaseDBCreateReq{
		Name: d.Get("name").(string),
	}

	log.Printf("[INFO] Creating database logical DB")
	databaseDB, _, err := client.Database.CreateDB(ctx, databaseID, req)
	if err != nil {
		return diag.Errorf("error creating database logical DB: %v", err)
	}

	d.SetId(databaseDB.Name)

	return resourceVultrDatabaseDBRead(ctx, d, meta)
}

func resourceVultrDatabaseDBRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	databaseDB, _, err := client.Database.GetDB(ctx, databaseID, d.Id())
	if err != nil {
		return diag.Errorf("error getting database logical DB (%s): %v", d.Id(), err)
	}

	if err := d.Set("name", databaseDB.Name); err != nil {
		return diag.Errorf("unable to set resource database logical DB `name` read value: %v", err)
	}

	return nil
}

func resourceVultrDatabaseDBDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting database logical DB (%s)", d.Id())

	databaseID := d.Get("database_id").(string)

	if err := client.Database.DeleteDB(ctx, databaseID, d.Id()); err != nil {
		return diag.Errorf("error destroying database logical DB %s : %v", d.Id(), err)
	}

	return nil
}
