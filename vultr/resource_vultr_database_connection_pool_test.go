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

func TestAccVultrDatabaseConnectionPoolBasic(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-connection-pool")

	name := "vultr_database_connection_pool.test_connection_pool"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrDatabaseConnectionPoolDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseConnectionPoolBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "database", "defaultdb"),
					resource.TestCheckResourceAttr(name, "username", "vultradmin"),
					resource.TestCheckResourceAttr(name, "mode", "transaction"),
					resource.TestCheckResourceAttr(name, "size", "3"),
				),
			},
		},
	})
}

func TestAccVultrDatabaseConnectionPoolUpdate(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-connection-pool-up")

	name := "vultr_database_connection_pool.test_connection_pool"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDatabaseConnectionPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseConnectionPoolBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "database", "defaultdb"),
					resource.TestCheckResourceAttr(name, "username", "vultradmin"),
					resource.TestCheckResourceAttr(name, "mode", "transaction"),
					resource.TestCheckResourceAttr(name, "size", "3"),
				),
			},
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseConnectionPoolBaseUpdatedMode(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "database", "defaultdb"),
					resource.TestCheckResourceAttr(name, "username", "vultradmin"),
					resource.TestCheckResourceAttr(name, "mode", "session"),
					resource.TestCheckResourceAttr(name, "size", "3"),
				),
			},
		},
	})
}
func TestAccVultrDatabaseConnectionPoolUpdateSize(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-connection-pool-up")

	name := "vultr_database_connection_pool.test_connection_pool"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDatabaseConnectionPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseConnectionPoolBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "database", "defaultdb"),
					resource.TestCheckResourceAttr(name, "username", "vultradmin"),
					resource.TestCheckResourceAttr(name, "mode", "transaction"),
					resource.TestCheckResourceAttr(name, "size", "3"),
				),
			},
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseConnectionPoolBaseUpdatedSize(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "database", "defaultdb"),
					resource.TestCheckResourceAttr(name, "username", "vultradmin"),
					resource.TestCheckResourceAttr(name, "mode", "transaction"),
					resource.TestCheckResourceAttr(name, "size", "5"),
				),
			},
		},
	})
}

func testAccCheckVultrDatabaseConnectionPoolDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_database_connection_pool" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.Database.GetConnectionPool(context.Background(), rs.Primary.Attributes["database_id"], rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Not a valid connection pool") || strings.Contains(err.Error(), "Not a valid DBaaS Subscription UUID") {
				return nil
			}
			return fmt.Errorf("error getting database connection pool: %s", err)
		}

		return fmt.Errorf("database connection pool %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrDatabaseConnectionPoolBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_connection_pool" "test_connection_pool" {
			database_id = vultr_database.test.id
			name = "%s"
			database = "defaultdb"
			username = "vultradmin"
			mode = "transaction"
			size = "3"
		} `, name)
}

func testAccVultrDatabaseConnectionPoolBaseUpdatedMode(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_connection_pool" "test_connection_pool" {
			database_id = vultr_database.test.id
			name = "%s"
			database = "defaultdb"
			username = "vultradmin"
			mode = "session"
			size = "3"
		} `, name)
}

func testAccVultrDatabaseConnectionPoolBaseUpdatedSize(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_connection_pool" "test_connection_pool" {
			database_id = vultr_database.test.id
			name = "%s"
			database = "defaultdb"
			username = "vultradmin"
			mode = "transaction"
			size = "5"
		} `, name)
}
