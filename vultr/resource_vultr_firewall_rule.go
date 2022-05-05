package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v2"
)

func resourceVultrFirewallRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrFirewallRuleCreate,
		ReadContext:   resourceVultrFirewallRuleRead,
		DeleteContext: resourceVultrFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVultrFirewallRuleImport,
		},
		Schema: map[string]*schema.Schema{
			"firewall_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"v4", "v6"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"icmp", "tcp", "udp", "gre", "ah", "esp"}, false),
			},
			"subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"subnet_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"", "cloudflare"}, false),
				Default:      "",
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
		},
	}
}

func resourceVultrFirewallRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Creating new firewall rule")

	protocol := d.Get("protocol").(string)

	if protocol != strings.ToLower(protocol) {
		return diag.Errorf("%q is required to be all lowercase", protocol)
	}

	fwRule := &govultr.FirewallRuleReq{
		IPType:     d.Get("ip_type").(string),
		Protocol:   protocol,
		Subnet:     d.Get("subnet").(string),
		SubnetSize: d.Get("subnet_size").(int),
		Port:       d.Get("port").(string),
		Source:     d.Get("source").(string),
		Notes:      d.Get("notes").(string),
	}

	firewallGroupID := d.Get("firewall_group_id").(string)

	rule, err := client.FirewallRule.Create(ctx, firewallGroupID, fwRule)
	if err != nil {
		return diag.Errorf("error creating firewall rule : %v", err)
	}

	d.SetId(strconv.Itoa(rule.ID))

	return resourceVultrFirewallRuleRead(ctx, d, meta)
}

func resourceVultrFirewallRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	ruleID, _ := strconv.Atoi(d.Id())
	fw, err := client.FirewallRule.Get(ctx, d.Get("firewall_group_id").(string), ruleID)
	if err != nil {
		return diag.Errorf("error getting firewall rule %s: %v", d.Get("firewall_group_id").(string), err)
	}

	d.Set("ip_type", fw.IPType)
	d.Set("protocol", fw.Protocol)
	d.Set("subnet", fw.Subnet)
	d.Set("subnet_size", fw.SubnetSize)
	d.Set("notes", fw.Notes)
	d.Set("port", fw.Port)
	d.Set("source", fw.Source)

	return nil
}

func resourceVultrFirewallRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("error converting firewall rule ID")
	}

	log.Printf("[INFO] Delete firewall rule : %s", d.Id())
	if err := client.FirewallRule.Delete(ctx, d.Get("firewall_group_id").(string), id); err != nil {
		return diag.Errorf("error destroying firewall rule %s: %v", d.Id(), err)
	}
	return nil
}

func resourceVultrFirewallRuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Client).govultrClient()

	importID := d.Id()
	commaIdx := strings.IndexByte(importID, ',')

	if commaIdx == -1 {
		return nil, fmt.Errorf(`invalid import format, expected "firewallGroupID,firewallRuleID"`)
	}
	fwGroup, ruleID := importID[:commaIdx], importID[commaIdx+1:]

	rule, _ := strconv.Atoi(ruleID)
	fw, err := client.FirewallRule.Get(ctx, fwGroup, rule)
	if err != nil {
		return nil, fmt.Errorf("firewall Rule %s not found for firewall group %s", ruleID, fwGroup)
	}

	d.SetId(strconv.Itoa(fw.ID))
	d.Set("firewall_group_id", fwGroup)
	return []*schema.ResourceData{d}, nil
}
