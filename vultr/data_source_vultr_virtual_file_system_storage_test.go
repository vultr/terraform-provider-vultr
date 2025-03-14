package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVultrVirtualFileSystemStorage(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-vfs-ds")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVultrVirtualFileSystemStorageConfig(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_virtual_file_system_storage.vfs", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_virtual_file_system_storage.vfs", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_virtual_file_system_storage.vfs", "cost"),
					resource.TestCheckResourceAttrSet("data.vultr_virtual_file_system_storage.vfs", "charges"),
					resource.TestCheckResourceAttrSet("data.vultr_virtual_file_system_storage.vfs", "status"),
					resource.TestCheckResourceAttrSet("data.vultr_virtual_file_system_storage.vfs", "size_gb"),
					resource.TestCheckResourceAttrSet("data.vultr_virtual_file_system_storage.vfs", "region"),
					resource.TestCheckResourceAttrSet("data.vultr_virtual_file_system_storage.vfs", "label"),
				),
			},
		},
	})
}

func testAccDataSourceVultrVirtualFileSystemStorageConfig(label string) string {
	return fmt.Sprintf(`

	resource "vultr_virtual_file_system_storage" "vfs" {
		region   = "ewr"
		size_gb     = 40
		label       = "%s"
	  }

	data "vultr_virtual_file_system_storage" "vfs-ds" {
	filter {
		name = "label"
		values = ["${vultr_virtual_file_system_storage.vfs.label}"]
		}
	}`, label)
}
