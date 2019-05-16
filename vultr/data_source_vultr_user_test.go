package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrUser_dataBase(t *testing.T) {

	rEmail := fmt.Sprintf("terraform-%s@vultr.com", acctest.RandString(4))
	name := "data.vultr_user.admin"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrUserConfig_base(rEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "email", rEmail),
					resource.TestCheckResourceAttr(name, "name", "Terraform AccTests"),
					resource.TestCheckResourceAttr(name, "acl.#", "10"),
					resource.TestCheckResourceAttr(name, "acl.0", "manage_users"),
					resource.TestCheckResourceAttr(name, "acl.1", "subscriptions"),
					resource.TestCheckResourceAttr(name, "acl.2", "billing"),
					resource.TestCheckResourceAttr(name, "acl.3", "support"),
					resource.TestCheckResourceAttr(name, "acl.4", "provisioning"),
					resource.TestCheckResourceAttr(name, "acl.5", "dns"),
					resource.TestCheckResourceAttr(name, "acl.6", "abuse"),
					resource.TestCheckResourceAttr(name, "acl.7", "upgrade"),
					resource.TestCheckResourceAttr(name, "acl.8", "firewall"),
					resource.TestCheckResourceAttr(name, "acl.9", "alerts"),
					resource.TestCheckResourceAttr(name, "api_enabled", "yes"),
					resource.TestCheckResourceAttrSet(name, "id"),
				),
			},
			{
				Config:      testAccVultrUserConfig_noResult(rEmail),
				ExpectError: regexp.MustCompile(fmt.Sprintf(".*%s: %s: no results were found", name, name)),
			},
		},
	})
}

func testAccVultrUserConfig_base(email string) string {
	return fmt.Sprintf(`
		data "vultr_user" "admin" {
			filter {
    			name = "email"
    			values = ["${vultr_user.admin.email}"]
  				}
			}

		resource "vultr_user" "admin" {
  			name = "Terraform AccTests",
  			email = "%s"
  			password = "password",
  			acl = [
            	"manage_users",
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

		`, email)
}

func testAccVultrUserConfig_noResult(email string) string {
	return fmt.Sprintf(`
		data "vultr_user" "admin" {
			filter {
    			name = "email"
    			values = ["%s"]
  				}
			}
		`, email)
}
