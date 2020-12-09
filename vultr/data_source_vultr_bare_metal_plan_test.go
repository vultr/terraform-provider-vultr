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
				Config: testAccCheckVultrBareMetalPlan("vbm-4c-32gb"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "cpu_count"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "cpu_model"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "cpu_threads"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "ram"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "disk"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "bandwidth"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "type"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "monthly_cost"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_plan.my_bm_plan", "locations.#"),
				),
			},
		},
	})
}

func testAccCheckVultrBareMetalPlan(name string) string {
	return fmt.Sprintf(`
		data "vultr_bare_metal_plan" "my_bm_plan" {
			filter {
				name = "id"
				values = ["%s"]
			}
		}`, name)
}
