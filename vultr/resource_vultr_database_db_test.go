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

func TestAccVultrDatabaseDBBasic(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-db")

	name := "vultr_database_db.test_db"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrDatabaseDBDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseDBBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
				),
			},
		},
	})
}

func testAccCheckVultrDatabaseDBDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_database_db" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.Database.GetDB(context.Background(), rs.Primary.Attributes["database_id"], rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Not a valid database db") || strings.Contains(err.Error(), "NNot a valid Database Subscription UUID") {
				return nil
			}
			return fmt.Errorf("error getting database db: %s", err)
		}

		return fmt.Errorf("database db %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrDatabaseDBBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_db" "test_db" {
			database_id = vultr_database.test.id
			name = "%s"
		} `, name)
}
