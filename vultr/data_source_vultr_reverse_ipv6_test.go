package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVultrReverseIPV6_basic(t *testing.T) {
	t.Parallel()

	name := "data.vultr_reverse_ipv6.test"

	rServerLabel := acctest.RandomWithPrefix("tf-vps-reverse-ipv6")
	reverse := fmt.Sprintf("host-%d.example.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrReverseIPV6(rServerLabel, reverse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "instance_id"),
					resource.TestCheckResourceAttrSet(name, "ip"),
					resource.TestCheckResourceAttr(name, "reverse", reverse),
				),
			},
		},
	})
}

func testAccDataSourceVultrReverseIPV6(rServerLabel, reverse string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "foo" {
			plan_id = "201"
			region_id = "6"
			os_id = "167"
			enable_ipv6 = true
			label = "%s"
		}

		resource "vultr_reverse_ipv6" "bar" {
			instance_id = "${vultr_server.foo.id}"
			ip = "${vultr_server.foo.v6_networks[0].v6_main_ip}"
			reverse = "%s"
		}

		data "vultr_reverse_ipv6" "test" {
			filter {
				name = "ip"
				values = ["${vultr_reverse_ipv6.bar.ip}"]
			}
		}
	`, rServerLabel, reverse)
}
