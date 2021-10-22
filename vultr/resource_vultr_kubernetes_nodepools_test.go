package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceVultrKubernetesNodePools(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-vke-rs")
	rNP := acctest.RandomWithPrefix("tf-vke-np")

	name := "vultr_kubernetes_node_pools.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel) + testAccVultrKubernetesNodePoolsBase(rNP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rNP),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "tag"),
					resource.TestCheckResourceAttr(name, "nodes.#", "1"),
					resource.TestCheckResourceAttr(name, "plan", "vc2-2c-4gb"),
				),
			},
		},
	})
}

func TestAccResourceVultrKubernetesNodePoolsUpdate(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-vke-rs")
	rNP := acctest.RandomWithPrefix("tf-vke-np")

	name := "vultr_kubernetes_node_pools.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel) + testAccVultrKubernetesNodePoolsBase(rNP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rNP),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "tag"),
					resource.TestCheckResourceAttr(name, "nodes.#", "1"),
					resource.TestCheckResourceAttr(name, "plan", "vc2-2c-4gb"),
				),
			},
			{
				Config: testAccVultrKubernetesBase(rLabel) + testAccVultrKubernetesNodePoolsUpdate(rNP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rNP),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "tag"),
					resource.TestCheckResourceAttr(name, "nodes.#", "2"),
					resource.TestCheckResourceAttr(name, "plan", "vc2-2c-4gb"),
				),
			},
		},
	})
}
func testAccVultrKubernetesNodePoolsBase(label string) string {
	return fmt.Sprintf(`
		resource "vultr_kubernetes_node_pools" "foo" {
    			cluster_id = vultr_kubernetes.foo.id
				node_quantity = 1
				plan = "vc2-2c-4gb"
    			label = "%s"
    			tag = "test23"
		}`, label)
}

func testAccVultrKubernetesNodePoolsUpdate(label string) string {
	return fmt.Sprintf(`
		resource "vultr_kubernetes_node_pools" "foo" {
    			cluster_id = vultr_kubernetes.foo.id
				node_quantity = 2
				plan = "vc2-2c-4gb"
    			label = "%s"
    			tag = "test23"
		}`, label)
}
