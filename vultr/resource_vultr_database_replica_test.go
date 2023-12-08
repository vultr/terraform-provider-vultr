package vultr

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrDatabaseReplicaBasic(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-replica")

	name := "vultr_database_replica.test_replica"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrDatabaseReplicaDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseReplicaBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "cluster_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr(name, "database_engine", "pg"),
					resource.TestCheckResourceAttr(name, "database_engine_version", "15"),
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "tag", "test tag"),
					resource.TestCheckResourceAttr(name, "maintenance_dow", "sunday"),
					resource.TestCheckResourceAttr(name, "maintenance_time", "01:00"),
					resource.TestCheckResourceAttr(name, "region", "EWR"),
					resource.TestCheckResourceAttr(name, "plan", "vultr-dbaas-startup-cc-1-55-2"),
					resource.TestCheckResourceAttr(name, "plan_disk", "55"),
					resource.TestCheckResourceAttr(name, "plan_ram", "2048"),
					resource.TestCheckResourceAttr(name, "plan_replicas", "0"),
					resource.TestCheckResourceAttr(name, "plan_vcpus", "1"),
					resource.TestCheckResourceAttr(name, "status", "Running"),
				),
			},
		},
	})
}

func TestAccVultrDatabaseReplicaUpdate(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-replica-up")

	name := "vultr_database_replica.test_replica"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDatabaseReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseReplicaBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "cluster_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr(name, "database_engine", "pg"),
					resource.TestCheckResourceAttr(name, "database_engine_version", "15"),
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "tag", "test tag"),
					resource.TestCheckResourceAttr(name, "maintenance_dow", "sunday"),
					resource.TestCheckResourceAttr(name, "maintenance_time", "01:00"),
					resource.TestCheckResourceAttr(name, "region", "EWR"),
					resource.TestCheckResourceAttr(name, "plan", "vultr-dbaas-startup-cc-1-55-2"),
					resource.TestCheckResourceAttr(name, "plan_disk", "55"),
					resource.TestCheckResourceAttr(name, "plan_ram", "2048"),
					resource.TestCheckResourceAttr(name, "plan_replicas", "0"),
					resource.TestCheckResourceAttr(name, "plan_vcpus", "1"),
					resource.TestCheckResourceAttr(name, "status", "Running"),
				),
			},
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseReplicaBaseUpdatedRegion(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "cluster_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr(name, "database_engine", "pg"),
					resource.TestCheckResourceAttr(name, "database_engine_version", "15"),
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "maintenance_dow", "sunday"),
					resource.TestCheckResourceAttr(name, "maintenance_time", "01:00"),
					resource.TestCheckResourceAttr(name, "region", "MIA"),
					resource.TestCheckResourceAttr(name, "plan", "vultr-dbaas-startup-cc-1-55-2"),
					resource.TestCheckResourceAttr(name, "plan_disk", "55"),
					resource.TestCheckResourceAttr(name, "plan_ram", "2048"),
					resource.TestCheckResourceAttr(name, "plan_replicas", "0"),
					resource.TestCheckResourceAttr(name, "plan_vcpus", "1"),
					resource.TestCheckResourceAttr(name, "status", "Running"),
				),
			},
		},
	})
}

func testAccCheckVultrDatabaseReplicaDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_database_replica" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.Database.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Not a valid Database Subscription UUID") {
				return nil
			}
			return fmt.Errorf("error getting database: %s", err)
		}

		return fmt.Errorf("database %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrDatabaseReplicaBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_replica" "test_replica" {
			database_id = vultr_database.test.id
			region = "ewr"
			label = "%s"
			tag = "test tag"
		} `, name)
}

func testAccVultrDatabaseReplicaBaseUpdatedRegion(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_replica" "test_replica" {
			database_id = vultr_database.test.id
			region = "mia"
			label = "%s"
			tag = "test tag"
		} `, name)
}
