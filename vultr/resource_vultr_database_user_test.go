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

func TestAccVultrDatabaseUserBasic(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-user")

	name := "vultr_database_user.test_user"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrDatabaseUserDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseUserBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "username", rName),
					resource.TestCheckResourceAttr(name, "password", "someRandomPW4928!z"),
				),
			},
		},
	})
}

func TestAccVultrDatabaseUserUpdate(t *testing.T) {
	t.Parallel()
	pName := acctest.RandomWithPrefix("tf-db-rs")
	rName := acctest.RandomWithPrefix("tf-db-user-up")

	name := "vultr_database_user.test_user"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseUserBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "username", rName),
					resource.TestCheckResourceAttr(name, "password", "someRandomPW4928!z"),
				),
			},
			{
				Config: testAccVultrDatabaseBase(pName) + testAccVultrDatabaseUserBaseUpdatedPassword(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "username", rName),
					resource.TestCheckResourceAttr(name, "password", "someNewPW8385@x"),
				),
			},
		},
	})
}

func testAccCheckVultrDatabaseUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_database_user" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.Database.GetUser(context.Background(), rs.Primary.Attributes["database_id"], rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Not a valid database user") || strings.Contains(err.Error(), "Not a valid DBaaS Subscription UUID") {
				return nil
			}
			return fmt.Errorf("error getting database user: %s", err)
		}

		return fmt.Errorf("database user %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrDatabaseUserBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_user" "test_user" {
			database_id = vultr_database.test.id
			username = "%s"
			password = "someRandomPW4928!z"
		} `, name)
}

func testAccVultrDatabaseUserBaseUpdatedPassword(name string) string {
	return fmt.Sprintf(`
		resource "vultr_database_user" "test_user" {
			database_id = vultr_database.test.id
			username = "%s"
			password = "someNewPW8385@x"
		} `, name)
}
