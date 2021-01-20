package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrSSHKeyCreate,
		Read:   resourceVultrSSHKeyRead,
		Update: resourceVultrSSHKeyUpdate,
		Delete: resourceVultrSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceVultrSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()
	sshReq := &govultr.SSHKeyReq{
		Name:   d.Get("name").(string),
		SSHKey: d.Get("ssh_key").(string),
	}

	key, err := client.SSHKey.Create(context.Background(), sshReq)
	if err != nil {
		return fmt.Errorf("error creating SSH key: %v", err)
	}

	d.SetId(key.ID)
	log.Printf("[INFO] SSH Key ID: %s", d.Id())

	return resourceVultrSSHKeyRead(d, meta)
}

func resourceVultrSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	key, err := client.SSHKey.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting SSH keys: %v", err)
	}

	d.Set("name", key.Name)
	d.Set("ssh_key", key.SSHKey)
	d.Set("date_created", key.DateCreated)

	return nil
}

func resourceVultrSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("error updating SSH key (%s): %v", d.Id(), err)
	}

	return resourceVultrSSHKeyRead(d, meta)
}

func resourceVultrSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting SSH Key: %s", d.Id())

	if err := client.SSHKey.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error destroying SSH key (%s): %v", d.Id(), err)
	}

	return nil
}
