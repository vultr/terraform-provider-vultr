package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrApplication(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrApplication("Docker on CentOS 7 x64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "deploy_name", "Docker on CentOS 7 x64"),
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "name", "Docker"),
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "short_name", "docker"),
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "id", "17"),
					resource.TestCheckResourceAttr(
						"data.vultr_application.docker", "surcharge", "0"),
				),
			},
			{
				Config:      testAccCheckVultrApplication_noresult("image_test"),
				ExpectError: regexp.MustCompile(`.* data.vultr_application.docker: data.vultr_application.docker: no results where found`),
			},
			{
				Config:      testAccCheckVultrApplication_tooManyResults("Docker"),
				ExpectError: regexp.MustCompile(`.* data.vultr_application.toomany: data.vultr_application.toomany: your search returned too many results : 2. Please refine your search to be more specific`),
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

func testAccCheckVultrApplication_noresult(name string) string {
	return fmt.Sprintf(`
		data "vultr_application" "docker" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}

func testAccCheckVultrApplication_tooManyResults(name string) string {
	return fmt.Sprintf(`
		data "vultr_application" "toomany" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}
