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
					resource.TestCheckResourceAttrSet(name, "reverse"),
					resource.TestCheckResourceAttrSet(name, "gateway"),
					resource.TestCheckResourceAttrSet(name, "netmask"),
				),
			},
		},
	})
}

func testAccDataSourceVultrServerIPv4(serverLabel string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "foo" {
			plan = "vc2-1c-1gb"
			region = "ewr"
			os_id = "167"
			label = "%s"
		}

		resource "vultr_server_ipv4" "bar" {
			instance_id = "${vultr_server.foo.id}"
		}

		data "vultr_server_ipv4" "test" {
			filter {
				name = "ip"
				values = ["${vultr_server_ipv4.bar.ip}"]
			}
		}
	`, serverLabel)
}
