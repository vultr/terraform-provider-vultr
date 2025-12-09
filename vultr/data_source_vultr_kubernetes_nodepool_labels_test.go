package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVultrKubernetesNodePoolLabels(t *testing.T) {
	skipCI(t)
	rLabel := acctest.RandomWithPrefix("tf-vke-rs")
	rNP := acctest.RandomWithPrefix("tf-vke-np")

	name := "data.vultr_kubernetes_node_pool_labels.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel) +
					testAccVultrKubernetesNodePoolsBase(rNP) +
					testAccVultrKubernetesNodePoolLabelBase() +
					testAccDataSourceVultrKubernetesNodePoolLabels(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "cluster_id"),
					resource.TestCheckResourceAttrSet(name, "nodepool_id"),
					resource.TestCheckResourceAttr(name, "labels.#", "1"),
					resource.TestCheckResourceAttr(name, "labels.0.key", "environment"),
					resource.TestCheckResourceAttr(name, "labels.0.value", "production"),
				),
			},
		},
	})
}

func testAccDataSourceVultrKubernetesNodePoolLabels() string {
	return `
		data "vultr_kubernetes_node_pool_labels" "test" {
			cluster_id  = vultr_kubernetes.foo.id
			nodepool_id = vultr_kubernetes_node_pools.foo.id

			depends_on = [vultr_kubernetes_node_pool_label.test]
		}`
}
