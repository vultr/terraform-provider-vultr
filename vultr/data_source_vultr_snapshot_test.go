package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVultrSnapshot(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrSnapshotConfigBasic(rDesc),
			},
			{
				Config: testAccVultrSnapshotConfigBasic(rDesc) + testAccDataSourceVultrSnapshotConfig(rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vultr_snapshot.my_snapshot", "description", rDesc),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "size"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "status"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "os_id"),
					resource.TestCheckResourceAttrSet("data.vultr_snapshot.my_snapshot", "app_id"),
				),
			},
			{
				Config:      testAccDataSourceVultrSnapshotConfig(rDesc),
				ExpectError: regexp.MustCompile(`.* data.vultr_snapshot.my_snapshot: data.vultr_snapshot.my_snapshot: no results were found`),
			},
		},
	})
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
