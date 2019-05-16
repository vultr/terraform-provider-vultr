package vultr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"testing"
	"time"
)

func TestAccVultrDnsRecord_basic(t *testing.T) {

	rString := acctest.RandString(6) + ".com"
	rSub := acctest.RandString(4) + rString
	name := "vultr_dns_record.a-record"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckVultrDnsDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDnsDomain_base(rString) + testAccVultrDnsRecord_base(rSub),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrDomainRecordExists,
					resource.TestCheckResourceAttr(name, "name", rSub),
					resource.TestCheckResourceAttr(name, "domain", rString),
					resource.TestCheckResourceAttr(name, "data", "10.0.0.1"),
					resource.TestCheckResourceAttr(name, "type", "A"),
					resource.TestCheckResourceAttr(name, "ttl", "3600"),
				),
			},
		},
	})
}

func testAccCheckVultrDomainRecordExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).govultrClient()


	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_dns_record" {
			continue
		}

		id := rs.Primary.ID
		domain := rs.Primary.Attributes["domain"]
		records, err  := client.DNSRecord.GetList(context.Background(), domain)

		if err != nil {
			return fmt.Errorf("Error getting dns record %s for domain %s : %v", id, domain, err)
		}

		exists := false
		for _, v  := range records {
			if strconv.Itoa(v.RecordID) == id {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Error getting dns record %s for domain %s : %v", id, domain, err)
		}
	}

	return nil
}


func testAccVultrDnsRecord_base(name string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`resource "vultr_dns_record" "a-record" {
  data = "10.0.0.1"
  domain = "${vultr_dns_domain.my-site.id}"
  name = "%s"
  type = "A"
  ttl = "3600"
}`, name)
}

