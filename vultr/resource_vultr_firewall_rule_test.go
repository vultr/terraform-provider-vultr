package vultr

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/vultr/govultr/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVultrFirewallRule_basic(t *testing.T) {

	rString := acctest.RandString(13)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRule_base(rString),
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

func TestAccVultrFirewallRule_update(t *testing.T) {
	rString := acctest.RandString(13)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRule_base(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vultr_firewall_rule.tcp", "firewall_group_id"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "port", "3048"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.tcp", "subnet", "10.0.0.0"),
				),
			},
			{
				Config: testAccVultrFirewallRule_update(rString),
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

func TestAccVultrFirewallRule_importBasic(t *testing.T) {

	rString := acctest.RandString(13)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRule_base(rString),
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
		groupId := rs.Primary.Attributes["firewall_group_id"]

		group, err := client.FirewallGroup.Get(context.Background(), groupId)
		if err != nil {
			return fmt.Errorf("error getting firewall group %s", err)
		}
		if reflect.DeepEqual(group, &govultr.FirewallGroup{}) {
			// the group & rules were deleted
			return nil
		}

		firewallRules, _, err := client.FirewallRule.List(context.Background(), group.ID, nil)
		if err != nil {
			return fmt.Errorf("error getting list of firewall rules: %s", err)
		}

		exists := false
		for i := range firewallRules {
			if strconv.Itoa(firewallRules[i].ID) == rs.Primary.ID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("firewall rule still exists : %s", rs.Primary.ID)
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
			ip_type = "v4"
			protocol = "tcp"
			subnet = "10.0.0.0"
			subnet_size = 32
			port = "3048"
		}`, desc)
}

func testAccVultrFirewallRule_update(desc string) string {
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
