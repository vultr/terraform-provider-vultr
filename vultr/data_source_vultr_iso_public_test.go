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
				Config: testAccVultrIsoPublicRead("7 x86_64 Minimal"),
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
				Config:      testAccVultrIsoPublicTooMany("Debian 9"),
				ExpectError: regexp.MustCompile(`Error: your search returned too many results. Please refine your search to be more specific`),
			},
		},
	})
}

func testAccVultrIsoPublicRead(description string) string {
	return fmt.Sprintf(`
		data "vultr_iso_public" "cent" {
			filter {
				name = "description"
				values = ["%s"]
			}
		}`, description)
}

func testAccVultrIsoPublicTooMany(name string) string {
	return fmt.Sprintf(`
		data "vultr_iso_public" "cent" {
			filter {
				name = "name"
				values = ["%s"]
			}
		}`, name)
}
