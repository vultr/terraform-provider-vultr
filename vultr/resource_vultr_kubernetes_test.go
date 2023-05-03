package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceVultrKubernetes(t *testing.T) {
	skipCI(t)
	rLabel := acctest.RandomWithPrefix("tf-vke-rs-")

	name := "vultr_kubernetes.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttr(name, "node_pools.#", "1"),
					resource.TestCheckResourceAttr(name, "node_pools.0.node_quantity", "1"),
					resource.TestCheckResourceAttr(name, "node_pools.0.plan", "vc2-2c-4gb"),
					resource.TestCheckResourceAttr(name, "node_pools.0.label", "tf-test-label"),
				),
			},
		},
	})
}

func TestAccResourceVultrKubernetesUpdate(t *testing.T) {
	skipCI(t)
	rLabel := acctest.RandomWithPrefix("tf-vke-rs-")

	name := "vultr_kubernetes.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttr(name, "node_pools.#", "1"),
					resource.TestCheckResourceAttr(name, "node_pools.0.node_quantity", "1"),
					resource.TestCheckResourceAttr(name, "node_pools.0.plan", "vc2-2c-4gb"),
					resource.TestCheckResourceAttr(name, "node_pools.0.label", "tf-test-label"),
				),
			},
			{
				Config: testAccVultrKubernetesUpdate(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttr(name, "node_pools.#", "1"),
					resource.TestCheckResourceAttr(name, "node_pools.0.node_quantity", "2"),
					resource.TestCheckResourceAttr(name, "node_pools.0.plan", "vc2-2c-4gb"),
					resource.TestCheckResourceAttr(name, "node_pools.0.label", "tf-test-label"),
					resource.TestCheckResourceAttr(name, "node_pools.0.auto_scaler", "true"),
					resource.TestCheckResourceAttr(name, "node_pools.0.min_nodes", "2"),
					resource.TestCheckResourceAttr(name, "node_pools.0.max_nodes", "3"),
				),
			},
		},
	})
}

func testAccVultrKubernetesBase(label string) string {
	return fmt.Sprintf(`
		resource "vultr_kubernetes" "foo" {
			region   = "ewr"
			label       = "%s"
			version = "v1.26.2+2"

			node_pools {
				node_quantity = 1
				plan = "vc2-2c-4gb"
    			label = "tf-test-label"
			}
		}`, label)
}

func testAccVultrKubernetesUpdate(label string) string {
	return fmt.Sprintf(`
		resource "vultr_kubernetes" "foo" {
			region   = "ewr"
			label       = "%s"
			version = "v1.26.2+2"

			node_pools {
				node_quantity = 2
				plan = "vc2-2c-4gb"
    			label = "tf-test-label"
				auto_scaler = true
				min_nodes = 2
				max_nodes = 3
			}
		}`, label)
}
