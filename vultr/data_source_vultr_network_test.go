package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceVultrNetwork(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-network-ds")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrNetworkConfig(rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_network.my_network", "description", rDesc),
					resource.TestCheckResourceAttrSet("data.vultr_network.my_network", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_network.my_network", "region_id"),
					resource.TestCheckResourceAttrSet("data.vultr_network.my_network", "cidr_block"),
				),
			},
		},
	})
}

func testAccDataSourceVultrNetworkConfig(description string) string {
	return fmt.Sprintf(`
		resource "vultr_network" "foo" {
			region_id   = 4
			description = "%s"
		}

		data "vultr_network" "my_network" {
			filter {
				name = "description"
				values = ["${vultr_network.foo.description}"]
			}
  		}`, description)
}
