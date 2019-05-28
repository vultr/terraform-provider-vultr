package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrServer(t *testing.T) {
	t.Parallel()
	rLabel := acctest.RandomWithPrefix("tf-test-ds")
	name := "data.vultr_server.server"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrServer(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "os"),
					resource.TestCheckResourceAttrSet(name, "ram"),
					resource.TestCheckResourceAttrSet(name, "disk"),
					resource.TestCheckResourceAttrSet(name, "main_ip"),
					resource.TestCheckResourceAttrSet(name, "vps_cpu_count"),
					resource.TestCheckResourceAttrSet(name, "location"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttrSet(name, "default_password"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttrSet(name, "pending_charges"),
					resource.TestCheckResourceAttrSet(name, "cost_per_month"),
					resource.TestCheckResourceAttrSet(name, "current_bandwidth"),
					resource.TestCheckResourceAttrSet(name, "allowed_bandwidth"),
					resource.TestCheckResourceAttrSet(name, "netmask_v4"),
					resource.TestCheckResourceAttrSet(name, "gateway_v4"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "power_status"),
					resource.TestCheckResourceAttrSet(name, "server_state"),
					resource.TestCheckResourceAttrSet(name, "plan_id"),
					resource.TestCheckResourceAttrSet(name, "label"),
					resource.TestCheckResourceAttrSet(name, "kvm_url"),
					resource.TestCheckResourceAttrSet(name, "auto_backups"),
					resource.TestCheckResourceAttrSet(name, "tag"),
					resource.TestCheckResourceAttrSet(name, "os_id"),
					resource.TestCheckResourceAttrSet(name, "app_id"),
					resource.TestCheckResourceAttrSet(name, "firewall_group_id"),
					resource.TestCheckResourceAttrSet(name, "v6_networks.#"),
				),
			},
		},
	})
}

func testAccCheckVultrServer(label string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "test" {
  			plan_id = "201"
  			region_id = "4"
  			os_id = "147"
  			label = "%s"
  			hostname = "testing-the-hostname"
  			enable_ipv6 = true
  			auto_backup = true
  			user_data = "unodostres!"
  			notify_activate = false
  			ddos_protection = true
  			tag = "even better tag"
		}

		data "vultr_server" "server" {
			filter {
				name = "label"
				values = ["${vultr_server.test.label}"]
			}
		}`, label)
}
