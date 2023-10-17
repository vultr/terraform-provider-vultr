package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrInstance(t *testing.T) {
	t.Parallel()
	rLabel := acctest.RandomWithPrefix("tf-test-ds")
	name := "data.vultr_instance.instance"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrInstance(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "os"),
					resource.TestCheckResourceAttrSet(name, "ram"),
					resource.TestCheckResourceAttrSet(name, "disk"),
					resource.TestCheckResourceAttrSet(name, "main_ip"),
					resource.TestCheckResourceAttrSet(name, "vcpu_count"),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttrSet(name, "allowed_bandwidth"),
					resource.TestCheckResourceAttrSet(name, "netmask_v4"),
					resource.TestCheckResourceAttrSet(name, "gateway_v4"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "power_status"),
					resource.TestCheckResourceAttrSet(name, "server_status"),
					resource.TestCheckResourceAttrSet(name, "plan"),
					resource.TestCheckResourceAttrSet(name, "label"),
					resource.TestCheckResourceAttrSet(name, "kvm"),
					resource.TestCheckResourceAttrSet(name, "os_id"),
					resource.TestCheckResourceAttrSet(name, "app_id"),
					resource.TestCheckResourceAttrSet(name, "v6_main_ip"),
					resource.TestCheckResourceAttrSet(name, "v6_network"),
					resource.TestCheckResourceAttrSet(name, "v6_network_size"),
					resource.TestCheckResourceAttr(name, "backups", "enabled"),
				),
			},
		},
	})
}

func testAccCheckVultrInstance(label string) string {
	return fmt.Sprintf(`
		resource "vultr_instance" "test" {
			plan = "vc2-1c-2gb"
			region = "sea"
			os_id = "167"
			label = "%s"
			hostname = "testing-the-hostname"
			enable_ipv6 = true
			ddos_protection = true
			backups = "enabled"
			backups_schedule{
				type = "weekly"
			}
		}

		data "vultr_instance" "instance" {
			filter {
				name = "label"
				values = ["${vultr_instance.test.label}"]
			}
		}`, label)
}
