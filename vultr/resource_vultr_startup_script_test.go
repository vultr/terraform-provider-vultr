package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrStartupScriptBasic(t *testing.T) {
	rInt := acctest.RandInt()
	rName := fmt.Sprintf("foo-%d", rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrStartupScriptDestroy,
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

		_, err := client.StartupScript.Get(context.Background(), scriptID)
		if err == nil {
			return fmt.Errorf("startup script still exists: %s", scriptID)
		}

		return nil
	}
	return nil
}

func testAccCheckVultrStartupScriptExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("script ID is not set")
		}

		scriptID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, err := client.StartupScript.Get(context.Background(), scriptID)
		if err != nil {
			return fmt.Errorf("startup script not found ID: %s", scriptID)
		}

		return nil
	}
}

func testAccVultrStartupScriptConfigBasic(rName string) string {
	return fmt.Sprintf(`
		resource "vultr_startup_script" "foo" {
			name = "%s"
			type = "pxe"
			script = "IyEvYmluL2Jhc2hcbmVjaG8gaGVsbG8gd29ybGQgPiAvcm9vdC9oZWxsbw=="
		}
	`, rName)
}

func testAccVultrStartupScriptConfigUpdate(rName string) string {
	return fmt.Sprintf(`
		resource "vultr_startup_script" "foo" {
			name = "%s"
			type = "boot"
			script = "IyEvYmluL2Jhc2hcbmVjaG8gaGVsbG8gd29ybGQgPiAvcm9vdC9oZWxsbw=="
		}
	`, rName)
}
