package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrOS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrOS("CentOS 6 x64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_os.centos", "name", "CentOS 6 x64"),
					resource.TestCheckResourceAttr(
						"data.vultr_os.centos", "arch", "x64"),
					resource.TestCheckResourceAttr(
						"data.vultr_os.centos", "windows", "false"),
					resource.TestCheckResourceAttr(
						"data.vultr_os.centos", "id", "127"),
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
