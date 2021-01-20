package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrNetworkNoCidrBlock(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-net-rs-nocdir")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrNetworkConfigNoCidrBlock(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrNetworkExists("vultr_private_network.foo"),
					resource.TestCheckResourceAttr("vultr_private_network.foo", "description", rDesc),
					resource.TestCheckResourceAttrSet("vultr_private_network.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_private_network.foo", "v4_subnet"),
				),
			},
		},
	})
}

func TestAccVultrNetworkWithCidrBlock(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-net-rs-cidr")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrNetworkConfigWithCidrBlock(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrNetworkExists("vultr_private_network.foo"),
					resource.TestCheckResourceAttr("vultr_private_network.foo", "description", rDesc),
					resource.TestCheckResourceAttr("vultr_private_network.foo", "v4_subnet", "10.0.0.0"),
					resource.TestCheckResourceAttrSet("vultr_private_network.foo", "date_created"),
				),
			},
		},
	})
}

func testAccCheckVultrNetworkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_network" {
			continue
		}

		networkID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, err := client.Network.Get(context.Background(), networkID)
		if err == nil {
			return fmt.Errorf("network still exists: %s", networkID)
		}
	}
	return nil
}

func testAccCheckVultrNetworkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("network ID is not set")
		}

		networkID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, err := client.Network.Get(context.Background(), networkID)
		if err != nil {
			return fmt.Errorf("network does not exist: %s", networkID)
		}

		return nil
	}
}

func testAccVultrNetworkConfigNoCidrBlock(rDesc string) string {
	return fmt.Sprintf(`
		resource "vultr_private_network" "foo" {
			region   = "atl"
			description = "%s"
		}
	`, rDesc)
}

func testAccVultrNetworkConfigWithCidrBlock(rDesc string) string {
	return fmt.Sprintf(`
		resource "vultr_private_network" "foo" {
			region   = "atl"
			description = "%s"
			v4_subnet  = "10.0.0.0"
			v4_subnet_mask = 24
		}
	`, rDesc)
}
