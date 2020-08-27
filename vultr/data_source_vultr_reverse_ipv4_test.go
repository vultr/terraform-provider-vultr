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

	serverLabel := acctest.RandomWithPrefix("tf-vps-reverse-ipv4")
	reverse := fmt.Sprintf("host-%d.example.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrReverseIPV4(serverLabel, reverse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "instance_id"),
					resource.TestCheckResourceAttrSet(name, "ip"),
					resource.TestCheckResourceAttr(name, "reverse", reverse),
				),
			},
		},
	})
}

func testAccDataSourceVultrReverseIPV4(serverLabel, reverse string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "foo" {
			plan_id = "201"
			region_id = "6"
			os_id = "167"
			label = "%s"
		}

		resource "vultr_reverse_ipv4" "bar" {
			instance_id = "${vultr_server.foo.id}"
			ip = "${vultr_server.foo.main_ip}"
			reverse = "%s"
		}

		data "vultr_reverse_ipv4" "test" {
			filter {
				name = "ip"
				values = ["${vultr_reverse_ipv4.bar.ip}"]
			}
		}
	`, serverLabel, reverse)
}
