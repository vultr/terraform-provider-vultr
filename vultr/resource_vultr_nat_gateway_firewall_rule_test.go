package vultr

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrNATGatewayFirewallRuleBasic(t *testing.T) {
	t.Parallel()
	pDesc := acctest.RandomWithPrefix("tf-vpc-rs")
	rLabel := acctest.RandomWithPrefix("tf-nat-gateway")
	rName := acctest.RandomWithPrefix("tf-nat-gateway-pfw-rule")

	name := "vultr_nat_gateway_firewall_rule.test_nat_gateway_firewall_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrNATGatewayFirewallRuleDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel) + testAccVultrNATGatewayPortForwardingRuleBase(rName),
			},
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel) + testAccVultrNATGatewayPortForwardingRuleBase(rName) + testAccVultrNATGatewayFirewallRuleBase(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "protocol", "tcp"),
					resource.TestCheckResourceAttr(name, "subnet", "1.2.3.4"),
					resource.TestCheckResourceAttr(name, "subnet_size", "24"),
					resource.TestCheckResourceAttr(name, "port", "123"),
					resource.TestCheckResourceAttr(name, "notes", "test rule"),
				),
			},
		},
	})
}

func TestAccVultrNATGatewayFirewallRuleUpdate(t *testing.T) {
	t.Parallel()
	pDesc := acctest.RandomWithPrefix("tf-vpc-rs")
	rLabel := acctest.RandomWithPrefix("tf-nat-gateway-up")
	rName := acctest.RandomWithPrefix("tf-nat-gateway-pfw-rule")

	name := "vultr_nat_gateway_firewall_rule.test_nat_gateway_firewall_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrNATGatewayFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel) + testAccVultrNATGatewayPortForwardingRuleBase(rName),
			},
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel) + testAccVultrNATGatewayPortForwardingRuleBase(rName) + testAccVultrNATGatewayFirewallRuleBase(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "protocol", "tcp"),
					resource.TestCheckResourceAttr(name, "subnet", "1.2.3.4"),
					resource.TestCheckResourceAttr(name, "subnet_size", "24"),
					resource.TestCheckResourceAttr(name, "port", "123"),
					resource.TestCheckResourceAttr(name, "notes", "test rule"),
				),
			},
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBaseUpdated(rLabel) + testAccVultrNATGatewayPortForwardingRuleBaseUpdated(rName) + testAccVultrNATGatewayFirewallRuleBaseUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "protocol", "tcp"),
					resource.TestCheckResourceAttr(name, "subnet", "1.2.3.4"),
					resource.TestCheckResourceAttr(name, "subnet_size", "24"),
					resource.TestCheckResourceAttr(name, "port", "123"),
					resource.TestCheckResourceAttr(name, "notes", "test rule updated"),
				),
			},
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel) + testAccVultrNATGatewayPortForwardingRuleBase(rName),
			},
		},
	})
}

func testAccCheckVultrNATGatewayFirewallRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_nat_gateway_firewall_rule" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.VPC.GetNATGatewayFirewallRule(context.Background(), rs.Primary.Attributes["vpc_id"], rs.Primary.Attributes["nat_gateway_id"], rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Invalid Firewall Rule ID.") || strings.Contains(err.Error(), "Invalid NAT Gateway ID.") || strings.Contains(err.Error(), "Invalid VPC ID.") {
				return nil
			}
			return fmt.Errorf("error getting nat gateway firewall rule: %s", err)
		}

		return fmt.Errorf("nat gateway firewall rule %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrNATGatewayFirewallRuleBase() string {
	return `
		resource "vultr_nat_gateway_firewall_rule" "test_nat_gateway_firewall_rule" {
			vpc_id = vultr_vpc.foo.id
			nat_gateway_id = vultr_nat_gateway.test_nat_gateway.id
			protocol = "tcp"
			subnet = "1.2.3.4"
			subnet_size = "24"
			port = "123"
			notes = "test rule"
		} `
}

func testAccVultrNATGatewayFirewallRuleBaseUpdated() string {
	return `
		resource "vultr_nat_gateway_firewall_rule" "test_nat_gateway_firewall_rule" {
			vpc_id = vultr_vpc.foo.id
			nat_gateway_id = vultr_nat_gateway.test_nat_gateway.id
			protocol = "tcp"
			subnet = "1.2.3.4"
			subnet_size = "24"
			port = "123"
			notes = "test rule updated"
		} `
}
