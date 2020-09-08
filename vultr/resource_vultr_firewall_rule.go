package vultr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vultr/govultr/v2"
	"log"
	"strconv"
)

func resourceVultrFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceVultrFirewallRuleCreate,
		Read:     resourceVultrFirewallRuleRead,
		Delete:   resourceVultrFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			//State: resourceVultrFirewallRuleImport,
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
				ValidateFunc: validation.StringInSlice([]string{"tcmp", "tcp", "udp", "gre"}, false),
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

func resourceVultrFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Creating new firewall rule")

	fwRule := &govultr.FirewallRuleReq{
		IPType:     d.Get("ip_type").(string),
		Protocol:   d.Get("protocol").(string),
		Subnet:     d.Get("subnet").(string),
		SubnetSize: d.Get("subnet_size").(int),
		Port:       d.Get("port").(string),
		Source:     d.Get("source").(string),
		Notes:      d.Get("notes").(string),
	}

	firewallGroupID := d.Get("firewall_group_id").(string)

	rule, err := client.FirewallRule.Create(context.Background(), firewallGroupID, fwRule)
	if err != nil {
		return fmt.Errorf("error creating firewall rule : %v", err)
	}

	d.SetId(strconv.Itoa(rule.ID))
	//
	//Port:
	//	"",
	//		Source:     "",
	//		Notes:      "",

	//protocol := d.Get("protocol").(string)
	//notes := d.Get("notes").(string)
	//
	//from, fromOk := d.GetOk("from_port")
	//to, toOk := d.GetOk("to_port")
	//
	//port := ""

	//if protocol != strings.ToLower(protocol) {
	//	return fmt.Errorf("%q is required to be all lowercase", protocol)
	//}
	//
	//if protocol == "tcp" || protocol == "udp" {
	//	if fromOk {
	//		if fromOk && toOk {
	//			port = fmt.Sprintf("%d:%d", from, to)
	//		} else {
	//			port = strconv.Itoa(from.(int))
	//		}
	//	} else {
	//		return fmt.Errorf("%q requires at requires at least from_port or from_port and to_port", protocol)
	//	}
	//}

	return resourceVultrFirewallRuleRead(d, meta)
}

func resourceVultrFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	// todo conver this to a single get call
	options := &govultr.ListOptions{
		PerPage: 1,
	}
	var fwRule *govultr.FirewallRule
	for {
		fw, meta, err := client.FirewallRule.List(context.Background(), d.Get("firewall_group_id").(string), options)
		if err != nil {
			return fmt.Errorf("1 error getting firewall rule %s: %v", d.Get("firewall_group_id").(string), err)
		}

		if strconv.Itoa(fw[0].ID) == d.Id() {
			fwRule = &fw[0]
			break
		}

		if meta.Links.Next == "" {
			log.Printf("[WARN] Removing firewall rule (%s) because it is gone", d.Id())
			d.SetId("")
			return fmt.Errorf("2 error getting firewall rule %s: %v", d.Get("firewall_group_id").(string), err)
		}

		options.Cursor = meta.Links.Next
	}

	d.Set("ip_type", fwRule.Type)
	d.Set("protocol", fwRule.Protocol)
	d.Set("subnet", fwRule.Subnet)
	d.Set("subnet_size", fwRule.SubnetSize)
	d.Set("notes", fwRule.Notes)
	d.Set("port", fwRule.Port)
	d.Set("source", fwRule.Source)

	return nil
}

func resourceVultrFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("error convering firewall rule ID")
	}

	log.Printf("[INFO] Delete firewall rule : %s", d.Id())
	if err := client.FirewallRule.Delete(context.Background(), d.Get("firewall_group_id").(string), id); err != nil {
		return fmt.Errorf("error destroying firewall rule %s: %v", d.Id(), err)
	}
	return nil
}

// todo
//func resourceVultrFirewallRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
//	client := meta.(*Client).govultrClient()
//
//	importID := d.Id()
//	commaIdx := strings.IndexByte(importID, ',')
//
//	if commaIdx == -1 {
//		return nil, fmt.Errorf(`invalid import format, expected "firewallGroupID,firewallRuleID"`)
//	}
//	fwGroup, ruleID := importID[:commaIdx], importID[commaIdx+1:]
//
//	options := &govultr.ListOptions{
//		PerPage: 25,
//	}
//
//	var rule *govultr.FirewallRule
//	for {
//		rules, meta, err := client.FirewallRule.List(context.Background(), fwGroup, options)
//		if err != nil {
//			return nil, fmt.Errorf("error getting Firewall Rules for Firewall Group %s: %v", fwGroup, err)
//		}
//		for _, v := range rules {
//			if strconv.Itoa(v.ID) == ruleID {
//				rule = &v
//				break
//			}
//		}
//
//		if rule == nil {
//			break
//		}
//
//		if meta.Links.Next == "" {
//			return nil, fmt.Errorf("firewall Rule %s not found for firewall group %s", ruleID, fwGroup)
//		}
//
//		options.Cursor = meta.Links.Next
//	}
//
//	d.SetId(strconv.Itoa(rule.ID))
//	d.Set("firewall_group_id", fwGroup)
//	return []*schema.ResourceData{d}, nil
//}
