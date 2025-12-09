package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceVultrKubernetesNodePoolLabel(t *testing.T) {
	skipCI(t)
	rLabel := acctest.RandomWithPrefix("tf-vke-rs")
	rNP := acctest.RandomWithPrefix("tf-vke-np")

	name := "vultr_kubernetes_node_pool_label.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel) + testAccVultrKubernetesNodePoolsBase(rNP) + testAccVultrKubernetesNodePoolLabelBase(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "environment"),
					resource.TestCheckResourceAttr(name, "value", "production"),
					resource.TestCheckResourceAttrSet(name, "cluster_id"),
					resource.TestCheckResourceAttrSet(name, "nodepool_id"),
				),
			},
		},
	})
}

func testAccVultrKubernetesNodePoolLabelBase() string {
	return `
		resource "vultr_kubernetes_node_pool_label" "test" {
			cluster_id  = vultr_kubernetes.foo.id
			nodepool_id = vultr_kubernetes_node_pools.foo.id
			key         = "environment"
			value       = "production"
		}`
}
