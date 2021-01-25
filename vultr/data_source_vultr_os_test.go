package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrOS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrOS("CentOS 7 x64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_os.centos", "name", "CentOS 7 x64"),
					resource.TestCheckResourceAttr(
						"data.vultr_os.centos", "arch", "x64"),
					resource.TestCheckResourceAttr(
						"data.vultr_os.centos", "id", "167"),
					resource.TestCheckResourceAttr(
						"data.vultr_os.centos", "family", "centos"),
				),
			},
		},
	})
}

func testAccCheckVultrOS(name string) string {
	return fmt.Sprintf(`data "vultr_os" "centos" {
 filter {
   name = "name"
   values = ["%s"]
 }
}`, name)
}
