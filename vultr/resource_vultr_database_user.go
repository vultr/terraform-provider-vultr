package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrDatabaseUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDatabaseUserCreate,
		ReadContext:   resourceVultrDatabaseUserRead,
		UpdateContext: resourceVultrDatabaseUserUpdate,
		DeleteContext: resourceVultrDatabaseUserDelete,
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
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encryption": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVultrDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	req := &govultr.DatabaseUserCreateReq{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		Encryption: d.Get("encryption").(string),
	}

	log.Printf("[INFO] Creating database user")
	databaseUser, _, err := client.Database.CreateUser(ctx, databaseID, req)
	if err != nil {
		return diag.Errorf("error creating database user: %v", err)
	}

	d.SetId(databaseUser.Username)

	return resourceVultrDatabaseUserRead(ctx, d, meta)
}

func resourceVultrDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	databaseUser, _, err := client.Database.GetUser(ctx, databaseID, d.Id())
	if err != nil {
		return diag.Errorf("error getting database user (%s): %v", d.Id(), err)
	}

	if err := d.Set("username", databaseUser.Username); err != nil {
		return diag.Errorf("unable to set resource database user `username` read value: %v", err)
	}

	if err := d.Set("password", databaseUser.Password); err != nil {
		return diag.Errorf("unable to set resource database user `password` read value: %v", err)
	}

	if databaseUser.Encryption != "" {
		var encryptionRaw string
		switch databaseUser.Encryption {
		case "Legacy (MySQL 5.x)":
			encryptionRaw = "mysql_native_password"
		default:
			encryptionRaw = "caching_sha2_password"
		}
		if err := d.Set("encryption", encryptionRaw); err != nil {
			return diag.Errorf("unable to set resource database user `encryption` read value: %v", err)
		}
	}

	return nil
}

func resourceVultrDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)

	if d.HasChange("password") {
		log.Printf("[INFO] Updating Password")
		_, newVal := d.GetChange("password")
		password := newVal.(string)
		req := &govultr.DatabaseUserUpdateReq{
			Password: password,
		}
		if _, _, err := client.Database.UpdateUser(ctx, databaseID, d.Id(), req); err != nil {
			return diag.Errorf("error updating database user %s : %s", d.Id(), err.Error())
		}
	}

	return resourceVultrDatabaseUserRead(ctx, d, meta)
}

func resourceVultrDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting database user (%s)", d.Id())

	databaseID := d.Get("database_id").(string)

	if err := client.Database.DeleteUser(ctx, databaseID, d.Id()); err != nil {
		return diag.Errorf("error destroying database user %s : %v", d.Id(), err)
	}

	return nil
}
