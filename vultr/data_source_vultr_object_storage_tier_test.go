package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrObjectStorageTier(t *testing.T) {
	t.Parallel()
	name := "data.vultr_object_storage_tier.obs_tier"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrObjectStorageTier(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "price"),
					resource.TestCheckResourceAttrSet(name, "locations.#"),
					resource.TestCheckResourceAttrSet(name, "slug"),
					resource.TestCheckResourceAttrSet(name, "rate_limit_bytes"),
					resource.TestCheckResourceAttrSet(name, "rate_limit_operations"),
				),
			},
		},
	})
}

func testAccCheckVultrObjectStorageTier() string {
	return `
		data "vultr_object_storage_tier" "obs_tier" {
			filter {
				name = "id"
				values = ["4"]
			}
		}`
}
