package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVultrVPC2(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-vpc2-ds")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrVPC2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrVPC2Config(rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_vpc2.my_vpc2", "description", rDesc),
					resource.TestCheckResourceAttrSet("data.vultr_vpc2.my_vpc2", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_vpc2.my_vpc2", "region"),
					resource.TestCheckResourceAttrSet("data.vultr_vpc2.my_vpc2", "ip_block"),
					resource.TestCheckResourceAttrSet("data.vultr_vpc2.my_vpc2", "prefix_length"),
				),
			},
		},
	})
}

func testAccDataSourceVultrVPC2Config(description string) string {
	return fmt.Sprintf(`
		resource "vultr_vpc2" "foo" {
			region   = "ewr"
			description = "%s"
		}

		data "vultr_vpc2" "my_vpc2" {
			filter {
				name = "description"
				values = ["${vultr_vpc2.foo.description}"]
			}
		}`, description)
}
