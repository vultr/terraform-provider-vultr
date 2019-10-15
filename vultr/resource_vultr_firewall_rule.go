package vultr

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVultrFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrFirewallRuleCreate,
		Read:   resourceVultrFirewallRuleRead,
		Delete: resourceVultrFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"firewall_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": { // type
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"to_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"from_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Creating new firewall rule")

	_, ipNet, err := net.ParseCIDR(d.Get("network").(string))

	if err != nil {
		return fmt.Errorf("Error parsing %q for firewall rule: %v", "cidr_block", err)
	}

	firewallGroupID := d.Get("firewall_group_id").(string)
	protocol := d.Get("protocol").(string)
	notes := d.Get("notes").(string)

	from, fromOk := d.GetOk("from_port")
	to, toOk := d.GetOk("to_port")

	port := ""

	if protocol != strings.ToLower(protocol) {
		return fmt.Errorf("%q is required to be all lowercase", protocol)
	}

	if protocol == "tcp" || protocol == "udp" {
		if fromOk {
			if fromOk && toOk {
				port = fmt.Sprintf("%d:%d", from, to)
			} else {
				port = strconv.Itoa(from.(int))
			}
		} else {
			return fmt.Errorf("%q requires at requires at least from_port or from_port and to_port", protocol)
		}
	}

	rule, err := client.FirewallRule.Create(context.Background(), firewallGroupID, protocol, port, ipNet.String(), notes)

	if err != nil {
		return fmt.Errorf("Err creating firewall group : %v", err)
	}

	d.SetId(strconv.Itoa(rule.RuleNumber))

	if ipNet.IP.To4() != nil {
		d.Set("ip_type", "v4")
	} else {
		d.Set("ip_type", "v6")
	}

	return resourceVultrFirewallRuleRead(d, meta)
}

func resourceVultrFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	firewallRuleList, err := client.FirewallRule.ListByIPType(context.Background(), d.Get("firewall_group_id").(string), d.Get("ip_type").(string))

	if err != nil {
		return fmt.Errorf("Error getting firewall rule %s: %v", d.Get("firewall_group_id").(string), err)
	}

	counter := 0
	flag := false
	for _, v := range firewallRuleList {
		if d.Id() == strconv.Itoa(v.RuleNumber) {
			flag = true
			break
		}
		counter++
	}

	if !flag {
		log.Printf("[WARN] Removing firewall rule (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("protocol", firewallRuleList[counter].Protocol)
	d.Set("network", firewallRuleList[counter].Network.String())
	d.Set("notes", firewallRuleList[counter].Notes)

	if firewallRuleList[counter].Network.IP.To4() != nil {
		d.Set("ip_type", "v4")
	} else {
		d.Set("ip_type", "v6")
	}

	// break up the ports
	from, to, err := splitFirewallRule(firewallRuleList[counter].Port)
	if err != nil {
		return fmt.Errorf("Error parsing port range for firewall rule (%s): %v", d.Id(), err)
	}

	d.Set("from_port", from)
	d.Set("to_port", to)

	return nil
}

func resourceVultrFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Delete firewall rule : %s", d.Id())
	err := client.FirewallRule.Delete(context.Background(), d.Get("firewall_group_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("Error destroying firewall rule %s: %v", d.Id(), err)
	}
	return nil
}

func splitFirewallRule(portRange string) (int, int, error) {
	if len(portRange) == 0 || strings.TrimSpace(portRange) == "-" {
		return 0, 0, nil
	}
	ports := strings.Split(portRange, "-")

	switch len(ports) {
	case 1:
		from, err := strconv.Atoi(strings.TrimSpace(ports[0]))
		if err != nil {
			return 0, 0, err
		}
		return from, 0, nil

	case 2:
		from, err := strconv.Atoi(strings.TrimSpace(ports[0]))
		if err != nil {
			return 0, 0, err
		}

		to, err := strconv.Atoi(strings.TrimSpace(ports[1]))
		if err != nil {
			return 0, 0, err
		}
		return from, to, nil

	default:
		return 0, 0, nil
	}

}
