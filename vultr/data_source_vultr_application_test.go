package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrApplication(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrApplication("Docker on Ubuntu 22.04"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "name", "Docker"),
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "short_name", "docker"),
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "type", "marketplace"),
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "vendor", "vultr-labs"),
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "image_id", "docker"),
				),
			},
		},
	})
}

func testAccCheckVultrApplication(deployName string) string {
	return fmt.Sprintf(`
		data "vultr_application" "docker" {
			filter {
				name = "deploy_name"
				values = ["%s"]
			}
		}`, deployName)
}
