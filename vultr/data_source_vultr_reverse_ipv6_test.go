package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVultrReverseIPV6Basic(t *testing.T) {
	t.Parallel()

	name := "data.vultr_reverse_ipv6.test"

	rServerLabel := acctest.RandomWithPrefix("tf-vps-reverse-ipv6")
	reverse := fmt.Sprintf("host-%d.example.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
		resource "vultr_instance" "foo" {
			plan = "vc2-1c-2gb"
			region = "sea"
			os_id = "167"
			enable_ipv6 = true
			label = "%s"
		}

		resource "vultr_reverse_ipv6" "bar" {
			instance_id = "${vultr_instance.foo.id}"
			ip = "${vultr_instance.foo.v6_main_ip}"
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
