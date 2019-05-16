package vultr

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrFirewallGroup(t *testing.T) {

	rString := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrFirewallGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallGroup_base(rString),
			},
			{
				Config: testAccVultrFirewallGroup_base(rString) + testAccVultrFirewallGroup_read(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_firewall_group.fwg", "description", rString),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "date_modified"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "instance_count"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "rule_count"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "max_rule_count"),
					resource.TestCheckResourceAttrSet("data.vultr_firewall_group.fwg", "id"),
				),
			},
			{
				Config:      testAccVultrFirewallGroup_noresult(rString),
				ExpectError: regexp.MustCompile(`.* data.vultr_firewall_group.fwg: data.vultr_firewall_group.fwg: no results were found`),
			},
		},
	})
}

func testAccVultrFirewallGroup_read(description string) string {
	return fmt.Sprintf(`
		data "vultr_firewall_group" "fwg" {
  			filter {
    			name = "description"
    			values = ["%s"]
  			}
		}`, description)
}

func testAccVultrFirewallGroup_noresult(description string) string {
	return fmt.Sprintf(`
		data "vultr_firewall_group" "fwg" {
  			filter {
    			name = "description"
    			values = ["%s"]
  			}
		}`, description)
}
