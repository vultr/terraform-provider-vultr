package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVultrFirewallGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrFirewallGroupCreate,
		Read:   resourceVultrFirewallGroupRead,
		Update: resourceVultrFirewallGroupUpdate,
		Delete: resourceVultrFirewallGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceVultrFirewallGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	description := d.Get("description").(string)

	log.Printf("[INFO] Creating new firewall group")
	fwGroup, err := client.FirewallGroup.Create(context.Background(), description)

	if err != nil {
		return fmt.Errorf("Error creating firewall group: %v", err)
	}

	d.SetId(fwGroup.FirewallGroupID)

	return resourceVultrFirewallGroupRead(d, meta)
}

func resourceVultrFirewallGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	group, err := client.FirewallGroup.Get(context.Background(), d.Id())

	if err != nil {

		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] Removing firewall group (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting firewall group %s : %v", d.Id(), err)
	}

	d.Set("description", group.Description)
	d.Set("date_created", group.DateCreated)
	d.Set("date_modified", group.DateModified)
	d.Set("instance_count", group.InstanceCount)
	d.Set("rule_count", group.RuleCount)
	d.Set("max_rule_count", group.MaxRuleCount)
	return nil
}

func resourceVultrFirewallGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating firewall group: %s", d.Id())

	if d.HasChange("description") {
		err := client.FirewallGroup.ChangeDescription(context.Background(), d.Id(), d.Get("description").(string))

		if err != nil {
			return fmt.Errorf("Error updating firewall group %s : %v", d.Id(), err)
		}
	}

	return resourceVultrFirewallGroupRead(d, meta)
}

func resourceVultrFirewallGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Destroying firewall group: %s", d.Id())

	err := client.FirewallGroup.Delete(context.Background(), d.Id())

	if err != nil {
		return fmt.Errorf("Error destroying firewall group %s: %v", d.Id(), err)
	}
	return nil
}
