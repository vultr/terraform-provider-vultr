package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrSSHKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrSSHKeyCreate,
		ReadContext:   resourceVultrSSHKeyRead,
		UpdateContext: resourceVultrSSHKeyUpdate,
		DeleteContext: resourceVultrSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ssh_key": {
				Type:     schema.TypeString,
				Required: true,
			},

			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	sshReq := &govultr.SSHKeyReq{
		Name:   d.Get("name").(string),
		SSHKey: d.Get("ssh_key").(string),
	}

	key, err := client.SSHKey.Create(ctx, sshReq)
	if err != nil {
		return diag.Errorf("error creating SSH key: %v", err)
	}

	d.SetId(key.ID)
	log.Printf("[INFO] SSH Key ID: %s", d.Id())

	return resourceVultrSSHKeyRead(ctx, d, meta)
}

func resourceVultrSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	key, err := client.SSHKey.Get(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting SSH keys: %v", err)
	}

	d.Set("name", key.Name)
	d.Set("ssh_key", key.SSHKey)
	d.Set("date_created", key.DateCreated)

	return nil
}

func resourceVultrSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	key := &govultr.SSHKeyReq{}

	if d.HasChange("name") {
		key.Name = d.Get("name").(string)
	}

	if d.HasChange("ssh_key") {
		key.SSHKey = d.Get("ssh_key").(string)
	}

	log.Printf("[INFO] Updating SSH Key: %s", d.Id())
	if err := client.SSHKey.Update(context.Background(), d.Id(), key); err != nil {
		return diag.Errorf("error updating SSH key (%s): %v", d.Id(), err)
	}

	return resourceVultrSSHKeyRead(ctx, d, meta)
}

func resourceVultrSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting SSH Key: %s", d.Id())

	if err := client.SSHKey.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying SSH key (%s): %v", d.Id(), err)
	}

	return nil
}
