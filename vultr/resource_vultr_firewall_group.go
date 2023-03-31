package vultr

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrFirewallGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrFirewallGroupCreate,
		ReadContext:   resourceVultrFirewallGroupRead,
		UpdateContext: resourceVultrFirewallGroupUpdate,
		DeleteContext: resourceVultrFirewallGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Default:  "",
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

func resourceVultrFirewallGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	fwReq := &govultr.FirewallGroupReq{Description: d.Get("description").(string)}
	log.Printf("[INFO] Creating new firewall group")
	fwGroup, _, err := client.FirewallGroup.Create(ctx, fwReq)
	if err != nil {
		return diag.Errorf("error creating firewall group: %v", err)
	}

	d.SetId(fwGroup.ID)

	return resourceVultrFirewallGroupRead(ctx, d, meta)
}

func resourceVultrFirewallGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	group, _, err := client.FirewallGroup.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "\"status\":404") {
			log.Printf("[WARN] Removing firewall group (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting firewall group %s : %v", d.Id(), err)
	}

	if err := d.Set("description", group.Description); err != nil {
		return diag.Errorf("unable to set resource firewall_group `description` read value: %v", err)
	}
	if err := d.Set("date_created", group.DateCreated); err != nil {
		return diag.Errorf("unable to set resource firewall_group `date_created` read value: %v", err)
	}
	if err := d.Set("date_modified", group.DateModified); err != nil {
		return diag.Errorf("unable to set resource firewall_group `date_modified` read value: %v", err)
	}
	if err := d.Set("instance_count", group.InstanceCount); err != nil {
		return diag.Errorf("unable to set resource firewall_group `instance_count` read value: %v", err)
	}
	if err := d.Set("rule_count", group.RuleCount); err != nil {
		return diag.Errorf("unable to set resource firewall_group `rule_count` read value: %v", err)
	}
	if err := d.Set("max_rule_count", group.MaxRuleCount); err != nil {
		return diag.Errorf("unable to set resource firewall_group `max_rule_count` read value: %v", err)
	}
	return nil
}

func resourceVultrFirewallGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating firewall group: %s", d.Id())

	fwReq := &govultr.FirewallGroupReq{Description: d.Get("description").(string)}
	if err := client.FirewallGroup.Update(ctx, d.Id(), fwReq); err != nil {
		return diag.Errorf("error updating firewall group %s : %v", d.Id(), err)
	}

	return resourceVultrFirewallGroupRead(ctx, d, meta)
}

func resourceVultrFirewallGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting firewall group: %s", d.Id())

	if err := client.FirewallGroup.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying firewall group %s: %v", d.Id(), err)
	}
	return nil
}
