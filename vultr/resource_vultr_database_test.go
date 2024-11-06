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

func TestAccVultrDatabaseBasic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-db-rs")

	name := "vultr_database.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrDatabaseDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "cluster_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr(name, "database_engine", "pg"),
					resource.TestCheckResourceAttr(name, "database_engine_version", "15"),
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "maintenance_dow", "sunday"),
					resource.TestCheckResourceAttr(name, "maintenance_time", "01:00"),
					resource.TestCheckResourceAttr(name, "region", "SEA"),
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

func TestAccVultrDatabaseUpdate(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-db-rs-up")

	name := "vultr_database.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "cluster_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr(name, "database_engine", "pg"),
					resource.TestCheckResourceAttr(name, "database_engine_version", "15"),
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "maintenance_dow", "sunday"),
					resource.TestCheckResourceAttr(name, "maintenance_time", "01:00"),
					resource.TestCheckResourceAttr(name, "region", "SEA"),
					resource.TestCheckResourceAttr(name, "plan", "vultr-dbaas-startup-cc-1-55-2"),
					resource.TestCheckResourceAttr(name, "plan_disk", "55"),
					resource.TestCheckResourceAttr(name, "plan_ram", "2048"),
					resource.TestCheckResourceAttr(name, "plan_replicas", "0"),
					resource.TestCheckResourceAttr(name, "plan_vcpus", "1"),
					resource.TestCheckResourceAttr(name, "status", "Running"),
				),
			},
			{
				Config: testAccVultrDatabaseBaseUpdatedRegion(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "cluster_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr(name, "database_engine", "pg"),
					resource.TestCheckResourceAttr(name, "database_engine_version", "15"),
					resource.TestCheckResourceAttr(name, "label", rName),
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

func TestAccVultrDatabaseUpdatePlan(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-db-rs-upnid")

	name := "vultr_database.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "cluster_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr(name, "database_engine", "pg"),
					resource.TestCheckResourceAttr(name, "database_engine_version", "15"),
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "maintenance_dow", "sunday"),
					resource.TestCheckResourceAttr(name, "maintenance_time", "01:00"),
					resource.TestCheckResourceAttr(name, "region", "SEA"),
					resource.TestCheckResourceAttr(name, "plan", "vultr-dbaas-startup-cc-1-55-2"),
					resource.TestCheckResourceAttr(name, "plan_disk", "55"),
					resource.TestCheckResourceAttr(name, "plan_ram", "2048"),
					resource.TestCheckResourceAttr(name, "plan_replicas", "0"),
					resource.TestCheckResourceAttr(name, "plan_vcpus", "1"),
					resource.TestCheckResourceAttr(name, "status", "Running"),
				),
			},
			{
				Config: testAccVultrDatabaseBaseUpdatePlan(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "cluster_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr(name, "database_engine", "pg"),
					resource.TestCheckResourceAttr(name, "database_engine_version", "15"),
					resource.TestCheckResourceAttr(name, "label", rName),
					resource.TestCheckResourceAttr(name, "maintenance_dow", "sunday"),
					resource.TestCheckResourceAttr(name, "maintenance_time", "01:00"),
					resource.TestCheckResourceAttr(name, "region", "SEA"),
					resource.TestCheckResourceAttr(name, "plan", "vultr-dbaas-business-cc-1-55-2"),
					resource.TestCheckResourceAttr(name, "plan_disk", "55"),
					resource.TestCheckResourceAttr(name, "plan_ram", "2048"),
					resource.TestCheckResourceAttr(name, "plan_replicas", "1"),
					resource.TestCheckResourceAttr(name, "plan_vcpus", "1"),
					resource.TestCheckResourceAttr(name, "status", "Running"),
				),
			},
		},
	})
}

func testAccCheckVultrDatabaseDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_database" {
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

func testAccVultrDatabaseBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database" "test" {
			database_engine = "pg"
			database_engine_version = "15"
			region = "sea"
			plan = "vultr-dbaas-startup-cc-1-55-2"
			label = "%s"
			cluster_time_zone = "America/New_York"
			maintenance_dow = "sunday"
			maintenance_time = "01:00"
		} `, name)
}

func testAccVultrDatabaseBaseUpdatedRegion(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database" "test" {
			database_engine = "pg"
			database_engine_version = "15"
			region = "ewr"
			plan = "vultr-dbaas-startup-cc-1-55-2"
			label = "%s"
			cluster_time_zone = "America/New_York"
			maintenance_dow = "sunday"
			maintenance_time = "01:00"
		} `, name)
}

func testAccVultrDatabaseBaseUpdatePlan(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database" "test" {
			database_engine = "pg"
			database_engine_version = "15"
			region = "sea"
			plan = "vultr-dbaas-business-cc-1-55-2"
			label = "%s"
			cluster_time_zone = "America/New_York"
			maintenance_dow = "sunday"
			maintenance_time = "01:00"
		} `, name)
}

func testAccVultrDatabaseKafkaBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database" "test" {
			database_engine = "kafka"
			database_engine_version = "3.7"
			region = "sea"
			plan = "vultr-dbaas-startup-3x-occ-so-2-30-2"
			label = "%s"
			tag = "test tag"
		} `, name)
}
