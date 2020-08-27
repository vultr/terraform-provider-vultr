package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceVultrServerIPv4_basic(t *testing.T) {
	t.Parallel()

	name := "data.vultr_server_ipv4.test"

	serverLabel := acctest.RandomWithPrefix("tf-ds-vps-server-ipv4")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrServerIPv4(serverLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "instance_id"),
					resource.TestCheckResourceAttrSet(name, "ip"),
				),
			},
		},
	})
}

func testAccDataSourceVultrServerIPv4(serverLabel string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "foo" {
			plan_id = "201"
			region_id = "6"
			os_id = "167"
			label = "%s"
		}

		resource "vultr_server_ipv4" "bar" {
			instance_id = "123456"
		}

		data "vultr_server_ipv4" "test" {
			filter {
				name = "ip"
				values = ["${vultr_server_ipv4.bar.ip}"]
			}
		}
	`, serverLabel)
}
