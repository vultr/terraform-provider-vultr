package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVultrUser_base(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVultrUser_create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "email", "terraform-acceptance@vultr.com"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "name", "Terraform AccTests"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.#", "12"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.0", "manage_users"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.1", "subscriptions_view"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.2", "subscriptions"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.3", "billing"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.4", "support"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.5", "provisioning"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.6", "dns"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.7", "abuse"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.8", "upgrade"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.9", "firewall"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.10", "alerts"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.11", "objstore"),
					//resource.TestCheckResourceAttr(
					//	"vultr_user.admin", "api_enabled", "true"),
					resource.TestCheckResourceAttrSet("vultr_user.admin", "api_key"),
				),
			},
			{
				Config: testAccResourceVultrUser_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "email", "terraform-acceptance@vultr.com"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "name", "Terraform Update Name"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.#", "11"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.0", "manage_users"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.1", "subscriptions_view"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.2", "subscriptions"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.3", "billing"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.4", "support"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.5", "provisioning"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.6", "dns"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.7", "abuse"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.8", "upgrade"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.9", "firewall"),
					resource.TestCheckResourceAttr(
						"vultr_user.admin", "acl.10", "alerts"),
					//resource.TestCheckResourceAttr(
					//	"vultr_user.admin", "api_enabled", "false"),
					resource.TestCheckResourceAttrSet("vultr_user.admin", "api_key"),
				),
			},
		},
	})
}

func testAccCheckVultrUsersDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*Client).govultrClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_user" {
			continue
		}

		if _, err := client.User.Get(context.Background(), rs.Primary.ID); err == nil {
			return fmt.Errorf("user still exists : %s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccResourceVultrUser_create() string {
	return `resource "vultr_user" "admin" {
name = "Terraform AccTests"
email = "terraform-acceptance@vultr.com"
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
          "alerts",
			"objstore"
]
#api_enabled = true
}`
}

func testAccResourceVultrUser_update() string {
	return `resource "vultr_user" "admin" {
name = "Terraform Update Name"
email = "terraform-acceptance@vultr.com"
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
#api_enabled = false
}`
}
