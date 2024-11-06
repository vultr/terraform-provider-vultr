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

func TestAccVultrDatabaseQuotaBasic(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	uName := acctest.RandomWithPrefix("tf-db-user")
	rName := acctest.RandomWithPrefix("tf-db-quota")

	name := "vultr_database_quota.test_quota"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrDatabaseQuotaDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseKafkaBase(pName) + testAccVultrDatabaseUserKafkaBase(uName) + testAccVultrDatabaseQuotaBase(rName, uName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "client_id", rName),
					resource.TestCheckResourceAttr(name, "consumer_byte_rate", "12345"),
					resource.TestCheckResourceAttr(name, "producer_byte_rate", "23456"),
					resource.TestCheckResourceAttr(name, "request_percentage", "20"),
					resource.TestCheckResourceAttr(name, "user", uName),
				),
			},
		},
	})
}

func testAccCheckVultrDatabaseQuotaDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_database_quota" {
			continue
		}

		quotaID := strings.Split(rs.Primary.ID, "|")

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.Database.GetQuota(context.Background(), rs.Primary.Attributes["database_id"], quotaID[0], quotaID[1])
		if err != nil {
			if strings.Contains(err.Error(), "Not a valid database quota") || strings.Contains(err.Error(), "Not a valid Database Subscription UUID") {
				return nil
			}
			return fmt.Errorf("error getting database quota: %s", err)
		}

		return fmt.Errorf("database quota %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrDatabaseQuotaBase(clientID string, user string) string {
	return fmt.Sprintf(`
		resource "vultr_database_quota" "test_quota" {
			database_id = vultr_database.test.id
			client_id = "%s"
			consumer_byte_rate = "12345"
			producer_byte_rate = "23456"
			request_percentage = "20"
			user = "%s"
		} `, clientID, user)
}
