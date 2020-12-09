package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrUser_dataBase(t *testing.T) {

	rEmail := fmt.Sprintf("terraform-%s@vultr.com", acctest.RandString(4))
	name := "data.vultr_user.admin"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrUserConfig_base(rEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "email", rEmail),
					resource.TestCheckResourceAttr(name, "name", "Terraform AccTests"),
					resource.TestCheckResourceAttr(name, "acl.#", "11"),
					resource.TestCheckResourceAttr(name, "acl.0", "manage_users"),
					resource.TestCheckResourceAttr(name, "acl.1", "subscriptions_view"),
					resource.TestCheckResourceAttr(name, "acl.2", "subscriptions"),
					resource.TestCheckResourceAttr(name, "acl.3", "billing"),
					resource.TestCheckResourceAttr(name, "acl.4", "support"),
					resource.TestCheckResourceAttr(name, "acl.5", "provisioning"),
					resource.TestCheckResourceAttr(name, "acl.6", "dns"),
					resource.TestCheckResourceAttr(name, "acl.7", "abuse"),
					resource.TestCheckResourceAttr(name, "acl.8", "upgrade"),
					resource.TestCheckResourceAttr(name, "acl.9", "firewall"),
					resource.TestCheckResourceAttr(name, "acl.10", "alerts"),
					resource.TestCheckResourceAttr(name, "api_enabled", "true"),
					resource.TestCheckResourceAttrSet(name, "id"),
				),
			},
		},
	})
}

func testAccVultrUserConfig_base(email string) string {
	return fmt.Sprintf(`
		resource "vultr_user" "admin" {
			name = "Terraform AccTests"
			email = "%s"
			password = "password"
			acl = [
				"manage_users",
				"subscriptions_view",
				"subscriptions",
				"billing",
				"support",
				"provisioning",
				"dns",
				"abuse",
				"upgrade",
				"firewall",
				"alerts"
			]
			api_enabled = true
		}

		data "vultr_user" "admin" {
			filter {
			name = "email"
			values = ["${vultr_user.admin.email}"]
				}
			}

	`, email)
}
