package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrReservedIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIp_read("reserved-ip"),
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
				Config:      testAccVultrReservedIp_noResult("foobar"),
				ExpectError: regexp.MustCompile(`.* data.vultr_reserved_ip.rs: data.vultr_reserved_ip.rs: no results were found`),
			},
		},
	})
}

func testAccVultrReservedIp_read(label string) string {
	return fmt.Sprintf(`
		data "vultr_reserved_ip" "rs" {
		filter {
    	name = "label"
    	values = ["%s"]
  	}
	}`, label)
}

func testAccVultrReservedIp_noResult(label string) string {
	return fmt.Sprintf(`
		data "vultr_reserved_ip" "rs" {
		filter {
    	name = "label"
    	values = ["%s"]
  	}
	}`, label)
}
