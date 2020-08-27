package vultr

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVultrServerIPv4_basic(t *testing.T) {
	t.Parallel()

	name := "vultr_server_ipv4.test"

	serverLabel := acctest.RandomWithPrefix("tf-vps-server-ipv4")
	reboot := "false"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrServerIPV4Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrServerIPV4(serverLabel, reboot),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckvultrServerIPV4Exists(name),
					resource.TestCheckResourceAttrSet(name, "instance_id"),
					resource.TestCheckResourceAttrSet(name, "ip"),
					resource.TestCheckResourceAttr(name, "reboot", reboot),
				),
			},
		},
	})
}

func testAccCheckVultrServerIPV4Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_server_ipv4" {
			continue
		}

		exists, err := vultrServerIPV4Exists(rs)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("IPv4 still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckvultrServerIPV4Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("IPv4 not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("IPv4 ID is not set")
		}

		exists, err := vultrServerIPV4Exists(rs)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("IPv4 does not exist: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccVultrServerIPV4(serverLabel, reboot string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "foo" {
			plan_id = "201"
			region_id = "6"
			os_id = "167"
			label = "%s"
		}

		resource "vultr_server_ipv4" "test" {
			instance_id = "${vultr_server.foo.id}"
			ip = "123.123.123.123"
			reboot = "%s"
		}
	`, serverLabel, reboot)
}

func vultrServerIPV4Exists(rs *terraform.ResourceState) (bool, error) {
	client := testAccProvider.Meta().(*Client).govultrClient()

	instanceID, ok := rs.Primary.Attributes["instance_id"]
	if !ok {
		return false, errors.New("Error getting instance ID")
	}

	ipv4s, err := client.Server.IPV4Info(context.Background(), instanceID, true)
	if err != nil {
		return false, fmt.Errorf("Error getting IPv4s: %v", err)
	}

	ip := rs.Primary.ID

	for i := range ipv4s {
		if ipv4s[i].IP == ip {
			return true, nil
		}
	}

	return false, nil
}
