package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVultrVPC(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-vpc-ds")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrNetworkConfig(rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_vpc.my_vpc", "description", rDesc),
					resource.TestCheckResourceAttrSet("data.vultr_vpc.my_vpc", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_vpc.my_vpc", "region"),
					resource.TestCheckResourceAttrSet("data.vultr_vpc.my_vpc", "v4_subnet"),
					resource.TestCheckResourceAttrSet("data.vultr_vpc.my_vpc", "v4_subnet_mask"),
				),
			},
		},
	})
}

func testAccDataSourceVultrVPCConfig(description string) string {
	return fmt.Sprintf(`
		resource "vultr_vpc" "foo" {
			region   = "ewr"
			description = "%s"
		}

		data "vultr_vpc" "my_vpc" {
			filter {
				name = "description"
				values = ["${vultr_vpc.foo.description}"]
			}
		}`, description)
}
