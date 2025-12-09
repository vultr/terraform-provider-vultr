package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVultrKubernetesNodePoolTaints(t *testing.T) {
	skipCI(t)
	rLabel := acctest.RandomWithPrefix("tf-vke-rs")
	rNP := acctest.RandomWithPrefix("tf-vke-np")

	name := "data.vultr_kubernetes_node_pool_taints.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel) +
					testAccVultrKubernetesNodePoolsBase(rNP) +
					testAccVultrKubernetesNodePoolTaintBase() +
					testAccDataSourceVultrKubernetesNodePoolTaints(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "cluster_id"),
					resource.TestCheckResourceAttrSet(name, "nodepool_id"),
					resource.TestCheckResourceAttr(name, "taints.#", "1"),
					resource.TestCheckResourceAttr(name, "taints.0.key", "workload"),
					resource.TestCheckResourceAttr(name, "taints.0.value", "gpu"),
					resource.TestCheckResourceAttr(name, "taints.0.effect", "NoSchedule"),
				),
			},
		},
	})
}

func testAccDataSourceVultrKubernetesNodePoolTaints() string {
	return `
		data "vultr_kubernetes_node_pool_taints" "test" {
			cluster_id  = vultr_kubernetes.foo.id
			nodepool_id = vultr_kubernetes_node_pools.foo.id

			depends_on = [vultr_kubernetes_node_pool_taint.test]
		}`
}
