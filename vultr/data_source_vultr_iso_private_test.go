package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrIsoPrivate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrIsoPrivate_read("neon-user-current.iso"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_iso_private.neon", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_private.neon", "size"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_private.neon", "status"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_private.neon", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_private.neon", "filename"),
				),
			},
			{
				Config:      testAccVultrIsoPrivate_noResults("Debian 9"),
				ExpectError: regexp.MustCompile(`.* data.vultr_iso_private.neon: data.vultr_iso_private.neon: no results were found`),
			},
		},
	})
}

func testAccVultrIsoPrivate_read(description string) string {
	return fmt.Sprintf(`data "vultr_iso_private" "neon" {
  filter {
    name = "filename"
    values = ["%s"]
  }
}`, description)
}

func testAccVultrIsoPrivate_noResults(name string) string {
	return fmt.Sprintf(`data "vultr_iso_private" "neon" {
  filter {
    name = "name"
    values = ["%s"]
  }
}`, name)
}
