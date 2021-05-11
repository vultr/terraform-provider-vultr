package vultr

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrInstanceBasic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-vps-rs")

	name := "vultr_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrInstanceDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrInstanceBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 7 x64"),
					resource.TestCheckResourceAttr(name, "ram", "1024"),
					resource.TestCheckResourceAttr(name, "disk", "25"),
					resource.TestCheckResourceAttr(name, "region", "sea"),
					resource.TestCheckResourceAttr(name, "os_id", "167"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region", "sea"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
					resource.TestCheckResourceAttr(name, "backups", "enabled"),
					resource.TestCheckResourceAttr(name, "backups_schedule.#", "1"),
					resource.TestCheckResourceAttr(name, "backups_schedule.0.type", "weekly"),
					resource.TestCheckResourceAttr(name, "backups_schedule.0.dow", "4"),
					resource.TestCheckResourceAttr(name, "backups_schedule.0.hour", "11"),
				),
			},
		},
	})
}
func TestAccVultrInstanceUpdate(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-vps-rs-up")

	name := "vultr_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrInstanceBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 7 x64"),
					resource.TestCheckResourceAttr(name, "ram", "1024"),
					resource.TestCheckResourceAttr(name, "disk", "25"),
					resource.TestCheckResourceAttr(name, "os_id", "167"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region", "sea"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
				),
			},
			{
				Config: testAccVultrInstanceBaseUpdatedRegion(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 7 x64"),
					resource.TestCheckResourceAttr(name, "ram", "1024"),
					resource.TestCheckResourceAttr(name, "disk", "25"),
					resource.TestCheckResourceAttr(name, "os_id", "167"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region", "ewr"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
				),
			},
		},
	})
}

func TestAccVultrInstanceUpdateFirewall(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-vps-rs-upfw")

	name := "vultr_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrInstanceBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 7 x64"),
					resource.TestCheckResourceAttr(name, "ram", "1024"),
					resource.TestCheckResourceAttr(name, "disk", "25"),
					resource.TestCheckResourceAttr(name, "os_id", "167"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region", "sea"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
				),
			},
			{
				Config: testAccVultrInstanceBaseUpdateFirewall(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 7 x64"),
					resource.TestCheckResourceAttr(name, "ram", "1024"),
					resource.TestCheckResourceAttr(name, "disk", "25"),
					resource.TestCheckResourceAttr(name, "os_id", "167"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region", "sea"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
					resource.TestCheckResourceAttrSet(name, "firewall_group_id"),
				),
			},
		},
	})
}

//func TestAccVultrInstance_UpdateNetworkIDs(t *testing.T) {
//	t.Parallel()
//	rName := acctest.RandomWithPrefix("tf-vps-rs-upnid")
//
//	name := "vultr_instance.test"
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckVultrInstanceDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccVultrInstanceBase(rName),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(name, "label", rName),
//					resource.TestCheckResourceAttr(name, "os", "CentOS 7 x64"),
//					resource.TestCheckResourceAttr(name, "ram", "1024"),
//					resource.TestCheckResourceAttr(name, "disk", "25"),
//					resource.TestCheckResourceAttr(name, "os_id", "167"),
//					resource.TestCheckResourceAttr(name, "status", "active"),
//					resource.TestCheckResourceAttr(name, "power_status", "running"),
//					resource.TestCheckResourceAttr(name, "region", "sea"),
//					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
//				),
//			},
//			{
//				Config: testAccVultrInstanceBaseUpdateNetworkIDs(rName),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(name, "label", rName),
//					resource.TestCheckResourceAttr(name, "os", "CentOS 7 x64"),
//					resource.TestCheckResourceAttr(name, "ram", "1024"),
//					resource.TestCheckResourceAttr(name, "disk", "Virtual 25 GB"),
//					resource.TestCheckResourceAttr(name, "location", "Seattle"),
//					resource.TestCheckResourceAttr(name, "os_id", "167"),
//					resource.TestCheckResourceAttr(name, "status", "active"),
//					resource.TestCheckResourceAttr(name, "power_status", "running"),
//					resource.TestCheckResourceAttr(name, "region", "4"),
//					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
//					resource.TestCheckResourceAttr(name, "network_ids.#", "2"),
//				),
//			},
//		},
//	})
//}

func testAccCheckVultrInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_instance" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, err := client.Instance.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Server is pending destruction") {
				return nil
			}
			return fmt.Errorf("error getting instance: %s", err)
		}

		return fmt.Errorf("instance %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrInstanceBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_instance" "test" {
			plan = "vc2-1c-1gb"
			region = "sea"
			os_id = "167"
			label = "%s"
			hostname = "testing-the-hostname"
			enable_ipv6 = true
			activation_email = false
			ddos_protection = true
			tag = "even better tag"
			backups = "enabled"
			backups_schedule{
				type = "weekly"
				dow = 4
				hour = 11
			}
		} `, name)
}

func testAccVultrInstanceBaseUpdateFirewall(name string) string {
	return fmt.Sprintf(`
		resource "vultr_instance" "test" {
			plan = "vc2-1c-1gb"
			region = "sea"
			os_id = "167"
			label = "%s"
			hostname = "testing-the-hostname"
			enable_ipv6 = true
			activation_email = false
			ddos_protection = true
			tag = "even better tag"
			firewall_group_id = "${vultr_firewall_group.fwg.id}"
		}

		resource "vultr_firewall_group" "fwg" {
		  description = "my-cool-fw-dos"
		}
		`, name)
}

func testAccVultrInstanceBaseUpdateNetworkIDs(name string) string {
	return fmt.Sprintf(`
		resource "vultr_private_network" "foo" {
			region   = "sea"
			description = "foo"
			cidr_block  = "10.0.0.0/24"
		}

	resource "vultr_private_network" "bar" {
			region   = "sea"
			description = "bar"
			cidr_block  = "10.0.0.0/24"
		}

	resource "vultr_instance" "test" {
		plan = "vc2-1c-2gb"
		region = "sea"
		os_id = 167
		label = "%s"
		hostname = "testing-the-hostname"
		enable_ipv6 = true
		activation_email = false
		ddos_protection = true
		tag = "even better tag"
		network_ids = ["${vultr_private_network.foo.id}","${vultr_private_network.bar.id}"]
	}
	`, name)
}

func testAccVultrInstanceBaseUpdatedRegion(name string) string {
	return fmt.Sprintf(`
		resource "vultr_instance" "test" {
			plan = "vc2-1c-1gb"
			region = "ewr"
			os_id = 167
			label = "%s"
			hostname = "testing-the-hostname"
			enable_ipv6 = true
			activation_email = false
			ddos_protection = true
			tag = "even better tag"
		} `, name)
}
