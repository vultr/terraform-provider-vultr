package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVultrSnapshot(t *testing.T) {
	t.Parallel()
	rDesc := acctest.RandomWithPrefix("tf-test")
	rLabel := acctest.RandomWithPrefix("tf-test-vps")
	name := "data.vultr_snapshot.my_snapshot"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceVultrSnapshotBase(rLabel, rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", rDesc),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttrSet(name, "size"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "os_id"),
					resource.TestCheckResourceAttrSet(name, "app_id"),
				),
			},
			{
				Config:      testAccDataSourceVultrSnapshotConfig(rDesc),
				ExpectError: regexp.MustCompile(`.* data.vultr_snapshot.my_snapshot: data.vultr_snapshot.my_snapshot: no results were found`),
			},
		},
	})
}

func testAccDataSourceVultrSnapshotBase(vpsLabel, desc string) string {
	return fmt.Sprintf(`
	resource "vultr_server" "test" {
  			plan_id = "201"
  			region_id = "4"
			os_id = "167"
  			label = "%s"
  			hostname = "testing-the-hostname"
  			tag = "even better tag"
		}

		resource "vultr_snapshot" "foo" {
			vps_id       = "${vultr_server.test.id}"
			description  = "%s"
		}

		data "vultr_snapshot" "my_snapshot" {
    		filter {
    			name = "description"
    			values = ["${vultr_snapshot.foo.description}"]
			}
  		}
		`, vpsLabel, desc)
}

func testAccDataSourceVultrSnapshotConfig(description string) string {
	return fmt.Sprintf(`
		data "vultr_snapshot" "my_snapshot" {
    	filter {
    	name = "description"
    	values = ["%s"]
	}
  	}`, description)
}
