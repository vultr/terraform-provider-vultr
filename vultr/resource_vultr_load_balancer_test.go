package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceVultrLoadBalancer(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-lb-rs")

	name := "vultr_load_balancer.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrLoadBalancerBase(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "ipv4"),
				),
			},
		},
	})
}

func TestAccResourceVultrLoadBalancerUpdateHealth(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-lb-rs")

	name := "vultr_load_balancer.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrLoadBalancerBase(rLabel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "ipv4"),
					resource.TestCheckResourceAttr(name, "health_check.#", "1"),
					resource.TestCheckResourceAttr(name, "health_check.0.check_interval", "15"),
					resource.TestCheckResourceAttr(name, "health_check.0.healthy_threshold", "5"),
					resource.TestCheckResourceAttr(name, "health_check.0.path", "/"),
					resource.TestCheckResourceAttr(name, "health_check.0.port", "80"),
					resource.TestCheckResourceAttr(name, "health_check.0.protocol", "http"),
					resource.TestCheckResourceAttr(name, "health_check.0.response_timeout", "5"),
					resource.TestCheckResourceAttr(name, "health_check.0.unhealthy_threshold", "5"),
				),
			},
			{
				Config: testAccVultrLoadBalancerConfig_updateHealth(rLabel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "ipv4"),
					resource.TestCheckResourceAttr(name, "health_check.#", "1"),
					resource.TestCheckResourceAttr(name, "health_check.0.check_interval", "3"),
					resource.TestCheckResourceAttr(name, "health_check.0.healthy_threshold", "4"),
					resource.TestCheckResourceAttr(name, "health_check.0.path", "/test"),
					resource.TestCheckResourceAttr(name, "health_check.0.port", "1212"),
					resource.TestCheckResourceAttr(name, "health_check.0.protocol", "http"),
					resource.TestCheckResourceAttr(name, "health_check.0.response_timeout", "1"),
					resource.TestCheckResourceAttr(name, "health_check.0.unhealthy_threshold", "2"),
				),
			},
		},
	})
}

func testAccCheckVultrLoadBalancerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_load_balancer" {
			continue
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, err := client.LoadBalancer.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("load balancer still exists: %s", id)
		}
	}
	return nil
}

func testAccVultrLoadBalancerBase(label string) string {
	return fmt.Sprintf(`
		resource "vultr_load_balancer" "foo" {
			region   = "ewr"
			label       = "%s"

			forwarding_rules {
				frontend_protocol = "http"
				frontend_port     = 80
				backend_protocol  = "http"
				backend_port      = 80
			}
		}`, label)
}

func testAccVultrLoadBalancerConfig_updateHealth(label string) string {
	return fmt.Sprintf(`
		resource "vultr_load_balancer" "foo" {
			region   = "ewr"
			label       = "%s"

			forwarding_rules {
				frontend_protocol = "http"
				frontend_port     = 80
				backend_protocol  = "http"
				backend_port      = 80
			}

			health_check {
				path = "/test"
				port = "1212"
				protocol = "http"
				response_timeout = 1
				unhealthy_threshold =2
				check_interval = 3
				healthy_threshold =4
			}
		}`, label)
}
