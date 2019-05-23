package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVultrBlockStorage(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-test")
	rServerLabel := acctest.RandomWithPrefix("tf-vps-bs")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBlockStorageConfig(rLabel, rServerLabel),
			},
			{
				Config: testAccVultrBlockStorageConfig(rLabel, rServerLabel) + testAccDataSourceVultrBlockStorageConfig(rLabel),
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
				Config:      testAccDataSourceVultrBlockStorageConfig(rLabel),
				ExpectError: regexp.MustCompile(`.* data.vultr_block_storage.block: data.vultr_block_storage.block: no results were found`),
			},
		},
	})
}

func testAccDataSourceVultrBlockStorageConfig(label string) string {
	return fmt.Sprintf(`
	data "vultr_block_storage" "block" {
    	filter {
    		name = "label"
    		values = ["%s"]
  		}
	}`, label)
}
