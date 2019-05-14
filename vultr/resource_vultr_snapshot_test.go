package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrSnapshot_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrSnapshotConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrSnapshotExists("vultr_snapshot.foo"),
					resource.TestCheckResourceAttrSet("vultr_snapshot.foo", "vps_id"),
					resource.TestCheckResourceAttrSet("vultr_snapshot.foo", "description"),
					resource.TestCheckResourceAttrSet("vultr_snapshot.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_snapshot.foo", "size"),
					resource.TestCheckResourceAttrSet("vultr_snapshot.foo", "status"),
					resource.TestCheckResourceAttrSet("vultr_snapshot.foo", "os_id"),
					resource.TestCheckResourceAttrSet("vultr_snapshot.foo", "app_id"),
				),
			},
		},
	})
}

func testAccCheckVultrSnapshotDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_snapshot" {
			continue
		}

		snapshotID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		snapshots, err := client.Snapshot.GetList(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting snapshots: %s", err)
		}

		exists := false
		for i := range snapshots {
			if snapshots[i].SnapshotID == snapshotID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("Snapshot still exists: %s", snapshotID)
		}
	}
	return nil
}

func testAccCheckVultrSnapshotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Snapshot ID is not set")
		}

		snapshotID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		snapshots, err := client.Snapshot.GetList(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting snapshots: %s", err)
		}

		exists := false
		for i := range snapshots {
			if snapshots[i].SnapshotID == snapshotID {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Snapshot does not exist: %s", snapshotID)
		}

		return nil
	}
}

func testAccVultrSnapshotConfigBasic(rInt int) string {
	return fmt.Sprintf(`
		resource "vultr_snapshot" "foo" {
			vps_id       = "24609751"
			description  = "%d - created by terraform test"
		}
	`, rInt)
}
