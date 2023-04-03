package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
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

	key, _, err := client.SSHKey.Create(ctx, sshReq)
	if err != nil {
		return diag.Errorf("error creating SSH key: %v", err)
	}

	d.SetId(key.ID)
	log.Printf("[INFO] SSH Key ID: %s", d.Id())

	return resourceVultrSSHKeyRead(ctx, d, meta)
}

func resourceVultrSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	key, _, err := client.SSHKey.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid ssh key") {
			tflog.Warn(ctx, fmt.Sprintf("Removing ssh key (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting SSH keys: %v", err)
	}

	if err := d.Set("name", key.Name); err != nil {
		return diag.Errorf("unable to set resource ssh_key `name` read value: %v", err)
	}
	if err := d.Set("ssh_key", key.SSHKey); err != nil {
		return diag.Errorf("unable to set resource ssh_key `ssh_key` read value: %v", err)
	}
	if err := d.Set("date_created", key.DateCreated); err != nil {
		return diag.Errorf("unable to set resource ssh_key `date_created` read value: %v", err)
	}

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
