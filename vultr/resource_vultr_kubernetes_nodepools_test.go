package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestNodePoolCreateQuantity(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		nodeQuantity int
		autoScaler   bool
		minNodes     int
		want         int
		wantErr      bool
	}{
		{
			name:         "valid quantity with autoscaler",
			nodeQuantity: 1,
			autoScaler:   true,
			minNodes:     1,
			want:         1,
		},
		{
			name:         "fallback to min_nodes when quantity is zero",
			nodeQuantity: 0,
			autoScaler:   true,
			minNodes:     1,
			want:         1,
		},
		{
			name:         "fallback uses min_nodes value",
			nodeQuantity: 0,
			autoScaler:   true,
			minNodes:     2,
			want:         2,
		},
		{
			name:         "error when quantity is zero without autoscaler",
			nodeQuantity: 0,
			autoScaler:   false,
			minNodes:     1,
			wantErr:      true,
		},
		{
			name:         "error when quantity and min_nodes are zero",
			nodeQuantity: 0,
			autoScaler:   true,
			minNodes:     0,
			wantErr:      true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got, err := nodePoolCreateQuantity(c.nodeQuantity, c.autoScaler, c.minNodes)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected error, got quantity %d", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Fatalf("got quantity %d, want %d", got, c.want)
			}
		})
	}
}

func TestAccResourceVultrKubernetesNodePools(t *testing.T) {
	skipCI(t)
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

func TestAccResourceVultrKubernetesNodePoolsAutoScalerCreate(t *testing.T) {
	skipCI(t)
	rLabel := acctest.RandomWithPrefix("tf-vke-rs")
	rNP := acctest.RandomWithPrefix("tf-vke-np")

	name := "vultr_kubernetes_node_pools.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel) + testAccVultrKubernetesNodePoolsAutoScalerCreate(rNP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rNP),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttr(name, "node_quantity", "1"),
					resource.TestCheckResourceAttr(name, "nodes.#", "1"),
					resource.TestCheckResourceAttr(name, "plan", "vc2-2c-4gb"),
					resource.TestCheckResourceAttr(name, "auto_scaler", "true"),
					resource.TestCheckResourceAttr(name, "min_nodes", "1"),
					resource.TestCheckResourceAttr(name, "max_nodes", "2"),
				),
			},
		},
	})
}

func TestAccResourceVultrKubernetesNodePoolsUpdate(t *testing.T) {
	skipCI(t)
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

func TestAccResourceVultrKubernetesNodePoolsUpdateAutoScaler(t *testing.T) {
	skipCI(t)
	rLabel := acctest.RandomWithPrefix("tf-vke-rs")
	rNP := acctest.RandomWithPrefix("tf-vke-np")

	name := "vultr_kubernetes_node_pools.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrKubernetesBase(rLabel) + testAccVultrKubernetesNodePoolsUpdate(rNP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rNP),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "tag"),
					resource.TestCheckResourceAttr(name, "nodes.#", "2"),
					resource.TestCheckResourceAttr(name, "plan", "vc2-2c-4gb"),
					resource.TestCheckResourceAttr(name, "auto_scaler", "true"),
					resource.TestCheckResourceAttr(name, "min_nodes", "2"),
					resource.TestCheckResourceAttr(name, "max_nodes", "4"),
					resource.TestCheckResourceAttr(name, "labels.0.key", "test-label"),
					resource.TestCheckResourceAttr(name, "labels.0.value", "test-label-value-upd"),
					resource.TestCheckResourceAttr(name, "taints.0.key", "test-taint"),
					resource.TestCheckResourceAttr(name, "taints.0.value", "test-taint-value-upd"),
					resource.TestCheckResourceAttr(name, "taints.0.effect", "PreferNoSchedule"),
				),
			},
			{
				Config: testAccVultrKubernetesBase(rLabel) + testAccVultrKubernetesNodePoolsUpdateAutoScaler(rNP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rNP),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "tag"),
					resource.TestCheckResourceAttr(name, "nodes.#", "2"),
					resource.TestCheckResourceAttr(name, "plan", "vc2-2c-4gb"),
					resource.TestCheckResourceAttr(name, "auto_scaler", "false"),
					resource.TestCheckResourceAttr(name, "min_nodes", "3"),
					resource.TestCheckResourceAttr(name, "max_nodes", "5"),
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

				labels {
					key = "test-label"
					value = "test-label-value"
				}

				taints {
					key = "test-taint"
					value = "test-taint-value"
					effect = "PreferNoSchedule"
				}
		}`, label)
}

func testAccVultrKubernetesNodePoolsAutoScalerCreate(label string) string {
	return fmt.Sprintf(`
		resource "vultr_kubernetes_node_pools" "foo" {
    			cluster_id = vultr_kubernetes.foo.id
				node_quantity = 1
				plan = "vc2-2c-4gb"
    			label = "%s"
    			tag = "test23"
				auto_scaler = true
				min_nodes = 1
				max_nodes = 2
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
				auto_scaler = true
				min_nodes = 2
				max_nodes = 4

				labels {
					key = "test-label"
					value = "test-label-value-upd"
				}

				taints {
					key = "test-taint"
					value = "test-taint-value-upd"
					effect = "PreferNoSchedule"
				}
		}`, label)
}

func testAccVultrKubernetesNodePoolsUpdateAutoScaler(label string) string {
	return fmt.Sprintf(`
		resource "vultr_kubernetes_node_pools" "foo" {
    			cluster_id = vultr_kubernetes.foo.id
				node_quantity = 2
				plan = "vc2-2c-4gb"
    			label = "%s"
    			tag = "test23"
				auto_scaler = false
				min_nodes = 3
				max_nodes = 5

				labels {
					key = "test-label"
					value = "test-label-value"
				}

				taints {
					key = "test-taint"
					value = "test-taint-value"
					effect = "PreferNoSchedule"
				}
		}`, label)
}
