package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceVultrKubernetesNodePoolTaint(t *testing.T) {
	skipCI(t)
	rLabel := acctest.RandomWithPrefix("tf-vke-rs")
	rNP := acctest.RandomWithPrefix("tf-vke-np")

	name := "vultr_kubernetes_node_pool_taint.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel) + testAccVultrKubernetesNodePoolsBase(rNP) + testAccVultrKubernetesNodePoolTaintBase(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "workload"),
					resource.TestCheckResourceAttr(name, "value", "gpu"),
					resource.TestCheckResourceAttr(name, "effect", "NoSchedule"),
					resource.TestCheckResourceAttrSet(name, "cluster_id"),
					resource.TestCheckResourceAttrSet(name, "nodepool_id"),
				),
			},
		},
	})
}

func testAccVultrKubernetesNodePoolTaintBase() string {
	return `
		resource "vultr_kubernetes_node_pool_taint" "test" {
			cluster_id  = vultr_kubernetes.foo.id
			nodepool_id = vultr_kubernetes_node_pools.foo.id
			key         = "workload"
			value       = "gpu"
			effect      = "NoSchedule"
		}`
}
