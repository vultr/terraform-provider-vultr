package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceVultrReservedIP(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-rip-ds")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIP_read(rLabel),
			},
			{
				Config: testAccVultrReservedIP_read(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.foo", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.foo", "label"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.foo", "ip_type"),
				),
			},
		},
	})
}

func testAccVultrReservedIP_read(label string) string {
	return fmt.Sprintf(`
		resource "vultr_reserved_ip" "bar" {
		label = "%s"
		region = "sea"
		ip_type = "v4"
	}

		data "vultr_reserved_ip" "foo" {
			filter {
				name = "label"
				values = ["${vultr_reserved_ip.bar.label}"]
			}
		}
		`, label)
}
