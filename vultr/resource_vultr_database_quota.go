package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrDatabaseQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDatabaseQuotaCreate,
		ReadContext:   resourceVultrDatabaseQuotaRead,
		UpdateContext: resourceVultrDatabaseQuotaUpdate,
		DeleteContext: resourceVultrDatabaseQuotaDelete,
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
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumer_byte_rate": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"producer_byte_rate": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"request_percentage": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrDatabaseQuotaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)
	clientID := d.Get("client_id").(string)
	user := d.Get("user").(string)

	// Check quota list for duplicate client ID/user combination
	log.Printf("[INFO] Fetching database quota list")
	quotas, _, _, err := client.Database.ListQuotas(ctx, databaseID)
	if err != nil {
		return diag.Errorf("error creating database quota: %v", err)
	}

	// Return error if any combination matches
	for _, quota := range quotas {
		if quota.ClientID == clientID && quota.User == user {
			return diag.Errorf("error creating database quota: a quota with this client ID and user already exists")
		}
	}

	// Good to create now
	req := &govultr.DatabaseQuotaCreateReq{
		ClientID:          clientID,
		ConsumerByteRate:  d.Get("consumer_byte_rate").(int),
		ProducerByteRate:  d.Get("producer_byte_rate").(int),
		RequestPercentage: d.Get("request_percentage").(int),
		User:              user,
	}

	log.Printf("[INFO] Creating database quota")
	DatabaseQuota, _, err := client.Database.CreateQuota(ctx, databaseID, req)
	if err != nil {
		return diag.Errorf("error creating database quota: %v", err)
	}

	d.SetId(fmt.Sprintf("%s|%s", DatabaseQuota.ClientID, DatabaseQuota.User))

	return resourceVultrDatabaseQuotaRead(ctx, d, meta)
}

func resourceVultrDatabaseQuotaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)
	quotaID := strings.Split(d.Id(), "|")

	DatabaseQuota, _, err := client.Database.GetQuota(ctx, databaseID, quotaID[0], quotaID[1])
	if err != nil {
		return diag.Errorf("error getting database quota (%s): %v", d.Id(), err)
	}

	if err := d.Set("client_id", DatabaseQuota.ClientID); err != nil {
		return diag.Errorf("unable to set resource database quota `client_id` read value: %v", err)
	}

	if err := d.Set("consumer_byte_rate", DatabaseQuota.ConsumerByteRate); err != nil {
		return diag.Errorf("unable to set resource database quota `consumer_byte_rate` read value: %v", err)
	}

	if err := d.Set("producer_byte_rate", DatabaseQuota.ProducerByteRate); err != nil {
		return diag.Errorf("unable to set resource database quota `producer_byte_rate` read value: %v", err)
	}

	if err := d.Set("request_percentage", DatabaseQuota.RequestPercentage); err != nil {
		return diag.Errorf("unable to set resource database quota `request_percentage` read value: %v", err)
	}

	if err := d.Set("user", DatabaseQuota.User); err != nil {
		return diag.Errorf("unable to set resource database quota `user` read value: %v", err)
	}

	return nil
}

func resourceVultrDatabaseQuotaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	databaseID := d.Get("database_id").(string)
	quotaID := strings.Split(d.Id(), "|")

	req := &govultr.DatabaseQuotaUpdateReq{}
	log.Printf("[INFO] Updating database quota (%s)", d.Id())

	if d.HasChange("consumer_byte_rate") {
		log.Print("[INFO] Updating `consumer_byte_rate`")
		req.ConsumerByteRate = d.Get("consumer_byte_rate").(int)
	}

	if d.HasChange("producer_byte_rate") {
		log.Print("[INFO] Updating `producer_byte_rate`")
		req.ProducerByteRate = d.Get("producer_byte_rate").(int)
	}

	if d.HasChange("request_percentage") {
		log.Print("[INFO] Updating `request_percentage`")
		req.RequestPercentage = d.Get("request_percentage").(int)
	}

	if _, _, err := client.Database.UpdateQuota(ctx, databaseID, quotaID[0], quotaID[1], req); err != nil {
		return diag.Errorf("error updating database quota %s : %s", d.Id(), err.Error())
	}

	return resourceVultrDatabaseQuotaRead(ctx, d, meta)
}

func resourceVultrDatabaseQuotaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting database quota (%s)", d.Id())

	databaseID := d.Get("database_id").(string)
	quotaID := strings.Split(d.Id(), "|")

	if err := client.Database.DeleteQuota(ctx, databaseID, quotaID[0], quotaID[1]); err != nil {
		return diag.Errorf("error destroying database quota %s : %v", d.Id(), err)
	}

	return nil
}
