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
					resource.TestCheckResourceAttr("data.vultr_private_network.my_network", "description", rDesc),
					resource.TestCheckResourceAttrSet("data.vultr_private_network.my_network", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_private_network.my_network", "region"),
					resource.TestCheckResourceAttrSet("data.vultr_private_network.my_network", "v4_subnet"),
					resource.TestCheckResourceAttrSet("data.vultr_private_network.my_network", "v4_subnet_mask"),
				),
			},
		},
	})
}

func testAccDataSourceVultrNetworkConfig(description string) string {
	return fmt.Sprintf(`
		resource "vultr_private_network" "foo" {
			region   = "ewr"
			description = "%s"
		}

		data "vultr_private_network" "my_network" {
			filter {
				name = "description"
				values = ["${vultr_private_network.foo.description}"]
			}
		}`, description)
}
