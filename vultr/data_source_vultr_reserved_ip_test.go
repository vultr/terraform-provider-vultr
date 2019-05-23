package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVultrReservedIP(t *testing.T) {
	rServerLabel := acctest.RandomWithPrefix("tf-vps-rip-ds")
	rLabel := acctest.RandomWithPrefix("tf-test-")
	ipType := "v4"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
			},
			{
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType) + testAccVultrReservedIP_read(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.rs", "id"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.rs", "subnet_size"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.rs", "subnet"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.rs", "region_id"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.rs", "label"),
					resource.TestCheckResourceAttrSet("data.vultr_reserved_ip.rs", "ip_type"),
				),
			},
			{
				Config:      testAccVultrReservedIP_noResult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_reserved_ip.rs: data.vultr_reserved_ip.rs: no results were found`),
			},
		},
	})
}

func testAccVultrReservedIP_read(label string) string {
	return fmt.Sprintf(`
		data "vultr_reserved_ip" "rs" {
		filter {
    	name = "label"
    	values = ["%s"]
  	}
	}`, label)
}

func testAccVultrReservedIP_noResult(label string) string {
	return fmt.Sprintf(`
		data "vultr_reserved_ip" "rs" {
		filter {
    	name = "label"
    	values = ["%s"]
  	}
	}`, label)
}
