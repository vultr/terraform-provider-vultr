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

func TestAccVultrNATGatewayPortForwardingRuleBasic(t *testing.T) {
	t.Parallel()
	pDesc := acctest.RandomWithPrefix("tf-vpc-rs")
	rLabel := acctest.RandomWithPrefix("tf-nat-gateway")
	rName := acctest.RandomWithPrefix("tf-nat-gateway-pfw-rule")

	name := "vultr_nat_gateway_port_forwarding_rule.test_nat_gateway_port_forwarding_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrNATGatewayPortForwardingRuleDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel) + testAccVultrNATGatewayPortForwardingRuleBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "protocol", "tcp"),
					resource.TestCheckResourceAttr(name, "internal_port", "555"),
					resource.TestCheckResourceAttr(name, "internal_ip", "10.1.2.3"),
					resource.TestCheckResourceAttr(name, "external_port", "123"),
					resource.TestCheckResourceAttr(name, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccVultrNATGatewayPortForwardingRuleUpdate(t *testing.T) {
	t.Parallel()
	pDesc := acctest.RandomWithPrefix("tf-vpc-rs")
	rLabel := acctest.RandomWithPrefix("tf-nat-gateway-up")
	rName := acctest.RandomWithPrefix("tf-nat-gateway-pfw-rule")

	name := "vultr_nat_gateway_port_forwarding_rule.test_nat_gateway_port_forwarding_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrNATGatewayPortForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel) + testAccVultrNATGatewayPortForwardingRuleBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "protocol", "tcp"),
					resource.TestCheckResourceAttr(name, "internal_port", "555"),
					resource.TestCheckResourceAttr(name, "internal_ip", "10.1.2.3"),
					resource.TestCheckResourceAttr(name, "external_port", "123"),
					resource.TestCheckResourceAttr(name, "enabled", "true"),
				),
			},
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBaseUpdated(rLabel) + testAccVultrNATGatewayPortForwardingRuleBaseUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "protocol", "tcp"),
					resource.TestCheckResourceAttr(name, "internal_port", "555"),
					resource.TestCheckResourceAttr(name, "internal_ip", "10.1.2.4"),
					resource.TestCheckResourceAttr(name, "external_port", "123"),
					resource.TestCheckResourceAttr(name, "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckVultrNATGatewayPortForwardingRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_nat_gateway_port_forwarding_rule" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.VPC.GetNATGatewayPortForwardingRule(context.Background(), rs.Primary.Attributes["vpc_id"], rs.Primary.Attributes["nat_gateway_id"], rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Invalid Port Forwarding Rule ID.") || strings.Contains(err.Error(), "Invalid NAT Gateway ID.") || strings.Contains(err.Error(), "Invalid VPC ID.") {
				return nil
			}
			return fmt.Errorf("error getting nat gateway port forwarding rule: %s", err)
		}

		return fmt.Errorf("nat gateway port forwarding rule %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrNATGatewayPortForwardingRuleBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_nat_gateway_port_forwarding_rule" "test_nat_gateway_port_forwarding_rule" {
			vpc_id = vultr_vpc.foo.id
			nat_gateway_id = vultr_nat_gateway.test_nat_gateway.id
			name = "%s"
			protocol = "tcp"
			internal_port = "555"
			internal_ip = "10.1.2.3"
			external_port = "123"
			enabled = true
		} `, name)
}

func testAccVultrNATGatewayPortForwardingRuleBaseUpdated(name string) string {
	return fmt.Sprintf(`
		resource "vultr_nat_gateway_port_forwarding_rule" "test_nat_gateway_port_forwarding_rule" {
			vpc_id = vultr_vpc.foo.id
			nat_gateway_id = vultr_nat_gateway.test_nat_gateway.id
			name = "%s"
			protocol = "tcp"
			internal_port = "555"
			internal_ip = "10.1.2.4"
			external_port = "123"
			enabled = true
		} `, name)
}
