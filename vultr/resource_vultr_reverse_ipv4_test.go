package vultr

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrReverseIPV4Basic(t *testing.T) {
	t.Parallel()

	name := "vultr_reverse_ipv4.test"

	rServerLabel := acctest.RandomWithPrefix("tf-rs-vps-reverse-ipv4")
	reverse := fmt.Sprintf("host-%d.example.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrReverseIPV4Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReverseIPV4(rServerLabel, reverse),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReverseIPV4Exists(name),
					resource.TestCheckResourceAttrSet(name, "instance_id"),
					resource.TestCheckResourceAttrSet(name, "ip"),
					resource.TestCheckResourceAttrSet(name, "netmask"),
					resource.TestCheckResourceAttrSet(name, "gateway"),
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
			return fmt.Errorf("reverse IPv4 still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckVultrReverseIPV4Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("reverse IPv4 not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("reverse IPv4 ID is not set")
		}

		exists, err := vultrReverseIPV4Exists(rs)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("reverse IPv4 does not exist: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccVultrReverseIPV4(rServerLabel, reverse string) string {
	return fmt.Sprintf(`
		resource "vultr_instance" "foo" {
			plan = "vc2-1c-2gb"
			region = "sea"
			os_id = "167"
			label = "%s"
		}

		resource "vultr_reverse_ipv4" "test" {
			instance_id = "${vultr_instance.foo.id}"
			ip = "${vultr_instance.foo.main_ip}"
			reverse = "%s"
		}
	`, rServerLabel, reverse)
}

func vultrReverseIPV4Exists(rs *terraform.ResourceState) (bool, error) {
	client := testAccProvider.Meta().(*Client).govultrClient()

	instanceID, ok := rs.Primary.Attributes["instance_id"]
	if !ok {
		return false, errors.New("error getting instance ID")
	}

	reverseIPV4s, _, _, err := client.Instance.ListIPv4(context.Background(), instanceID, nil)
	if err != nil {
		return false, fmt.Errorf("error getting reverse IPv4s: %v", err)
	}

	ip := rs.Primary.ID

	for i := range reverseIPV4s {
		if reverseIPV4s[i].IP == ip {
			return true, nil
		}
	}

	return false, nil
}
