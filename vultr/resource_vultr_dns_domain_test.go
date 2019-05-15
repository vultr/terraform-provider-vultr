package vultr

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					resource.TestCheckResourceAttr(name, "server_ip", "10.0.0.0"),
				),
			},
			{
				Config: testAccVultrDnsDomain_update(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", rString),
					resource.TestCheckResourceAttr(name, "domain", rString),
					resource.TestCheckResourceAttr(name, "server_ip", "10.0.0.1"),
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
					resource.TestCheckResourceAttr(name, "server_ip", "10.0.0.0"),
				),
			},
			{
				Config: testAccVultrDnsDomain_update(newDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", newDomain),
					resource.TestCheckResourceAttr(name, "domain", newDomain),
					resource.TestCheckResourceAttr(name, "server_ip", "10.0.0.1"),
				),
			},
		},
	})
}

func testAccCheckVultrDnsDomainDestroy(s *terraform.State) error {
	time.Sleep(1 * time.Second)
	client := testAccProvider.Meta().(*Client).govultrClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_user" {
			continue
		}

		domains, err := client.DNSDomain.GetList(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting list of domains : %s", err)
		}

		exists := false
		for i := range domains {
			if domains[i].Domain == rs.Primary.ID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("User still exists : %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccVultrDnsDomain_base(domain string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`resource "vultr_dns_domain" "my-site" {
  domain = "%s",
  server_ip = "10.0.0.0"
}`, domain)
}

func testAccVultrDnsDomain_update(domain string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`resource "vultr_dns_domain" "my-site" {
  domain = "%s",
  server_ip = "10.0.0.1"
}`, domain)
}
