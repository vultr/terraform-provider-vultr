package vultr

import (
	"context"
	"fmt"
	"regexp"
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
				Config: testAccVultrFirewallGroup_base(rString),
			},
			{
				Config: testAccVultrFirewallGroup_base(rString) + testAccVultrFirewallRule_base(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "from_port", "3048"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "network", "10.0.0.0/32")),
			},
			{
				Config:      testAccVultrFirewallGroup_noresult(rString),
				ExpectError: regexp.MustCompile(`.* data.vultr_firewall_group.fwg: data.vultr_firewall_group.fwg: no results were found`),
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
				Config: testAccVultrFirewallGroup_base(rString) + testAccVultrFirewallRule_base(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "from_port", "3048"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "network", "10.0.0.0/32"),
				),
			},
			{
				Config: testAccVultrFirewallGroup_base(rString) + testAccVultrFirewallRule_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "udp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "from_port", "3046"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "network", "10.0.0.0/32"),
				),
			},
			{
				Config:      testAccVultrFirewallGroup_noresult(rString),
				ExpectError: regexp.MustCompile(`.* data.vultr_firewall_group.fwg: data.vultr_firewall_group.fwg: no results were found`),
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

func testAccVultrFirewallRule_base() string {
	return fmt.Sprint(`resource "vultr_firewall_rule" "tcp" {
		firewall_group_id = "${vultr_firewall_group.fwg.id}"
		protocol = "tcp"
		network = "10.0.0.0/32"
		from_port = "3048"
		}`)
}

func testAccVultrFirewallRule_update() string {
	return fmt.Sprint(`resource "vultr_firewall_rule" "tcp" {
		firewall_group_id = "${vultr_firewall_group.fwg.id}"
		protocol = "udp"
		network = "10.0.0.0/32"
		from_port = "3046"
		}`)
}
