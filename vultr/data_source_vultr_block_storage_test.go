package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrBlockStorage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBlockStorage_read("block-label"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "cost_pre_month"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "status"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "size_gb"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "region_id"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "label"),
				),
			},
			{
				Config:      testAccVultrBlockStorage_NoResult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_block_storage.block: data.vultr_block_storage.block: no results were found`),
			},
		},
	})
}

func testAccVultrBlockStorage_read(label string) string {
	return fmt.Sprintf(`
		data "vultr_block_storage" "block" {
    	filter {
    	name = "label"
    	values = ["%s"]
  		}
		}`, label)
}

func testAccVultrBlockStorage_NoResult(label string) string {
	return fmt.Sprintf(`
		data "vultr_block_storage" "block" {
    	filter {
    	name = "label"
    	values = ["%s"]
  		}
		}`, label)
}
