package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrPlan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrPlan("16384 MB RAM,2x110 GB SSD,20.00 TB BW"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_plan.plan16gb", "name", "16384 MB RAM,2x110 GB SSD,20.00 TB BW"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan16gb", "vcpu_count", "4"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan16gb", "ram", "16384"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan16gb", "disk", "110"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan16gb", "bandwidth", "20.00"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan16gb", "bandwidth_gb", "20480"),
					resource.TestCheckResourceAttr("data.vultr_plan.plan16gb", "plan_type", "DEDICATED"),
					resource.TestCheckResourceAttrSet("data.vultr_plan.plan16gb", "price_per_month"),
					resource.TestCheckResourceAttrSet("data.vultr_plan.plan16gb", "available_locations.#"),
				),
			},
			{
				Config:      testAccCheckVultrPlan_tooManyResults("110"),
				ExpectError: regexp.MustCompile(`your search returned too many results. Please refine your search to be more specific`),
			},
		},
	})
}

func testAccCheckVultrPlan(name string) string {
	return fmt.Sprintf(`
		data "vultr_plan" "plan16gb" {
    		filter {
    			name = "name"
    			values = ["%s"]
			}
  		}`, name)
}

func testAccCheckVultrPlan_tooManyResults(disk string) string {
	return fmt.Sprintf(`
		data "vultr_plan" "plan16gb" {
			filter {
    			name = "disk"
    			values = ["%s"]
			}
		}`, disk)
}
