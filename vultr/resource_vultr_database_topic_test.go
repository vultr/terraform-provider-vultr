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

func TestAccVultrDatabaseTopicBasic(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-topic")

	name := "vultr_database_topic.test_topic"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrDatabaseTopicDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseKafkaBase(pName) + testAccVultrDatabaseTopicBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "partitions", "3"),
					resource.TestCheckResourceAttr(name, "replication", "2"),
					resource.TestCheckResourceAttr(name, "retention_hours", "150"),
					resource.TestCheckResourceAttr(name, "retention_bytes", "150000"),
				),
			},
		},
	})
}

func TestAccVultrDatabaseTopicUpdate(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-topic-up")

	name := "vultr_database_topic.test_topic"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDatabaseTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseKafkaBase(pName) + testAccVultrDatabaseTopicBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "partitions", "3"),
					resource.TestCheckResourceAttr(name, "replication", "2"),
					resource.TestCheckResourceAttr(name, "retention_hours", "150"),
					resource.TestCheckResourceAttr(name, "retention_bytes", "150000"),
				),
			},
			{
				PreConfig: func() { time.Sleep(60 * time.Second) },
				Config:    testAccVultrDatabaseKafkaBase(pName) + testAccVultrDatabaseTopicBaseUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "partitions", "3"),
					resource.TestCheckResourceAttr(name, "replication", "2"),
					resource.TestCheckResourceAttr(name, "retention_hours", "160"),
					resource.TestCheckResourceAttr(name, "retention_bytes", "160000"),
				),
			},
		},
	})
}

func testAccCheckVultrDatabaseTopicDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_database_topic" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.Database.GetTopic(context.Background(), rs.Primary.Attributes["database_id"], rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Not a valid database topic") || strings.Contains(err.Error(), "Not a valid Database Subscription UUID") {
				return nil
			}
			return fmt.Errorf("error getting database topic: %s", err)
		}

		return fmt.Errorf("database topic %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrDatabaseTopicBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_topic" "test_topic" {
			database_id = vultr_database.test.id
			name = "%s"
			partitions = "3"
			replication = "2"
			retention_hours = "150"
			retention_bytes = "150000"
		} `, name)
}

func testAccVultrDatabaseTopicBaseUpdated(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_topic" "test_topic" {
			database_id = vultr_database.test.id
			name = "%s"
			partitions = "3"
			replication = "2"
			retention_hours = "160"
			retention_bytes = "160000"
		} `, name)
}
