package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrUser("terraform-acceptance@vultr.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "email", "terraform-acceptance@vultr.com"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "name", "Terraform AccTests"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "id", "ef35cd31c6de8"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.#", "10"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.0", "manage_users"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.1", "subscriptions"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.2", "billing"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.3", "support"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.4", "provisioning"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.5", "dns"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.6", "abuse"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.7", "upgrade"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.8", "firewall"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "acl.9", "alerts"),
					resource.TestCheckResourceAttr(
						"data.vultr_user.admin", "api_enabled", "yes"),
				),
			},
		},
	})
}

func testAccVultrUser(email string) string {
	return fmt.Sprintf(`data "vultr_user" "admin" {
  			filter {
    			name = "email"
    			values = ["%s"]
  			}
			}`, email)
}
