package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrBareMetalPlan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrBareMetalPlan("32768 MB RAM,4x 240 GB SSD,1.00 TB BW"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_bare_metal_plan.my_bm_plan", "name", "32768 MB RAM,4x 240 GB SSD,1.00 TB BW"),
					resource.TestCheckResourceAttr("data.vultr_bare_metal_plan.my_bm_plan", "cpu_count", "12"),
					resource.TestCheckResourceAttr("data.vultr_bare_metal_plan.my_bm_plan", "cpu_model", "E-2186G"),
					resource.TestCheckResourceAttr("data.vultr_bare_metal_plan.my_bm_plan", "ram", "32768"),
					resource.TestCheckResourceAttr("data.vultr_bare_metal_plan.my_bm_plan", "disk", "4x 240 GB SSD"),
					resource.TestCheckResourceAttr("data.vultr_bare_metal_plan.my_bm_plan", "bandwidth_tb", "1"),
					resource.TestCheckResourceAttr("data.vultr_bare_metal_plan.my_bm_plan", "plan_type", "SSD"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "deprecated"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "price_per_month"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "available_locations.#"),
				),
			},
			{
				Config:      testAccCheckVultrBareMetalPlan_tooManyResults(32768),
				ExpectError: regexp.MustCompile(`your search returned too many results. Please refine your search to be more specific`),
			},
		},
	})
}

func testAccCheckVultrBareMetalPlan(name string) string {
	return fmt.Sprintf(`
		data "vultr_bare_metal_plan" "my_bm_plan" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}

func testAccCheckVultrBareMetalPlan_noResult(name string) string {
	return fmt.Sprintf(`
		data "vultr_bare_metal_plan" "my_bm_plan" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}

func testAccCheckVultrBareMetalPlan_tooManyResults(ram int) string {
	return fmt.Sprintf(`
		data "vultr_bare_metal_plan" "my_bm_plan" {
    	filter {
    	name = "ram"
    	values = [%v]
	}
  	}`, ram)
}
