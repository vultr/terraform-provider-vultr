package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrLoadBalancer(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-test-ds")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrLoadBalancer(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_load_balancer.lb", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_load_balancer.lb", "status"),
					resource.TestCheckResourceAttrSet("data.vultr_load_balancer.lb", "region"),
					resource.TestCheckResourceAttrSet("data.vultr_load_balancer.lb", "label"),
				),
			},
		},
	})
}

func testAccCheckVultrLoadBalancer(label string) string {
	return fmt.Sprintf(`
		resource "vultr_load_balancer" "test" {
			region = "ewr"
			label = "%s"

			forwarding_rules {
				frontend_protocol = "http"
				frontend_port     = 80
				backend_protocol  = "http"
				backend_port      = 80
			}

		}

		data "vultr_load_balancer" "lb" {
			filter {
				name = "label"
				values = ["${vultr_load_balancer.test.label}"]
			}
		}`, label)
}
