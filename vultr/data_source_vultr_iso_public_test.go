package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrIsoPublic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrIsoPublic_read("7 x86_64 Minimal"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_iso_public.cent", "description", "7 x86_64 Minimal"),
					resource.TestCheckResourceAttr(
						"data.vultr_iso_public.cent", "name", "CentOS 7"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_public.cent", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_public.cent", "md5sum"),
				),
			},
			{
				Config:      testAccVultrIsoPublic_tooMany("Debian 9"),
				ExpectError: regexp.MustCompile(`errors during refresh: your search returned too many results. Please refine your search to be more specific`),
			},
		},
	})
}

func testAccVultrIsoPublic_read(description string) string {
	return fmt.Sprintf(`
		data "vultr_iso_public" "cent" {
			filter {
				name = "description"
				values = ["%s"]
			}
		}`, description)
}

func testAccVultrIsoPublic_tooMany(name string) string {
	return fmt.Sprintf(`
		data "vultr_iso_public" "cent" {
			filter {
				name = "name"
				values = ["%s"]
			}
		}`, name)
}
