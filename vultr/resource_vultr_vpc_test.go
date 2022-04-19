package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrVPCNoCidrBlock(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-vpc-rs-nocdir")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPCConfigNoCidrBlock(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrNetworkExists("vultr_vpc.foo"),
					resource.TestCheckResourceAttr("vultr_vpc.foo", "description", rDesc),
					resource.TestCheckResourceAttrSet("vultr_vpc.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_vpc.foo", "v4_subnet"),
				),
			},
		},
	})
}

func TestAccVultrVPCWithCidrBlock(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-vpc-rs-cidr")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPCConfigWithCidrBlock(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrNetworkExists("vultr_vpc.foo"),
					resource.TestCheckResourceAttr("vultr_vpc.foo", "description", rDesc),
					resource.TestCheckResourceAttr("vultr_vpc.foo", "v4_subnet", "10.0.0.0"),
					resource.TestCheckResourceAttrSet("vultr_vpc.foo", "date_created"),
				),
			},
		},
	})
}

func testAccCheckVultrVPCDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_vpc" {
			continue
		}

		vpcID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, err := client.VPC.Get(context.Background(), vpcID)
		if err == nil {
			return fmt.Errorf("vpc still exists: %s", vpcID)
		}
	}
	return nil
}

func testAccCheckVultrVPCExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("VPC ID is not set")
		}

		vpcID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, err := client.VPC.Get(context.Background(), vpcID)
		if err != nil {
			return fmt.Errorf("VPC does not exist: %s", vpcID)
		}

		return nil
	}
}

func testAccVultrVPCConfigNoCidrBlock(rDesc string) string {
	return fmt.Sprintf(`
		resource "vultr_vpc" "foo" {
			region   = "atl"
			description = "%s"
		}
	`, rDesc)
}

func testAccVultrVPCConfigWithCidrBlock(rDesc string) string {
	return fmt.Sprintf(`
		resource "vultr_vpc" "foo" {
			region   = "atl"
			description = "%s"
			v4_subnet  = "10.0.0.0"
			v4_subnet_mask = 24
		}
	`, rDesc)
}
