package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrKubernetes(t *testing.T) {
	skipCI(t)

	rLabel := acctest.RandomWithPrefix("tf-test-k8")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrKubernetes(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_kubernetes.k8", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_kubernetes.k8", "status"),
					resource.TestCheckResourceAttrSet("data.vultr_kubernetes.k8", "region"),
					resource.TestCheckResourceAttrSet("data.vultr_kubernetes.k8", "label"),
					resource.TestCheckResourceAttrSet("data.vultr_kubernetes.k8", "kube_config"),
					resource.TestCheckResourceAttr("data.vultr_kubernetes.k8", "node_pools.#", "1"),
				),
			},
		},
	})
}

func testAccCheckVultrKubernetes(label string) string {
	return fmt.Sprintf(`
		resource "vultr_kubernetes" "test" {
			region = "ewr"
			label = "%s"
			version = "v1.24.3+2"

			node_pools {
				node_quantity = 1
				plan = "vc2-2c-4gb"
    			label = "tf-test-label"
			}
		}

		data "vultr_kubernetes" "k8" {
			filter {
				name = "label"
				values = ["${vultr_kubernetes.test.label}"]
			}
		}`, label)
}
