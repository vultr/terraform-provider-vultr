package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrFirewallGroup(t *testing.T) {

	rDesc := acctest.RandomWithPrefix("tf-fwg-ds")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrFirewallGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallGroup_read(rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_firewall_group.fwg", "description", rDesc),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "date_modified"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "instance_count"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "rule_count"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "max_rule_count"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "id"),
				),
			},
		},
	})
}

func testAccVultrFirewallGroup_read(description string) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "fwg" {
			description = "%s"
		}

		data "vultr_firewall_group" "fwg" {
			filter {
				name = "description"
				values = ["${vultr_firewall_group.fwg.description}"]
			}
		}`, description)
}
