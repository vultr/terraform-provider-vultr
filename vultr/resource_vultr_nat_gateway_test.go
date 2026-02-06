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

func TestAccVultrNATGatewayBasic(t *testing.T) {
	t.Parallel()
	pDesc := acctest.RandomWithPrefix("tf-vpc-rs")
	rLabel := acctest.RandomWithPrefix("tf-nat-gateway")

	name := "vultr_nat_gateway.test_nat_gateway"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrNATGatewayDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttr(name, "tag", "some random tag"),
				),
			},
		},
	})
}

func TestAccVultrNATGatewayUpdate(t *testing.T) {
	t.Parallel()
	pDesc := acctest.RandomWithPrefix("tf-vpc-rs")
	rLabel := acctest.RandomWithPrefix("tf-nat-gateway-up")

	name := "vultr_nat_gateway.test_nat_gateway"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrNATGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBase(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttr(name, "tag", "some random tag"),
				),
			},
			{
				Config: testAccVultrVPCConfigBase(pDesc) + testAccVultrNATGatewayBaseUpdated(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttr(name, "tag", "some random tag updated"),
				),
			},
		},
	})
}

func testAccCheckVultrNATGatewayDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_nat_gateway" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.VPC.GetNATGateway(context.Background(), rs.Primary.Attributes["vpc_id"], rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Invalid NAT Gateway ID.") || strings.Contains(err.Error(), "Invalid VPC ID.") {
				return nil
			}
			return fmt.Errorf("error getting nat gateway: %s", err)
		}

		return fmt.Errorf("nat gateway %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrNATGatewayBase(label string) string {
	return fmt.Sprintf(`
		resource "vultr_nat_gateway" "test_nat_gateway" {
			vpc_id = vultr_vpc.foo.id
			label = "%s"
			tag = "some random tag"
		} `, label)
}

func testAccVultrNATGatewayBaseUpdated(label string) string {
	return fmt.Sprintf(`
		resource "vultr_nat_gateway" "test_nat_gateway" {
			vpc_id = vultr_vpc.foo.id
			label = "%s"
			tag = "some random tag updated"
		} `, label)
}
