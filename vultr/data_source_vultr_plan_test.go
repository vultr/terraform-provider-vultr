package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrPlan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrPlan("vc2-1c-1gb"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_plan.plan1gb", "vcpu_count", "1"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan1gb", "ram", "1024"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan1gb", "disk", "25"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan1gb", "bandwidth", "1024"),
					resource.TestCheckResourceAttrSet("data.vultr_plan.plan1gb", "monthly_cost"),
					resource.TestCheckResourceAttrSet("data.vultr_plan.plan1gb", "locations.#"),
				),
			},
		},
	})
}

func testAccCheckVultrPlan(name string) string {
	return fmt.Sprintf(`
		data "vultr_plan" "plan1gb" {
			filter {
				name = "id"
				values = ["%s"]
			}
		}`, name)
}
