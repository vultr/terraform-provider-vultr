package vultr

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVultrServer_Basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-vps-rs")

	name := "vultr_server.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrServerBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 6 i386"),
					resource.TestCheckResourceAttr(name, "ram", "1024 MB"),
					resource.TestCheckResourceAttr(name, "disk", "Virtual 25 GB"),
					resource.TestCheckResourceAttr(name, "location", "Seattle"),
					resource.TestCheckResourceAttr(name, "os_id", "147"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region_id", "4"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
				),
			},
		},
	})
}
func TestAccVultrServer_Update(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-vps-rs-up")

	name := "vultr_server.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrServerBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 6 i386"),
					resource.TestCheckResourceAttr(name, "ram", "1024 MB"),
					resource.TestCheckResourceAttr(name, "disk", "Virtual 25 GB"),
					resource.TestCheckResourceAttr(name, "location", "Seattle"),
					resource.TestCheckResourceAttr(name, "os_id", "147"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region_id", "4"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
				),
			},
			{
				Config: testAccVultrServerBaseUpdatedRegion(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 6 i386"),
					resource.TestCheckResourceAttr(name, "ram", "2048 MB"),
					resource.TestCheckResourceAttr(name, "disk", "Virtual 55 GB"),
					resource.TestCheckResourceAttr(name, "location", "New Jersey"),
					resource.TestCheckResourceAttr(name, "os_id", "147"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region_id", "1"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
				),
			},
		},
	})
}

func TestAccVultrServer_UpdateFirewall(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-vps-rs-upfw")

	name := "vultr_server.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrServerBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 6 i386"),
					resource.TestCheckResourceAttr(name, "ram", "1024 MB"),
					resource.TestCheckResourceAttr(name, "disk", "Virtual 25 GB"),
					resource.TestCheckResourceAttr(name, "location", "Seattle"),
					resource.TestCheckResourceAttr(name, "os_id", "147"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region_id", "4"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
				),
			},
			{
				Config: testAccVultrServerBaseUpdateFirewall(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 6 i386"),
					resource.TestCheckResourceAttr(name, "ram", "1024 MB"),
					resource.TestCheckResourceAttr(name, "disk", "Virtual 25 GB"),
					resource.TestCheckResourceAttr(name, "location", "Seattle"),
					resource.TestCheckResourceAttr(name, "os_id", "147"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region_id", "4"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
					resource.TestCheckResourceAttrSet(name, "firewall_group_id"),
				),
			},
		},
	})
}

func TestAccVultrServer_UpdateNetworkIDs(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-vps-rs-upnid")

	name := "vultr_server.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrServerBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 6 i386"),
					resource.TestCheckResourceAttr(name, "ram", "1024 MB"),
					resource.TestCheckResourceAttr(name, "disk", "Virtual 25 GB"),
					resource.TestCheckResourceAttr(name, "location", "Seattle"),
					resource.TestCheckResourceAttr(name, "os_id", "147"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region_id", "4"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
				),
			},
			{
				Config: testAccVultrServerBaseUpdateNetworkIDs(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "os", "CentOS 6 i386"),
					resource.TestCheckResourceAttr(name, "ram", "1024 MB"),
					resource.TestCheckResourceAttr(name, "disk", "Virtual 25 GB"),
					resource.TestCheckResourceAttr(name, "location", "Seattle"),
					resource.TestCheckResourceAttr(name, "os_id", "147"),
					resource.TestCheckResourceAttr(name, "status", "active"),
					resource.TestCheckResourceAttr(name, "power_status", "running"),
					resource.TestCheckResourceAttr(name, "region_id", "4"),
					resource.TestCheckResourceAttr(name, "tag", "even better tag"),
					resource.TestCheckResourceAttr(name, "network_ids.#", "2"),
				),
			},
		},
	})
}

func testAccCheckVultrServerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_server" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, err := client.Server.GetServer(context.Background(), rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Invalid server") {
				return nil
			}
			return fmt.Errorf("Error getting server: %s", err)
		}

		return fmt.Errorf("Server %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrServerBase(name string) string {
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
		} `, name)
}

func testAccVultrServerBaseUpdateFirewall(name string) string {
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
			firewall_group_id = "${vultr_firewall_group.fwg.id}"
		} 
		
		resource "vultr_firewall_group" "fwg" {
		  description = "my-cool-fw-dos"
		}
		`, name)
}

func testAccVultrServerBaseUpdateNetworkIDs(name string) string {
	return fmt.Sprintf(`
		resource "vultr_network" "foo" {
			region_id   = "4"
			description = "foo"
			cidr_block  = "10.0.0.0/24"
		}

    	resource "vultr_network" "bar" {
			region_id   = "4"
			description = "bar"
			cidr_block  = "10.0.0.0/24"
		}

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
        	network_ids = ["${vultr_network.foo.id}","${vultr_network.bar.id}"]
		}
		`, name)
}

func testAccVultrServerBaseUpdatedRegion(name string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "test" {
  			plan_id = "202"
  			region_id = "1"
  			os_id = "147"
  			label = "%s"
  			hostname = "testing-the-hostname"
  			enable_ipv6 = true
  			auto_backup = true
  			user_data = "unodostres!"
  			notify_activate = false
  			ddos_protection = true
  			tag = "even better tag"
		} `, name)
}
