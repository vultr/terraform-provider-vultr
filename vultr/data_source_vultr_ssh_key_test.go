package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrSSHKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrSSHKey("Terraform Test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_ssh_key.my_key", "name", "Terraform Test"),
					resource.TestCheckResourceAttrSet("data.vultr_ssh_key.my_key", "ssh_key"),
					resource.TestCheckResourceAttrSet("data.vultr_ssh_key.my_key", "date_created"),
				),
			},
			{
				Config:      testAccCheckVultrSSHKey_noResult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_ssh_key.my_key: data.vultr_ssh_key.my_key: no results were found`),
			},
		},
	})
}

func testAccCheckVultrSSHKey(name string) string {
	return fmt.Sprintf(`
		data "vultr_ssh_key" "my_key" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}

func testAccCheckVultrSSHKey_noResult(name string) string {
	return fmt.Sprintf(`
		data "vultr_ssh_key" "my_key" {
    	filter {
    	name = "name"
    	values = ["%s"]
	}
  	}`, name)
}
