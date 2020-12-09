package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrDnsDomain_dataBase(t *testing.T) {
	domain := fmt.Sprintf("%s.com", acctest.RandString(6))
	name := "data.vultr_dns_domain.my-site"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDnsDomainConfig(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttr(name, "domain", domain),
					resource.TestCheckResourceAttrSet(name, "date_created"),
				),
			},
		},
	})
}

func testAccVultrDnsDomainConfig(domain string) string {
	return fmt.Sprintf(`
			data "vultr_dns_domain" "my-site" {
				domain = "${vultr_dns_domain.my-site.id}"
			}

			resource "vultr_dns_domain" "my-site" {
				domain = "%s"
				ip = "10.0.0.0"
			}`, domain)
}
