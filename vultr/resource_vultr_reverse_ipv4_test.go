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

func TestAccVultrReverseIPV4_basic(t *testing.T) {
	t.Parallel()

	name := "vultr_reverse_ipv4.test"

	rServerLabel := acctest.RandomWithPrefix("tf-vps-reverse-ipv4")
	reverse := fmt.Sprintf("host-%d.example.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrReverseIPV4Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReverseIPV4(rServerLabel, reverse),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReverseIPV4Exists(name),
					resource.TestCheckResourceAttrSet(name, "instance_id"),
					resource.TestCheckResourceAttrSet(name, "ip"),
					resource.TestCheckResourceAttr(name, "reverse", reverse),
				),
			},
		},
	})
}

func testAccCheckVultrReverseIPV4Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_reverse_ipv4" {
			continue
		}

		exists, err := vultrReverseIPV4Exists(rs)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Reverse IPv4 still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVultrReverseIPV4Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Reverse IPv4 not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Reverse IPv4 ID is not set")
		}

		exists, err := vultrReverseIPV4Exists(rs)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("Reverse IPv4 does not exist: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccVultrReverseIPV4(rServerLabel, reverse string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "foo" {
			plan_id = "201"
			region_id = "6"
			os_id = "167"
			label = "%s"
		}

		resource "vultr_reverse_ipv4" "test" {
			instance_id = "${vultr_server.foo.id}"
			ip = "${vultr_server.foo.v6_networks[0].v6_main_ip}"
			reverse = "%s"
		}
	`, rServerLabel, reverse)
}

func vultrReverseIPV4Exists(rs *terraform.ResourceState) (bool, error) {
	client := testAccProvider.Meta().(*Client).govultrClient()

	instanceID, ok := rs.Primary.Attributes["instance_id"]
	if !ok {
		return false, errors.New("Error getting instance ID")
	}

	reverseIPV4s, err := client.Server.IPV4Info(context.Background(), instanceID, true)
	if err != nil {
		return false, fmt.Errorf("Error getting reverse IPv4s: %v", err)
	}

	ip := rs.Primary.ID

	for i := range reverseIPV4s {
		if reverseIPV4s[i].IP == ip {
			return true, nil
		}
	}

	return false, nil
}
