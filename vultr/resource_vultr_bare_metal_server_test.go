package vultr

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrBareMetalServer_basic(t *testing.T) {
	rInt := acctest.RandInt()
	rName := acctest.RandomWithPrefix("tf-test")
	rSSH, _, err := acctest.RandSSHKeyPair("foobar")
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrBareMetalServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBareMetalServerConfigBasic(rInt, rSSH, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBareMetalServerExists("vultr_bare_metal_server.foo"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "os"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "ram"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "disk"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "main_ip"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "cpu_count"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "default_password"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "status"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "netmask_v4"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "gateway_v4"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "plan_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "v6_networks.#"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "label"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "tag"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "os_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "app_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "startup_script_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "enable_ipv6"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "ssh_key_ids.#"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "userdata"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "notify_activate"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "hostname"),
				),
			},
			{
				// change the server from OS to application, update tag, label, and userdata
				Config: testAccVultrBareMetalServerConfigUpdate(rInt, rSSH, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBareMetalServerExists("vultr_bare_metal_server.foo"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "os"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "ram"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "disk"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "main_ip"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "cpu_count"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "default_password"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "status"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "netmask_v4"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "gateway_v4"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "plan_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "v6_networks.#"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "label"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "tag"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "os_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "app_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "startup_script_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "enable_ipv6"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "ssh_key_ids.#"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "userdata"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "notify_activate"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "hostname"),
				),
			},
		},
	})
}

func testAccCheckVultrBareMetalServerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_bare_metal_server" {
			continue
		}

		bmsID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()
		_, err := client.BareMetalServer.GetServer(context.Background(), bmsID)
		if err != nil {
			if strings.Contains(err.Error(), "Invalid server") {
				return nil
			}
			return fmt.Errorf("Error getting bare metal server: %s", err)
		}

		return fmt.Errorf("Bare metal server (%s) still exists", bmsID)
	}
	return nil
}

func testAccCheckVultrBareMetalServerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("bare metal server ID is not set")
		}

		bmsID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()
		_, err := client.BareMetalServer.GetServer(context.Background(), bmsID)
		if err != nil {
			if strings.Contains(err.Error(), "Invalid server") {
				return fmt.Errorf("Bare metal server (%s) does not exist", bmsID)
			}
			return fmt.Errorf("Error getting bare metal server: %s", err)
		}

		return nil
	}
}

func testAccVultrBareMetalServerConfigBasic(rInt int, rSSH, rName string) string {
	return testAccVultrSSHKeyConfigBasic(rInt, rSSH) + testAccVultrStartupScriptConfigBasic(rName) + fmt.Sprintf(`
		data "vultr_region" "singapore" {
			filter {
				name   = "name"
				values = ["Singapore"]
			}
		}

		resource "vultr_bare_metal_server" "foo" {
			region_id 		  = "${data.vultr_region.singapore.id}"
			os_id 			  = "270"
			plan_id           = "100"
			enable_ipv6       = true
			notify_activate   = false
			ssh_key_ids       = ["${vultr_ssh_key.foo.id}"]
			startup_script_id = "${vultr_startup_script.foo.id}"
			userdata          = "V2h5IHdvdWxkIHlvdSBkZWNvZGUgdGhpcz8gR0VUIEJBQ0sgVE8gV09SSyE="
			tag               = "%s"
			label             = "%s"
			hostname 		  = "%s"
		}
	`, rName, rName, rName)
}

func testAccVultrBareMetalServerConfigUpdate(rInt int, rSSH, rName string) string {
	return testAccVultrSSHKeyConfigBasic(rInt, rSSH) + testAccVultrStartupScriptConfigBasic(rName) + fmt.Sprintf(`
		data "vultr_region" "singapore" {
			filter {
				name   = "name"
				values = ["Singapore"]
			}
		}

		resource "vultr_bare_metal_server" "foo" {
			region_id 		  = "${data.vultr_region.singapore.id}"
			app_id 			  = "2"
			plan_id           = "100"
			enable_ipv6       = true
			notify_activate   = false
			ssh_key_ids       = ["${vultr_ssh_key.foo.id}"]
			startup_script_id = "${vultr_startup_script.foo.id}"
			userdata          = "V2h5Li4uPw=="
			tag               = "%s-update"
			label             = "%s-update"
			hostname 		  = "%s"
		}
	`, rName, rName, rName)
}
