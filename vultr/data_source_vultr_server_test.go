package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrServer(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-test-ds")
	name := "data.vultr_server.server"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrServerBase(rLabel),
			},
			{
				Config: testAccVultrServerBase(rLabel) + testAccCheckVultrServer(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "os"),
					resource.TestCheckResourceAttrSet(name, "ram"),
					resource.TestCheckResourceAttrSet(name, "disk"),
					resource.TestCheckResourceAttrSet(name, "main_ip"),
					resource.TestCheckResourceAttrSet(name, "vps_cpu_count"),
					resource.TestCheckResourceAttrSet(name, "location"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttrSet(name, "pending_charges"),
					resource.TestCheckResourceAttrSet(name, "cost_per_month"),
					resource.TestCheckResourceAttrSet(name, "current_bandwidth"),
					resource.TestCheckResourceAttrSet(name, "allowed_bandwidth"),
					resource.TestCheckResourceAttrSet(name, "netmask_v4"),
					resource.TestCheckResourceAttrSet(name, "gateway_v4"),
					resource.TestCheckResourceAttrSet(name, "power_status"),
					resource.TestCheckResourceAttrSet(name, "server_status"),
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
			{
				Config:      testAccCheckVultrServer_noResult(rLabel),
				ExpectError: regexp.MustCompile(fmt.Sprintf(".*%s: %s: no results were found", name, name)),
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
