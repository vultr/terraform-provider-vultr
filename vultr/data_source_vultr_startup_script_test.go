package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrStartupScript(t *testing.T) {

	rName := fmt.Sprintf("%s-terraform-test", acctest.RandString(10))
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
		data "vultr_startup_script" "my_script" {
    		filter {
    			name = "name"
    			values = ["${vultr_startup_script.foo.name}"]
			}
  		}

		resource "vultr_startup_script" "foo" {
			name = "%s"
			type = "pxe"
			script = "#!/bin/bash\necho hello world > /root/hello"
		}
		`, name)
}
