package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrApi(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrApi(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_api_key.api", "acl.#"),
					resource.TestCheckResourceAttrSet("data.vultr_api_key.api", "email"),
					resource.TestCheckResourceAttrSet("data.vultr_api_key.api", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_api_key.api", "name"),
				),
			},
		},
	})
}

func testAccVultrApi() string {
	return `data "vultr_api_key" "api" {}`
}
