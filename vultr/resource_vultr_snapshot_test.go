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
	t.Parallel()
	rInt := acctest.RandInt()
	desc := fmt.Sprintf("%d - created by Terraform test", rInt)
	rServerLabel := acctest.RandomWithPrefix("tf-vps-snap")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrSnapshotConfigBasic(rServerLabel, desc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrSnapshotExists("vultr_snapshot.foo"),
					resource.TestCheckResourceAttrSet("vultr_snapshot.foo", "vps_id"),
					resource.TestCheckResourceAttr("vultr_snapshot.foo", "description", desc),
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

func testAccVultrSnapshotConfigBasic(rServerLabel, desc string) string {
	return fmt.Sprintf(`
		resource "vultr_server" "snap" {
        	label = "%s"
        	region_id = "1"
        	plan_id = 201
        	os_id = 147
    	}
		resource "vultr_snapshot" "foo" {
			vps_id       = "${vultr_server.snap.id}"
			description  = "%s"
		}
	`, rServerLabel, desc)
}
