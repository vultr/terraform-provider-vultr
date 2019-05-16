package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceVultrBlockStorage(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-test")
	rLabelUpdate := acctest.RandomWithPrefix("tf-test-update")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBlockStorageConfig(rLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBlockStorageExists("vultr_block_storage.foo"),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "size_gb", "10"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "cost_per_month"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "status"),
				),
			},
			{
				Config: testAccVultrBlockStorageConfig_attach(rLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBlockStorageExists("vultr_block_storage.foo"),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "size_gb", "10"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "cost_per_month"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "status"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "attached_id"),
				),
			},
			{
				Config: testAccVultrBlockStorageConfig_updateLabel(rLabelUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBlockStorageExists("vultr_block_storage.foo"),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "label", rLabelUpdate),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "size_gb", "10"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "cost_per_month"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "status"),
				),
			},
			{
				Config: testAccVultrBlockStorageConfig_resize(rLabelUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBlockStorageExists("vultr_block_storage.foo"),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "label", rLabelUpdate),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "size_gb", "15"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "cost_per_month"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "status"),
				),
			},
			{
				// test detach by unsetting the attached_id
				Config: testAccVultrBlockStorageConfig_detach(rLabelUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBlockStorageExists("vultr_block_storage.foo"),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "label", rLabelUpdate),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "size_gb", "15"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "cost_per_month"),
					resource.TestCheckResourceAttrSet("vultr_block_storage.foo", "status"),
				),
			},
		},
	})
}

func testAccCheckVultrBlockStorageDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_block_storage" {
			continue
		}

		bsID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		bses, err := client.BlockStorage.GetList(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting block storages: %s", err)
		}

		exists := false
		for i := range bses {
			if bses[i].BlockStorageID == bsID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("Block storage still exists: %s", bsID)
		}
	}
	return nil
}

func testAccCheckVultrBlockStorageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Block storage ID is not set")
		}

		bsID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		bses, err := client.BlockStorage.GetList(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting block storages: %s", err)
		}

		exists := false
		for i := range bses {
			if bses[i].BlockStorageID == bsID {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Block storage does not exist: %s", bsID)
		}

		return nil
	}
}

func testAccVultrBlockStorageConfig(label string) string {
	return fmt.Sprintf(`
	data "vultr_region" "new_jersey" {
		filter {
		  name   = "name"
		  values = ["New Jersey"]
		}
	}
	resource "vultr_block_storage" "foo" {
		region_id   = "${data.vultr_region.new_jersey.id}"
		size_gb     = 10
		label       = "%s"
	  }
   `, label)
}

func testAccVultrBlockStorageConfig_attach(label string) string {
	return fmt.Sprintf(`
	data "vultr_region" "new_jersey" {
		filter {
		  name   = "name"
		  values = ["New Jersey"]
		}
	}
	resource "vultr_block_storage" "foo" {
		region_id   = "${data.vultr_region.new_jersey.id}"
		size_gb     = 10
		label       = "%s"
		attached_id = "24609751"
	  }
   `, label)
}

func testAccVultrBlockStorageConfig_updateLabel(label string) string {
	return fmt.Sprintf(`
	data "vultr_region" "new_jersey" {
		filter {
		  name   = "name"
		  values = ["New Jersey"]
		}
	}
	resource "vultr_block_storage" "foo" {
		region_id   = "${data.vultr_region.new_jersey.id}"
		size_gb     = 10
		label       = "%s"
		attached_id = "24609751"
	  }
   `, label)
}

func testAccVultrBlockStorageConfig_resize(label string) string {
	return fmt.Sprintf(`
	data "vultr_region" "new_jersey" {
		filter {
		  name   = "name"
		  values = ["New Jersey"]
		}
	}
	resource "vultr_block_storage" "foo" {
		region_id   = "${data.vultr_region.new_jersey.id}"
		size_gb     = 15
		label       = "%s"
		attached_id = "24609751"
	  }
   `, label)
}

func testAccVultrBlockStorageConfig_detach(label string) string {
	return fmt.Sprintf(`
	data "vultr_region" "new_jersey" {
		filter {
		  name   = "name"
		  values = ["New Jersey"]
		}
	}
	resource "vultr_block_storage" "foo" {
		region_id   = "${data.vultr_region.new_jersey.id}"
		size_gb     = 15
		label       = "%s"
	  }
   `, label)
}
