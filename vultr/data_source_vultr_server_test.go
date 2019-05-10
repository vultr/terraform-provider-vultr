package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrServer("server-label"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "os"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "ram"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "disk"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "main_ip"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "vps_cpu_count"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "location"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "region_id"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "pending_charges"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "cost_per_month"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "current_bandwidth"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "allowed_bandwidth"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "netmask_v4"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "gateway_v4"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "power_status"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "server_status"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "plan_id"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "label"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "kvm_url"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "auto_backups"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "tag"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "os_id"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "app_id"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "firewall_group_id"),
					resource.TestCheckResourceAttrSet("data.vultr_server.server", "v6_networks.#"),
				),
			},
			{
				Config:      testAccCheckVultrServer_noResult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_server.server: data.vultr_server.server: no results were found`),
			},
		},
	})
}

func testAccCheckVultrServer(label string) string {
	return fmt.Sprintf(`data "vultr_server" "server" {
		filter {
		name = "label"
		values = ["%s"]
		}
		}`, label)
}

func testAccCheckVultrServer_noResult(label string) string {
	return fmt.Sprintf(`data "vultr_server" "server" {
		filter {
		name = "label"
		values = ["%s"]
		}
		}`, label)
}
