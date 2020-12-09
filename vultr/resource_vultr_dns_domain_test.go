package vultr

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVultrDnsDomain_basic(t *testing.T) {

	rString := acctest.RandString(6) + ".com"
	name := "vultr_dns_domain.my-site"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrDnsDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDnsDomain_base(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", rString),
					resource.TestCheckResourceAttr(name, "domain", rString),
					resource.TestCheckResourceAttr(name, "ip", "10.0.0.0"),
				),
			},
			{
				Config: testAccVultrDnsDomain_update(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", rString),
					resource.TestCheckResourceAttr(name, "domain", rString),
					resource.TestCheckResourceAttr(name, "ip", "10.0.0.1"),
				),
			},
		},
	})
}

func TestAccVultrDnsDomain_newDomainForce(t *testing.T) {
	rString := acctest.RandString(6) + ".com"
	newDomain := acctest.RandString(6) + ".com"
	name := "vultr_dns_domain.my-site"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrDnsDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDnsDomain_base(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", rString),
					resource.TestCheckResourceAttr(name, "domain", rString),
					resource.TestCheckResourceAttr(name, "ip", "10.0.0.0"),
				),
			},
			{
				Config: testAccVultrDnsDomain_update(newDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", newDomain),
					resource.TestCheckResourceAttr(name, "domain", newDomain),
					resource.TestCheckResourceAttr(name, "ip", "10.0.0.1"),
				),
			},
		},
	})
}

func testAccCheckVultrDnsDomainDestroy(s *terraform.State) error {
	time.Sleep(1 * time.Second)
	client := testAccProvider.Meta().(*Client).govultrClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_dns_domain" {
			continue
		}

		_, err := client.Domain.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("domain still exists : %s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccVultrDnsDomain_base(domain string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`
		resource "vultr_dns_domain" "my-site" {
			domain = "%s"
			ip = "10.0.0.0"
		}`, domain)
}

func testAccVultrDnsDomain_update(domain string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`
		resource "vultr_dns_domain" "my-site" {
			domain = "%s"
			ip = "10.0.0.1"
		}`, domain)
}
