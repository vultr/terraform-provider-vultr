package vultr

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrDNSDomainBasic(t *testing.T) {

	rString := acctest.RandString(6) + ".com"
	name := "vultr_dns_domain.my-site"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDNSDomainBase(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", rString),
					resource.TestCheckResourceAttr(name, "domain", rString),
					resource.TestCheckResourceAttr(name, "ip", "10.0.0.0"),
				),
			},
			{
				Config: testAccVultrDNSDomainUpdate(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", rString),
					resource.TestCheckResourceAttr(name, "domain", rString),
					resource.TestCheckResourceAttr(name, "ip", "10.0.0.1"),
				),
			},
		},
	})
}

func TestAccVultrDNSDomainNewDomainForce(t *testing.T) {
	rString := acctest.RandString(6) + ".com"
	newDomain := acctest.RandString(6) + ".com"
	name := "vultr_dns_domain.my-site"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDNSDomainBase(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", rString),
					resource.TestCheckResourceAttr(name, "domain", rString),
					resource.TestCheckResourceAttr(name, "ip", "10.0.0.0"),
				),
			},
			{
				Config: testAccVultrDNSDomainUpdate(newDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", newDomain),
					resource.TestCheckResourceAttr(name, "domain", newDomain),
					resource.TestCheckResourceAttr(name, "ip", "10.0.0.1"),
				),
			},
		},
	})
}

func testAccCheckVultrDNSDomainDestroy(s *terraform.State) error {
	time.Sleep(1 * time.Second)
	client := testAccProvider.Meta().(*Client).govultrClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_dns_domain" {
			continue
		}

		_, _, err := client.Domain.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("domain still exists : %s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccVultrDNSDomainBase(domain string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`
		resource "vultr_dns_domain" "my-site" {
			domain = "%s"
			ip = "10.0.0.0"
		}`, domain)
}

func testAccVultrDNSDomainUpdate(domain string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`
		resource "vultr_dns_domain" "my-site" {
			domain = "%s"
			ip = "10.0.0.1"
		}`, domain)
}
