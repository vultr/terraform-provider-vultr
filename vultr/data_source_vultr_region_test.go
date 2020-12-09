package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrRegion("mia"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_region.miami", "id", "mia"),
					resource.TestCheckResourceAttr("data.vultr_region.miami", "country", "US"),
					resource.TestCheckResourceAttr("data.vultr_region.miami", "continent", "North America"),
					resource.TestCheckResourceAttr("data.vultr_region.miami", "city", "Miami"),
				),
			},
			{
				Config:      testAccCheckVultrRegion_tooManyResults("US"),
				ExpectError: regexp.MustCompile(`your search returned too many results. Please refine your search to be more specific`),
			},
		},
	})
}

func testAccCheckVultrRegion(name string) string {
	return fmt.Sprintf(`
		data "vultr_region" "miami" {
			filter {
				name = "id"
				values = ["%s"]
			}
		}`, name)
}

func testAccCheckVultrRegion_tooManyResults(country string) string {
	return fmt.Sprintf(`
		data "vultr_region" "miami" {
			filter {
				name = "country"
				values = ["%s"]
			}
		}`, country)
}
