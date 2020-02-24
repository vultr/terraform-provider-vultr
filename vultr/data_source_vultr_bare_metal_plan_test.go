package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrBareMetalPlan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrBareMetalPlan("32768 MB RAM,2x 240 GB SSD,5.00 TB BW"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "name"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "cpu_count"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "cpu_model"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "ram"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "disk"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "bandwidth_tb"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "plan_type"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "deprecated"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "price_per_month"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "available_locations.#"),
				),
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
