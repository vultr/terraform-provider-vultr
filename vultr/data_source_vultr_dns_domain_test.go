package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrDnsDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDnsDomain_read("domain-test.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_dns_domain.my-site", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_dns_domain.my-site", "domain"),
					resource.TestCheckResourceAttrSet("data.vultr_dns_domain.my-site", "date_created"),
				),
			},
			{
				Config:      testAccVultrDnsDomain_noResults("bad-domain.com"),
				ExpectError: regexp.MustCompile(`.* data.vultr_dns_domain.my-site: data.vultr_dns_domain.my-site: no results were found`),
			},
			{
				Config:      testAccVultrDnsDomain_noDomain(),
				ExpectError: regexp.MustCompile(`config is invalid: data.vultr_dns_domain.my-site: "domain": required field is not set`),
			},
			{
				Config:      testAccVultrDnsDomain_emptyDomain(),
				ExpectError: regexp.MustCompile(`config is invalid: data.vultr_dns_domain.my-site: domain must not be empty`),
			},
		},
	})
}

func testAccVultrDnsDomain_read(domain string) string {
	return fmt.Sprintf(`data "vultr_dns_domain" "my-site" {
  domain = "%s"
}`, domain)
}

func testAccVultrDnsDomain_noResults(name string) string {
	return fmt.Sprintf(`data "vultr_dns_domain" "my-site" {
  domain = "%s"
}`, name)
}

func testAccVultrDnsDomain_noDomain() string {
	return `data "vultr_dns_domain" "my-site" {}`
}

func testAccVultrDnsDomain_emptyDomain() string {
	return `data "vultr_dns_domain" "my-site" {domain = ""}`
}
