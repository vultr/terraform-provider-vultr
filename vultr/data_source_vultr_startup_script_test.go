package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrStartupScript(t *testing.T) {

	rName := acctest.RandomWithPrefix("tf-startup-ds")
	name := "data.vultr_startup_script.my_script"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrStartupScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrStartupScriptConfig_base(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "type", "pxe"),
					resource.TestCheckResourceAttrSet(name, "script"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttrSet(name, "date_modified"),
				),
			},
		},
	})
}

func testAccCheckVultrStartupScriptConfig_base(name string) string {
	return fmt.Sprintf(`
		resource "vultr_startup_script" "foo" {
			name = "%s"
			type = "pxe"
			script = "IyEvYmluL2Jhc2hcbmVjaG8gaGVsbG8gd29ybGQgPiAvcm9vdC9oZWxsbw=="
		}

		data "vultr_startup_script" "my_script" {
		filter {
			name = "name"
			values = ["${vultr_startup_script.foo.name}"]
			}
		}
		`, name)
}
