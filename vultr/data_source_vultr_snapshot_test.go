package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrSnapshot(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrSnapshot("Terraform Test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_snapshot.my_snapshot", "description", "Terraform Test"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "size"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "status"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "os_id"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "app_id"),
				),
			},
			{
				Config:      testAccCheckVultrSnapshot_noResult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_snapshot.my_snapshot: data.vultr_snapshot.my_snapshot: no results were found`),
			},
		},
	})
}

func testAccCheckVultrSnapshot(description string) string {
	return fmt.Sprintf(`
		data "vultr_snapshot" "my_snapshot" {
    	filter {
    	name = "description"
    	values = ["%s"]
	}
  	}`, description)
}

func testAccCheckVultrSnapshot_noResult(description string) string {
	return fmt.Sprintf(`
		data "vultr_snapshot" "my_snapshot" {
    	filter {
    	name = "description"
    	values = ["%s"]
	}
  	}`, description)
}
