package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrIsoPublic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrIsoPublicRead("124 x64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_iso_public.finnix", "description", "124 x64"),
					resource.TestCheckResourceAttr(
						"data.vultr_iso_public.finnix", "name", "Finnix"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_public.finnix", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_iso_public.finnix", "md5sum"),
				),
			},
		},
	})
}

func testAccVultrIsoPublicRead(description string) string {
	return fmt.Sprintf(`
		data "vultr_iso_public" "finnix" {
			filter {
				name = "description"
				values = ["%s"]
			}
		}`, description)
}
