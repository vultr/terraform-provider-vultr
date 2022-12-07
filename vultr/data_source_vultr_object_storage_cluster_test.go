package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrObjectStorageCluster(t *testing.T) {
	t.Parallel()
	name := "data.vultr_object_storage_cluster.s3"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrObjectStorageCluster(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "hostname"),
					resource.TestCheckResourceAttrSet(name, "deploy"),
				),
			},
		},
	})
}

func testAccCheckVultrObjectStorageCluster() string {
	return `
		data "vultr_object_storage_cluster" "s3" {
			filter {
				name = "region"
				values = ["ewr"]
			}
		}`
}
