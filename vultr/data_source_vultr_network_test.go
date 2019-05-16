package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVultrNetwork(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrNetworkConfig_noCidrBlock(rDesc),
			},
			{
				Config: testAccVultrNetworkConfig_noCidrBlock(rDesc) + testAccDataSourceVultrNetworkConfig(rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_network.my_network", "description", rDesc),
					resource.TestCheckResourceAttrSet("data.vultr_network.my_network", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_network.my_network", "region_id"),
					resource.TestCheckResourceAttrSet("data.vultr_network.my_network", "cidr_block"),
				),
			},
			{
				Config:      testAccDataSourceVultrNetworkConfig(rDesc),
				ExpectError: regexp.MustCompile(`.* data.vultr_network.my_network: data.vultr_network.my_network: no results were found`),
			},
		},
	})
}

func testAccDataSourceVultrNetworkConfig(description string) string {
	return fmt.Sprintf(`
	data "vultr_network" "my_network" {
		filter {
			name = "description"
			values = ["%s"]
		}
  	}`, description)
}
