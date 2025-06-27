package vultr

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrDatabaseConnectorBasic(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-connector")

	name := "vultr_database_connector.test_connector"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrDatabaseConnectorDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseKafkaBase(pName) + testAccVultrDatabaseConnectorBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "class", "com.couchbase.connect.kafka.CouchbaseSinkConnector"),
					resource.TestCheckResourceAttr(name, "topics", "tf-db-topic"),
					resource.TestCheckResourceAttr(name, "config", "{\"couchbase.seed.nodes\":\"3\",\"couchbase.username\":\"some_username\",\"couchbase.password\":\"some_password\"}"),
				),
			},
		},
	})
}

func TestAccVultrDatabaseConnectorUpdate(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-connector-up")

	name := "vultr_database_connector.test_connector"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDatabaseConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseKafkaBase(pName) + testAccVultrDatabaseConnectorBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "class", "com.couchbase.connect.kafka.CouchbaseSinkConnector"),
					resource.TestCheckResourceAttr(name, "topics", "tf-db-topic"),
					resource.TestCheckResourceAttr(name, "config", "{\"couchbase.seed.nodes\":\"3\",\"couchbase.username\":\"some_username\",\"couchbase.password\":\"some_password\"}"),
				),
			},
			{
				PreConfig: func() { time.Sleep(60 * time.Second) },
				Config:    testAccVultrDatabaseKafkaBase(pName) + testAccVultrDatabaseConnectorBaseUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "class", "com.couchbase.connect.kafka.CouchbaseSinkConnector"),
					resource.TestCheckResourceAttr(name, "topics", "tf-db-topic-2"),
					resource.TestCheckResourceAttr(name, "config", "{\"couchbase.seed.nodes\":\"3\",\"couchbase.username\":\"some_username\",\"couchbase.password\":\"some_password\"}"),
				),
			},
		},
	})
}

func testAccCheckVultrDatabaseConnectorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_database_connector" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.Database.GetConnector(context.Background(), rs.Primary.Attributes["database_id"], rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Not a valid database connector") || strings.Contains(err.Error(), "Not a valid Database Subscription UUID") {
				return nil
			}
			return fmt.Errorf("error getting database connector: %s", err)
		}

		return fmt.Errorf("database connector %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrDatabaseConnectorBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_connector" "test_connector" {
			database_id = vultr_database.test.id
			name = "%s"
			class = "com.couchbase.connect.kafka.CouchbaseSinkConnector"
			topics = "tf-db-topic"
			config = jsonencode({
				"couchbase.seed.nodes" = "3"
				"couchbase.username" = "some_username"
				"couchbase.password" = "some_password"
			})
		} `, name)
}

func testAccVultrDatabaseConnectorBaseUpdated(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_connector" "test_connector" {
			database_id = vultr_database.test.id
			name = "%s"
			class = "com.couchbase.connect.kafka.CouchbaseSinkConnector"
			topics = "tf-db-topic-2"
			config = jsonencode({
				"couchbase.seed.nodes" = "3"
				"couchbase.username" = "some_username"
				"couchbase.password" = "some_password"
			})
		} `, name)
}
