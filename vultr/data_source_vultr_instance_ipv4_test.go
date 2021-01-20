package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVultrInstanceIPv4_basic(t *testing.T) {
	t.Parallel()

	name := "data.vultr_instance_ipv4.test"

	serverLabel := acctest.RandomWithPrefix("tf-ds-vps-instance-ipv4")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrInstanceIPv4(serverLabel),
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

func testAccDataSourceVultrInstanceIPv4(serverLabel string) string {
	return fmt.Sprintf(`
		resource "vultr_instance" "foo" {
			plan = "vc2-1c-1gb"
			region = "ewr"
			os_id = "167"
			label = "%s"
		}

		resource "vultr_instance_ipv4" "bar" {
			instance_id = "${vultr_instance.foo.id}"
		}

		data "vultr_instance_ipv4" "test" {
			filter {
				name = "ip"
				values = ["${vultr_instance_ipv4.bar.ip}"]
			}
		}
	`, serverLabel)
}
