package vultr

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vultr/govultr"
)

func TestAccVultrFirewallGroup_basic(t *testing.T) {

	rString := acctest.RandString(10)
	updatedString := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrFirewallGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallGroup_base(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrFirewallGroupExists("vultr_firewall_group.fwg"),
					resource.TestCheckResourceAttr(
						"vultr_firewall_group.fwg", "description", rString),
				),
			},
			{
				Config: testAccVultrFirewallGroup_update(updatedString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrFirewallGroupExists("vultr_firewall_group.fwg"),
					resource.TestCheckResourceAttr(
						"vultr_firewall_group.fwg", "description", updatedString),
				),
			},
		},
	})
}

func testAccCheckVultrFirewallGroupDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*Client).govultrClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_firewall_group" {
			continue
		}

		group, err := client.FirewallGroup.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Firewall group still exists : %s", err)
		}

		fwGroup := &govultr.FirewallGroup{}
		if !reflect.DeepEqual(group, fwGroup) {
			return fmt.Errorf("Firewall group still exists : %s", err)
		}
	}

	return nil
}

func testAccCheckVultrFirewallGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("FirewallGroupID is not set")
		}

		keyID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		key, err := client.FirewallGroup.Get(context.Background(), keyID)
		if err != nil {
			return fmt.Errorf("Error getting firewall group : %s", err)
		}

		if key.FirewallGroupID != keyID {
			return fmt.Errorf("FirewallGroup does not exist: %s", keyID)
		}

		return nil
	}
}

func testAccVultrFirewallGroup_base(description string) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "fwg" {
  			description = "%s"
		}`, description)
}

func testAccVultrFirewallGroup_update(description string) string {
	return fmt.Sprintf(`
	resource "vultr_firewall_group" "fwg" {
  		description = "%s"
	}`, description)
}
