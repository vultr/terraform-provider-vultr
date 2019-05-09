package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrFirewallGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallGroup_read("My FireWall Group"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_firewall_group.fwg", "description", "My FireWall Group"),
					resource.TestCheckResourceAttr(
						"data.vultr_firewall_group.fwg", "date_created", "2019-05-08 19:19:31"),
					resource.TestCheckResourceAttr(
						"data.vultr_firewall_group.fwg", "date_modified", "2019-05-08 19:19:31"),
					resource.TestCheckResourceAttr(
						"data.vultr_firewall_group.fwg", "instance_count", "0"),
					resource.TestCheckResourceAttr(
						"data.vultr_firewall_group.fwg", "rule_count", "0"),
					resource.TestCheckResourceAttr(
						"data.vultr_firewall_group.fwg", "max_rule_count", "50"),
					resource.TestCheckResourceAttr(
						"data.vultr_firewall_group.fwg", "id", "2e353f07"),
				),
			},
			{
				Config:      testAccVultrFirewallGroup_noresult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_firewall_group.fwg: data.vultr_firewall_group.fwg: no results were found`),
			},
		},
	})
}

func testAccVultrFirewallGroup_read(description string) string {
	return fmt.Sprintf(`data "vultr_firewall_group" "fwg" {
  filter {
    name = "description"
    values = ["%s"]
  }
}`, description)
}

func testAccVultrFirewallGroup_noresult(description string) string {
	return fmt.Sprintf(`data "vultr_firewall_group" "fwg" {
  filter {
    name = "description"
    values = ["%s"]
  }
}`, description)
}
