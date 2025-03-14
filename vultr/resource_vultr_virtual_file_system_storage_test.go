package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceVultrVirtualFileSystemStorage(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-vfs-rs")
	rServerLabel := acctest.RandomWithPrefix("tf-vps-vfs")
	rLabelUpdate := acctest.RandomWithPrefix("tf-vfs-rs-test-update")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrVirtualFileSystemStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVirtualFileSystemStorageConfig(rLabel, rServerLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrVirtualFileSystemStorageExists("vultr_virtual_file_system_storage.vfs"),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "size_gb", "40"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "region"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "tags.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "cost"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "charges"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "disk_type"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "status"),
				),
			},
			{
				Config: testAccVultrVirtualFileSystemStorageConfigAttach(rLabel, rServerLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrVirtualFileSystemStorageExists("vultr_virtual_file_system_storage.vfs"),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "size_gb", "40"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "region"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "attached_instances.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "tags.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "cost"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "charges"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "disk_type"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "status"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "attachments.#"),
				),
			},
			{
				Config: testAccVultrVirtualFileSystemStorageConfigUpdateLabel(rLabelUpdate, rServerLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrVirtualFileSystemStorageExists("vultr_virtual_file_system_storage.vfs"),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "label", rLabelUpdate),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "size_gb", "40"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "region"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "attached_instances.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "tags.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "cost"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "charges"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "disk_type"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "status"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "attachments.#"),
				),
			},
			{
				Config: testAccVultrVirtualFileSystemStorageConfigResize(rLabelUpdate, rServerLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrVirtualFileSystemStorageExists("vultr_virtual_file_system_storage.vfs"),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "size_gb", "45"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "region"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "attached_instances.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "tags.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "cost"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "charges"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "disk_type"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "status"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "attachments.#"),
				),
			},
			{
				// test detach by unsetting the attached_to_instance
				Config: testAccVultrVirtualFileSystemStorageConfigDetach(rLabelUpdate, rServerLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrVirtualFileSystemStorageExists("vultr_virtual_file_system_storage.vfs"),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "label", rLabelUpdate),
					resource.TestCheckResourceAttr("vultr_virtual_file_system_storage.vfs", "size_gb", "45"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "region"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "attached_instances.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "tags.#"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "cost"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "charges"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "disk_type"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "status"),
					resource.TestCheckResourceAttrSet("vultr_virtual_file_system_storage.vfs", "attachments.#"),
				),
			},
		},
	})
}

func testAccCheckVultrVirtualFileSystemStorageDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_virtual_file_system_storage" {
			continue
		}

		vfsID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		if _, _, err := client.VirtualFileSystemStorage.Get(context.Background(), vfsID); err == nil {
			return fmt.Errorf("vfs storage still exists: %s", vfsID)
		}
	}

	return nil
}

func testAccCheckVultrVirtualFileSystemStorageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vfs storage ID is not set")
		}

		vfsID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		if _, _, err := client.VirtualFileSystemStorage.Get(context.Background(), vfsID); err != nil {
			return fmt.Errorf("vfs storage does not exist: %s", vfsID)
		}

		return nil
	}
}

func testAccVultrVirtualFileSystemStorageConfig(label, serverLabel string) string {
	return fmt.Sprintf(`
	resource "vultr_virtual_file_system_storage" "vfs" {
		region = "ewr"
		size_gb = 40
		label = "%s"
		tags = ["terraform"]
	}

	resource "vultr_instance" "ip" {
		label = "%s"
		region = "ewr"
		plan = "vc2-4c-8gb"
		os_id = 1743
	}
  `, label, serverLabel)
}

func testAccVultrVirtualFileSystemStorageConfigAttach(label, serverLabel string) string {
	return fmt.Sprintf(`
	resource "vultr_virtual_file_system_storage" "vfs" {
		region = "ewr"
		size_gb = 40
		label = "%s"
		tags = ["terraform"]
		attached_instances = ["${vultr_instance.ip.id}"]
	}

	resource "vultr_instance" "ip" {
		label = "%s"
		region = "ewr"
		plan = "vc2-4c-8gb"
		os_id = 1743
	}
  `, label, serverLabel)
}

func testAccVultrVirtualFileSystemStorageConfigUpdateLabel(label, serverLabel string) string {
	return fmt.Sprintf(`
	resource "vultr_virtual_file_system_storage" "vfs" {
		region = "ewr"
		size_gb = 40
		label = "%s"
		tags = ["terraform"]
		attached_instances = ["${vultr_instance.ip.id}"]
	}

	resource "vultr_instance" "ip" {
		label = "%s"
		region = "ewr"
		plan = "vc2-4c-8gb"
		os_id = 1743
	}
  `, label, serverLabel)
}

func testAccVultrVirtualFileSystemStorageConfigResize(label, serverLabel string) string {
	return fmt.Sprintf(`
	resource "vultr_virtual_file_system_storage" "vfs" {
		region = "ewr"
		size_gb = 45
		label = "%s"
		tags = ["terraform"]
		attached_instances = ["${vultr_instance.ip.id}"]
	}

	resource "vultr_instance" "ip" {
		label = "%s"
		region = "ewr"
		plan = "vc2-4c-8gb"
		os_id = 1743
	}
  `, label, serverLabel)
}

func testAccVultrVirtualFileSystemStorageConfigDetach(label, serverLabel string) string {
	return fmt.Sprintf(`
	resource "vultr_virtual_file_system_storage" "vfs" {
		region = "ewr"
		size_gb = 45
		label = "%s"
		tags = ["terraform"]
	}

	resource "vultr_instance" "ip" {
		label = "%s"
		region = "ewr"
		plan = "vc2-4c-8gb"
		os_id = 1743
	}
  `, label, serverLabel)
}
