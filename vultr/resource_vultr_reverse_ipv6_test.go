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

func TestAccVultrReverseIPV6_basic(t *testing.T) {
	t.Parallel()

	name := "vultr_reverse_ipv6.test"

	rServerLabel := acctest.RandomWithPrefix("tf-vps-reverse-ipv6")
	reverse := fmt.Sprintf("host-%d.example.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrReverseIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReverseIPV6(rServerLabel, reverse),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReverseIPV6Exists(name),
					resource.TestCheckResourceAttrSet(name, "instance_id"),
					resource.TestCheckResourceAttrSet(name, "ip"),
					resource.TestCheckResourceAttr(name, "reverse", reverse),
				),
			},
		},
	})
}

func testAccCheckVultrReverseIPV6Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_reverse_ipv6" {
			continue
		}

		exists, err := vultrReverseIPV6Exists(rs)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("reverse IPv6 still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVultrReverseIPV6Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("reverse IPv6 not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("reverse IPv6 ID is not set")
		}

		exists, err := vultrReverseIPV6Exists(rs)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("reverse IPv6 does not exist: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccVultrReverseIPV6(rServerLabel, reverse string) string {
	return fmt.Sprintf(`
		resource "vultr_instance" "foo" {
			plan = "vc2-1c-1gb"
			region = "sea"
			os_id = "167"
			enable_ipv6 = true
			label = "%s"
		}

		resource "vultr_reverse_ipv6" "test" {
			instance_id = "${vultr_instance.foo.id}"
			ip = "${vultr_instance.foo.v6_main_ip}"
			reverse = "%s"
		}
	`, rServerLabel, reverse)
}

func vultrReverseIPV6Exists(rs *terraform.ResourceState) (bool, error) {
	client := testAccProvider.Meta().(*Client).govultrClient()

	instanceID, ok := rs.Primary.Attributes["instance_id"]
	if !ok {
		return false, errors.New("error getting instance ID")
	}

	reverseIPV6s, err := client.Instance.ListReverseIPv6(context.Background(), instanceID)
	if err != nil {
		return false, fmt.Errorf("error getting reverse IPv6s: %v", err)
	}

	ip := rs.Primary.ID

	for i := range reverseIPV6s {
		if reverseIPV6s[i].IP == ip {
			return true, nil
		}
	}

	return false, nil
}
