package vultr

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceVultrLoadBalancer(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-lb-rs")
	rLabelUpdate := acctest.RandomWithPrefix("tf-test-update")
	protocol := forwardingRule()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrLoadBalancerConfig(rLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrLoadBalancerExists("vultr_load_balancer.foo"),
					resource.TestCheckResourceAttr("vultr_load_balancer.foo", "label", rLabel),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "status"),
				),
			},
			{
				Config: testAccVultrLoadBalancerConfig_attachNode(rLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrLoadBalancerExists("vultr_load_balancer.foo"),
					resource.TestCheckResourceAttr("vultr_load_balancer.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_load_balancer.foo", "attached_nodes", "[123456, 654321]"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "status"),
				),
			},
			{
				Config: testAccVultrLoadBalancerConfig_updateLabel(rLabelUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrLoadBalancerExists("vultr_load_balancer.foo"),
					resource.TestCheckResourceAttr("vultr_load_balancer.foo", "label", rLabelUpdate),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "status"),
				),
			},
			{
				Config: testAccVultrLoadBalancerConfig_updateFR(rLabel, protocol),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrLoadBalancerExists("vultr_load_balancer.foo"),
					// resource.TestCheckResourceAttr("vultr_load_balancer.foo", "forwarding_rule", protocol),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "status"),
				),
			},
			{
				Config: testAccVultrLoadBalancerConfig_detachInstance(rLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrLoadBalancerExists("vultr_load_balancer.foo"),
					resource.TestCheckResourceAttr("vultr_load_balancer.foo", "attached_nodes", "[654321]"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_load_balancer.foo", "status"),
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

		lbs, err := client.LoadBalancer.List(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting load balancers: %s", err)
		}

		exists := false
		for i := range lbs {
			lbID, _ := strconv.Atoi(id)
			if lbs[i].ID == lbID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("Load balancer still exists: %s", id)
		}
	}
	return nil
}

func testAccCheckVultrLoadBalancerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Load balancer ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		lbs, err := client.LoadBalancer.List(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting load balancers: %s", err)
		}

		exists := false
		for i := range lbs {
			lbID, _ := strconv.Atoi(id)
			if lbs[i].ID == lbID {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Load balancer does not exist: %s", id)
		}

		return nil
	}
}

func testAccVultrLoadBalancerConfig(label string) string {
	return fmt.Sprintf(`
	resource "vultr_load_balancer" "foo" {
		region_id   = "1"
		label       = "%s"

		forwarding_rules {
			frontend_protocol = "http"
			frontend_port     = 80
			backend_protocol  = "http"
			backend_port      = 80
		}
	  }
   `, label)
}

func testAccVultrLoadBalancerConfig_attachNode(label string) string {
	return fmt.Sprintf(`
	resource "vultr_load_balancer" "foo" {
		region_id   = "1"
		label       = "%s"

		forwarding_rules {
			frontend_protocol = "http"
			frontend_port     = 80
			backend_protocol  = "http"
			backend_port      = 80
		}

		attached_nodes = [123456, 654321]
	  }
   `, label)
}

func testAccVultrLoadBalancerConfig_updateLabel(label string) string {
	return fmt.Sprintf(`
	resource "vultr_load_balancer" "foo" {
		region_id   = "1"
		label       = "%s"

		forwarding_rules {
			frontend_protocol = "http"
			frontend_port     = 80
			backend_protocol  = "http"
			backend_port      = 80
		}
	  }
   `, label)
}

func testAccVultrLoadBalancerConfig_updateFR(label, protocol string) string {
	return fmt.Sprintf(`
	resource "vultr_load_balancer" "foo" {
		region_id   = "1"
		label       = "%s"

		forwarding_rules {
			frontend_protocol = "http"
			frontend_port     = 80
			backend_protocol  = "tcp"
			backend_port      = 80
		}
		`+protocol+`
	  }
   `, label, protocol)
}

func testAccVultrLoadBalancerConfig_detachInstance(label string) string {
	return fmt.Sprintf(`
	resource "vultr_load_balancer" "foo" {
		region_id   = "1"
		label       = "%s"

		attached_nodes = [654321]
	  }
   `, label)
}

func forwardingRule() string {
	return fmt.Sprintf(`
	forwarding_rules {
		frontend_protocol = "http"
		frontend_port     = 80
		backend_protocol  = "tcp"
		backend_port      = 80
	}`)
}
