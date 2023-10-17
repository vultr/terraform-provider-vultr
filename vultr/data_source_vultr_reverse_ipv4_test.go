package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVultrReverseIPV4Basic(t *testing.T) {
	t.Parallel()

	name := "data.vultr_reverse_ipv4.test"
	serverLabel := acctest.RandomWithPrefix("tf-ds-vps-reverse-ipv4")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrReverseIPV4(serverLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "reverse"),
					resource.TestCheckResourceAttrSet(name, "netmask"),
					resource.TestCheckResourceAttrSet(name, "gateway"),
					resource.TestCheckResourceAttrSet(name, "instance_id"),
				),
			},
		},
	})
}

func testAccDataSourceVultrReverseIPV4(serverLabel string) string {
	return fmt.Sprintf(`
		variable "hostname" {
			description = "hostname"
			default     = "vultrusercontent.com"
		}

		resource "vultr_instance" "foo" {
			plan = "vc2-1c-2gb"
			region = "sea"
			os_id = "167"
			label = "%s"
		}

		resource "vultr_reverse_ipv4" "bar" {
			instance_id = "${vultr_instance.foo.id}"
			ip = "${vultr_instance.foo.main_ip}"
			reverse = "${vultr_instance.foo.main_ip}${var.hostname}"
		}

		data "vultr_reverse_ipv4" "test" {
			filter {
				name = "ip"
				values = ["${vultr_reverse_ipv4.bar.ip}"]
			}
		}
	`, serverLabel)
}
