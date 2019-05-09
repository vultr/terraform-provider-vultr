package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrStartupScript(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrStartupScript("Terraform Test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_startup_script.my_script", "name", "Terraform Test"),
					resource.TestCheckResourceAttr("data.vultr_startup_script.my_script", "type", "boot"),
					resource.TestCheckResourceAttrSet("data.vultr_startup_script.my_script", "script"),
					resource.TestCheckResourceAttrSet("data.vultr_startup_script.my_script", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_startup_script.my_script", "date_modified"),
				),
			},
			{
				Config:      testAccCheckVultrStartupScript_noResult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_startup_script.my_script: data.vultr_startup_script.my_script: no results were found`),
			},
			{
				Config:      testAccCheckVultrStartupScript_tooManyResults("boot"),
				ExpectError: regexp.MustCompile(`.* data.vultr_startup_script.my_script: data.vultr_startup_script.my_script: your search returned too many results. Please refine your search to be more specific`),
			},
		},
	})
}

func testAccCheckVultrStartupScript(name string) string {
	return fmt.Sprintf(`
		data "vultr_startup_script" "my_script" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}

func testAccCheckVultrStartupScript_noResult(name string) string {
	return fmt.Sprintf(`
		data "vultr_startup_script" "my_script" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}

func testAccCheckVultrStartupScript_tooManyResults(disk string) string {
	return fmt.Sprintf(`
		data "vultr_startup_script" "my_script" {
    	filter {
    	name = "type"
    	values = ["%s"]
	}
  	}`, disk)
}
