package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrRegion("Miami"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_region.miami", "id", "39"),
					resource.TestCheckResourceAttr("data.vultr_region.miami", "name", "Miami"),
					resource.TestCheckResourceAttr("data.vultr_region.miami", "country", "US"),
					resource.TestCheckResourceAttr("data.vultr_region.miami", "continent", "North America"),
					resource.TestCheckResourceAttr("data.vultr_region.miami", "state", "FL"),
					resource.TestCheckResourceAttr("data.vultr_region.miami", "regioncode", "MIA"),
					resource.TestCheckResourceAttrSet("data.vultr_region.miami", "ddos_protection"),
					resource.TestCheckResourceAttrSet("data.vultr_region.miami", "block_storage"),
				),
			},
			{
				Config:      testAccCheckVultrRegion_noResult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_region.miami: data.vultr_region.miami: no results were found`),
			},
			{
				Config:      testAccCheckVultrRegion_tooManyResults("US"),
				ExpectError: regexp.MustCompile(`.* data.vultr_region.miami: data.vultr_region.miami: your search returned too many results : 8. Please refine your search to be more specific`),
			},
		},
	})
}

func testAccCheckVultrRegion(name string) string {
	return fmt.Sprintf(`
		data "vultr_region" "miami" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}

func testAccCheckVultrRegion_noResult(name string) string {
	return fmt.Sprintf(`
		data "vultr_region" "miami" {
    	filter {
    	name = "name"
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
