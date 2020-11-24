package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceVultrReverseIPV4_basic(t *testing.T) {
	t.Parallel()

	name := "data.vultr_reverse_ipv4.test"
	serverLabel := acctest.RandomWithPrefix("tf-ds-vps-reverse-ipv4")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrReverseIPV4(serverLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "subnet_size"),
					resource.TestCheckResourceAttrSet(name, "subnet"),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "label"),
					resource.TestCheckResourceAttrSet(name, "ip_type"),
				),
			},
		},
	})
}

func testAccDataSourceVultrReverseIPV4(serverLabel string) string {
	return fmt.Sprintf(`
		variable "hostname" {
			description = "hostname"
			default     = "vultr.com"
		}

		resource "vultr_instance" "foo" {
			plan = "vc2-1c-1gb"
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
