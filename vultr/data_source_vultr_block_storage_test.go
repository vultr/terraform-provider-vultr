package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceVultrBlockStorage(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-bs-ds")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrBlockStorageConfig(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "cost"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "status"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "size_gb"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "region"),
					resource.TestCheckResourceAttrSet("data.vultr_block_storage.block", "label"),
				),
			},
		},
	})
}

func testAccDataSourceVultrBlockStorageConfig(label string) string {
	return fmt.Sprintf(`

	resource "vultr_block_storage" "foo" {
		region   = "ewr"
		size_gb     = 10
		label       = "%s"
	  }

	data "vultr_block_storage" "block" {
	filter {
		name = "label"
		values = ["${vultr_block_storage.foo.label}"]
		}
	}`, label)
}
