package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrNetwork_noCidrBlock(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-net-rs-nocdir")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrNetworkConfig_noCidrBlock(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrNetworkExists("vultr_network.foo"),
					resource.TestCheckResourceAttr("vultr_network.foo", "description", rDesc),
					resource.TestCheckResourceAttrSet("vultr_network.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_network.foo", "cidr_block"),
				),
			},
		},
	})
}

func TestAccVultrNetwork_withCidrBlock(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-net-rs-cidr")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrNetworkConfig_withCidrBlock(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrNetworkExists("vultr_network.foo"),
					resource.TestCheckResourceAttr("vultr_network.foo", "description", rDesc),
					resource.TestCheckResourceAttr("vultr_network.foo", "cidr_block", "10.0.0.0/24"),
					resource.TestCheckResourceAttrSet("vultr_network.foo", "date_created"),
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

		nets, err := client.Network.GetList(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting SSH nets: %s", err)
		}

		exists := false
		for i := range nets {
			if nets[i].NetworkID == networkID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("Network still exists: %s", networkID)
		}
	}
	return nil
}

func testAccCheckVultrNetworkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network ID is not set")
		}

		networkID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		nets, err := client.Network.GetList(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting Networks: %s", err)
		}

		exists := false
		for i := range nets {
			if nets[i].NetworkID == networkID {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Network does not exist: %s", networkID)
		}

		return nil
	}
}

func testAccVultrNetworkConfig_noCidrBlock(rDesc string) string {
	return fmt.Sprintf(`
		data "vultr_region" "atlanta" {
			filter {
			name   = "name"
			values = ["Atlanta"]
			}
		}
		resource "vultr_network" "foo" {
			region_id   = "${data.vultr_region.atlanta.id}"
			description = "%s"
		}
	`, rDesc)
}

func testAccVultrNetworkConfig_withCidrBlock(rDesc string) string {
	return fmt.Sprintf(`
		data "vultr_region" "atlanta" {
			filter {
			name   = "name"
			values = ["Atlanta"]
			}
		}
		resource "vultr_network" "foo" {
			region_id   = "${data.vultr_region.atlanta.id}"
			description = "%s"
			cidr_block  = "10.0.0.0/24"
		}
	`, rDesc)
}
