package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrFirewallRuleBasic(t *testing.T) {

	rString := acctest.RandString(13)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRuleBase(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrFirewallGroupExists("vultr_firewall_group.fwg"),
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "port", "3048"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "subnet", "10.0.0.0")),
			},
		},
	})
}

func TestAccVultrFirewallRuleIcmp(t *testing.T) {

	rString := acctest.RandString(13)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRuleIcmp(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrFirewallGroupExists("vultr_firewall_group.fwg"),
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "icmp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "subnet", "0.0.0.0")),
			},
		},
	})
}

func TestAccVultrFirewallRuleUpdate(t *testing.T) {
	rString := acctest.RandString(13)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRuleBase(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "port", "3048"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "subnet", "10.0.0.0"),
				),
			},
			{
				Config: testAccVultrFirewallRuleUpdate(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "udp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "port", "3048"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "subnet", "10.0.0.0"),
				),
			},
		},
	})
}

func TestAccVultrFirewallRuleImportBasic(t *testing.T) {

	rString := acctest.RandString(13)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRuleBase(rString),
			},
			{
				ResourceName:      "vultr_firewall_rule.tcp",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testFirewallImportID("vultr_firewall_group.fwg", "vultr_firewall_rule.tcp"),
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

		groupID := rs.Primary.Attributes["firewall_group_id"]

		// If the group exists, something went wrong, probably
		_,_, groupErr := client.FirewallGroup.Get(context.Background(), groupID)
		if groupErr == nil {

			// group and rules don't throw an error from the api so the resources still exist
			_, _,_, rulesErr := client.FirewallRule.List(context.Background(), groupID, nil)
			if rulesErr == nil {
				return fmt.Errorf("firewall rules still exist: %s", rulesErr)
			}
		}
	}

	return nil
}

func testAccVultrFirewallRuleBase(desc string) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "fwg" {
			description = "%s"
		}

		resource "vultr_firewall_rule" "tcp" {
			firewall_group_id = "${vultr_firewall_group.fwg.id}"
			ip_type = "v4"
			protocol = "tcp"
			subnet = "10.0.0.0"
			subnet_size = 32
			port = "3048"
		}`, desc)
}

func testAccVultrFirewallRuleIcmp(desc string) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "fwg" {
			description = "%s"
		}

		resource "vultr_firewall_rule" "tcp" {
			firewall_group_id = "${vultr_firewall_group.fwg.id}"
			ip_type = "v4"
			protocol = "icmp"
			subnet = "0.0.0.0"
			subnet_size = 0
		}`, desc)
}

func testAccVultrFirewallRuleUpdate(desc string) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "fwg" {
			description = "%s"
		}

		resource "vultr_firewall_rule" "tcp" {
			firewall_group_id = "${vultr_firewall_group.fwg.id}"
			ip_type = "v4"
			protocol = "udp"
			subnet = "10.0.0.0"
			subnet_size = 32
			port = "3048"
		}`, desc)
}

func testFirewallImportID(g, r string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[g]
		if !ok {
			return "", fmt.Errorf("not found: %s", g)
		}

		rs2, ok := s.RootModule().Resources[r]
		if !ok {
			return "", fmt.Errorf("not found: %s", r)
		}

		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["id"], rs2.Primary.Attributes["id"]), nil
	}
}
