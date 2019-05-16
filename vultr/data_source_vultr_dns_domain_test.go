package vultr

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
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
			{
				Config:      testAccVultrDnsDomain_noResults(domain),
				ExpectError: regexp.MustCompile(`.* data.vultr_dns_domain.my-site: data.vultr_dns_domain.my-site: no results were found`),
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
  				domain = "%s",
  				server_ip = "10.0.0.0"
			}`, domain)
}

func testAccVultrDnsDomain_noResults(name string) string {
	return fmt.Sprintf(`
		data "vultr_dns_domain" "my-site" {
 			domain = "%s"
		}`, name)
}
