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

func TestAccVultrBareMetalServerBasic(t *testing.T) {
	t.Parallel()
	rInt := acctest.RandInt()
	rName := acctest.RandomWithPrefix("tf-bms-rs")
	rSSH, _, err := acctest.RandSSHKeyPair("foobar")
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrBareMetalServerDestroy,
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
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "status"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "netmask_v4"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "gateway_v4"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "plan"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "v6_network"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "label"),
					resource.TestCheckResourceAttr("vultr_bare_metal_server.foo", "tags.#", "1"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "os_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "app_id"),
				),
			},
			{
				// update label, and user_data
				Config: testAccVultrBareMetalServerConfigUpdate(rInt, rSSH, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBareMetalServerExists("vultr_bare_metal_server.foo"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "os"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "ram"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "disk"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "main_ip"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "cpu_count"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "status"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "netmask_v4"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "gateway_v4"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "plan"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "label"),
					resource.TestCheckResourceAttr("vultr_bare_metal_server.foo", "tags.#", "2"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "os_id"),
					resource.TestCheckResourceAttrSet("vultr_bare_metal_server.foo", "app_id"),
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
		_,_, err := client.BareMetalServer.Get(context.Background(), bmsID)
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
		_,_, err := client.BareMetalServer.Get(context.Background(), bmsID)
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
		resource "vultr_bare_metal_server" "foo" {
			region = "ams"
			os_id = 159
			plan = "vbm-4c-32gb"
			enable_ipv6 = true
			activation_email = false
			ssh_key_ids = ["${vultr_ssh_key.foo.id}"]
			script_id = "${vultr_startup_script.foo.id}"
			user_data = "my user data"
			label = "%s"
			hostname = "%s"
			tags = [ "test tag" ]
		}
	`, rName, rName)
}

func testAccVultrBareMetalServerConfigUpdate(rInt int, rSSH, rName string) string {
	return testAccVultrSSHKeyConfigBasic(rInt, rSSH) + testAccVultrStartupScriptConfigBasic(rName) + fmt.Sprintf(`
		resource "vultr_bare_metal_server" "foo" {
			region = "ams"
			os_id = 159
			plan = "vbm-4c-32gb"
			activation_email = false
			ssh_key_ids = ["${vultr_ssh_key.foo.id}"]
			script_id = "${vultr_startup_script.foo.id}"
			user_data = "my user data"
			label = "%s-update"
			hostname = "%s"
			tags = [ "test tag", "another tag" ]
		}
	`, rName, rName)
}
