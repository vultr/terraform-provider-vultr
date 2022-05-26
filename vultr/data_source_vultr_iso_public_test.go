package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrIsoPublic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrIsoPublicRead("10.12 x64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_iso_public.deb", "description", "10.12 x64"),
					resource.TestCheckResourceAttr(
						"data.vultr_iso_public.deb", "name", "Debian Buster"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_public.deb", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_public.deb", "md5sum"),
				),
			},
			{
				Config:      testAccVultrIsoPublicTooMany("Debian Buster"),
				ExpectError: regexp.MustCompile(`Error: your search returned too many results. Please refine your search to be more specific`),
			},
		},
	})
}

func testAccVultrIsoPublicRead(description string) string {
	return fmt.Sprintf(`
		data "vultr_iso_public" "deb" {
			filter {
				name = "description"
				values = ["%s"]
			}
		}`, description)
}

func testAccVultrIsoPublicTooMany(name string) string {
	return fmt.Sprintf(`
		data "vultr_iso_public" "deb" {
			filter {
				name = "name"
				values = ["%s"]
			}
		}`, name)
}
