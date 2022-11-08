package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrInstances(t *testing.T) {
	t.Parallel()
	rLabel := acctest.RandomWithPrefix("tf-test-ds")
	name := "data.vultr_instances.instances"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrInstances(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "instances.0.os"),
					resource.TestCheckResourceAttrSet(name, "instances.0.ram"),
					resource.TestCheckResourceAttrSet(name, "instances.0.disk"),
					resource.TestCheckResourceAttrSet(name, "instances.0.main_ip"),
					resource.TestCheckResourceAttrSet(name, "instances.0.vcpu_count"),
					resource.TestCheckResourceAttrSet(name, "instances.0.region"),
					resource.TestCheckResourceAttrSet(name, "instances.0.date_created"),
					resource.TestCheckResourceAttrSet(name, "instances.0.allowed_bandwidth"),
					resource.TestCheckResourceAttrSet(name, "instances.0.netmask_v4"),
					resource.TestCheckResourceAttrSet(name, "instances.0.gateway_v4"),
					resource.TestCheckResourceAttrSet(name, "instances.0.status"),
					resource.TestCheckResourceAttrSet(name, "instances.0.power_status"),
					resource.TestCheckResourceAttrSet(name, "instances.0.server_status"),
					resource.TestCheckResourceAttrSet(name, "instances.0.plan"),
					resource.TestCheckResourceAttrSet(name, "instances.0.label"),
					resource.TestCheckResourceAttrSet(name, "instances.0.kvm"),
					resource.TestCheckResourceAttrSet(name, "instances.0.tag"),
					resource.TestCheckResourceAttrSet(name, "instances.0.os_id"),
					resource.TestCheckResourceAttrSet(name, "instances.0.app_id"),
					resource.TestCheckResourceAttrSet(name, "instances.0.v6_main_ip"),
					resource.TestCheckResourceAttrSet(name, "instances.0.v6_network"),
					resource.TestCheckResourceAttrSet(name, "instances.0.v6_network_size"),
					resource.TestCheckResourceAttr(name, "instances.0.backups", "enabled"),
				),
			},
		},
	})
}

func testAccCheckVultrInstances(label string) string {
	return fmt.Sprintf(`
		resource "vultr_instance" "test" {
			plan = "vc2-1c-1gb"
			region = "sea"
			os_id = "167"
			label = "%s"
			hostname = "testing-the-hostname"
			enable_ipv6 = true
			ddos_protection = true
			tag = "even better tag"
			backups = "enabled"
			backups_schedule{
				type = "weekly"
			}
		}

		data "vultr_instances" "instances" {
			filter {
				name = "label"
				values = ["${vultr_instance.test.label}"]
			}
		}`, label)
}
