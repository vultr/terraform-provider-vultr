package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrStartupScript_basic(t *testing.T) {
	rInt := acctest.RandInt()
	rName := fmt.Sprintf("foo-%d", rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrStartupScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrStartupScriptConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrStartupScriptExists("vultr_startup_script.foo"),
					resource.TestCheckResourceAttr("vultr_startup_script.foo", "name", rName),
					resource.TestCheckResourceAttr("vultr_startup_script.foo", "type", "pxe"),
					resource.TestCheckResourceAttrSet("vultr_startup_script.foo", "script"),
					resource.TestCheckResourceAttrSet("vultr_startup_script.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_startup_script.foo", "date_modified"),
				),
			},
			{
				Config: testAccVultrStartupScriptConfigUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrStartupScriptExists("vultr_startup_script.foo"),
					resource.TestCheckResourceAttr("vultr_startup_script.foo", "name", rName),
					resource.TestCheckResourceAttr("vultr_startup_script.foo", "type", "boot"),
					resource.TestCheckResourceAttrSet("vultr_startup_script.foo", "script"),
					resource.TestCheckResourceAttrSet("vultr_startup_script.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_startup_script.foo", "date_modified"),
				),
			},
		},
	})
}

func testAccCheckVultrStartupScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_startup_script" {
			continue
		}

		scriptID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		scripts, err := client.StartupScript.List(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting startup scripts: %s", err)
		}

		exists := false
		for i := range scripts {
			if scripts[i].ScriptID == scriptID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("Startup script still exists: %s", scriptID)
		}

		return nil
	}
	return nil
}

func testAccCheckVultrStartupScriptExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Script ID is not set")
		}

		scriptID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		scripts, err := client.StartupScript.List(context.Background())

		if err != nil {
			return fmt.Errorf("Error getting startup scripts: %s", err)
		}

		exists := false
		for i := range scripts {
			if scripts[i].ScriptID == scriptID {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Startup script not found ID: %s", scriptID)
		}

		return nil
	}
}

func testAccVultrStartupScriptConfigBasic(rName string) string {
	return fmt.Sprintf(`
		resource "vultr_startup_script" "foo" {
			name = "%s"
			type = "pxe"
			script = "#!/bin/bash\necho hello world > /root/hello"
		}
	`, rName)
}

func testAccVultrStartupScriptConfigUpdate(rName string) string {
	return fmt.Sprintf(`
		resource "vultr_startup_script" "foo" {
			name = "%s"
			type = "boot"
			script = "#!/bin/bash\necho hello world > /root/hello"
		}
	`, rName)
}
