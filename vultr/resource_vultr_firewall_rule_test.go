package vultr

import (
	"context"
	"fmt"
	"github.com/vultr/govultr"
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrFirewallRule_basic(t *testing.T) {

	rString := acctest.RandString(13)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRule_base(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "from_port", "3048"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "network", "10.0.0.0/32")),
			},
		},
	})
}

func TestAccVultrFirewallRule_update(t *testing.T) {
	rString := acctest.RandString(13)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRule_base(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "from_port", "3048"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "network", "10.0.0.0/32"),
				),
			},
			{
				Config: testAccVultrFirewallRule_update(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "udp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "from_port", "3046"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "network", "10.0.0.0/32"),
				),
			},
		},
	})
}

func testAccCheckVultrFirewallRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).govultrClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_firewall_rule" {
			continue
		}
		groupId := rs.Primary.Attributes["firewall_group_id"]
		ipType := rs.Primary.Attributes["ip_type"]

		group, err := client.FirewallGroup.Get(context.Background(), groupId)

		if reflect.DeepEqual(group, &govultr.FirewallGroup{}) {
			// the group & rules were deleted
			return nil
		}

		firewallRules, err := client.FirewallRule.GetList(context.Background(), groupId, ipType)
		if err != nil {
			return fmt.Errorf("Error getting list of firewall rules: %s", err)
		}

		exists := false
		for i := range firewallRules {
			if strconv.Itoa(firewallRules[i].RuleNumber) == rs.Primary.ID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("Firewall rule still exists : %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccVultrFirewallRule_base(desc string) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "fwg" {
  			description = "%s"
		}			

		resource "vultr_firewall_rule" "tcp" {
			firewall_group_id = "${vultr_firewall_group.fwg.id}"
			protocol = "tcp"
			network = "10.0.0.0/32"
			from_port = "3048"
		}`, desc)
}

func testAccVultrFirewallRule_update(desc string) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "fwg" {
  			description = "%s"
		}		

		resource "vultr_firewall_rule" "tcp" {
			firewall_group_id = "${vultr_firewall_group.fwg.id}"
			protocol = "udp"
			network = "10.0.0.0/32"
			from_port = "3046"
		}`, desc)
}
