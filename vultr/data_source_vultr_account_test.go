package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrAccount(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_account.account", "name"),
					resource.TestCheckResourceAttrSet("data.vultr_account.account", "email"),
					resource.TestCheckResourceAttrSet("data.vultr_account.account", "balance"),
					resource.TestCheckResourceAttrSet("data.vultr_account.account", "pending_charges"),
					resource.TestCheckResourceAttrSet("data.vultr_account.account", "last_payment_date"),
					resource.TestCheckResourceAttrSet("data.vultr_account.account", "last_payment_amount"),
					resource.TestCheckResourceAttrSet("data.vultr_account.account", "acl.#"),
				),
			},
		},
	})
}

func testAccVultrAccount() string {
	return `data "vultr_account" "account" {}`
}
