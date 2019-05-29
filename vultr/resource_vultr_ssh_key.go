package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
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

	name := d.Get("name").(string)
	sshKey := d.Get("ssh_key").(string)

	key, err := client.SSHKey.Create(context.Background(), name, sshKey)
	if err != nil {
		return fmt.Errorf("Error creating SSH key: %v", err)
	}

	d.SetId(key.SSHKeyID)
	log.Printf("[INFO] SSH Key ID: %s", d.Id())

	return resourceVultrSSHKeyRead(d, meta)
}

func resourceVultrSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	keys, err := client.SSHKey.List(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting SSH keys: %v", err)
	}

	var key *govultr.SSHKey
	for i := range keys {
		if keys[i].SSHKeyID == d.Id() {
			key = &keys[i]
			break
		}
	}

	if key == nil {
		log.Printf("[WARN] Vultr SSH key (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", key.Name)
	d.Set("ssh_key", key.Key)
	d.Set("date_created", key.DateCreated)

	return nil
}

func resourceVultrSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	key := &govultr.SSHKey{
		SSHKeyID: d.Id(),
		Key:      d.Get("ssh_key").(string),
		Name:     d.Get("name").(string),
	}

	log.Printf("[INFO] Updating SSH Key: %s", d.Id())
	if err := client.SSHKey.Update(context.Background(), key); err != nil {
		return fmt.Errorf("Error updating SSH key (%s): %v", d.Id(), err)
	}

	return resourceVultrSSHKeyRead(d, meta)
}

func resourceVultrSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting SSH Key: %s", d.Id())
	if err := client.SSHKey.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("Error destroying SSH key (%s): %v", d.Id(), err)
	}

	return nil
}
