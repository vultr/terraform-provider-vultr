package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVultrBareMetalServer(t *testing.T) {
	rInt := acctest.RandInt()
	rName := acctest.RandomWithPrefix("tf-test")
	rSSH, _, err := acctest.RandSSHKeyPair("foobar")
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBareMetalServerConfigBasic(rInt, rSSH, rName),
			},
			{
				Config: testAccVultrBareMetalServerConfigBasic(rInt, rSSH, rName) + testAccCheckVultrBareMetalServer(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "os"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "ram"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "disk"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "main_ip"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "cpu_count"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "location"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "region_id"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "default_password"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "netmask_v4"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "gateway_v4"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "plan_id"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "label"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "tag"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "os_id"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "app_id"),
					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "v6_networks.#"),
				),
			},
			{
				Config:      testAccCheckVultrBareMetalServer(rName),
				ExpectError: regexp.MustCompile(`.* data.vultr_bare_metal_server.server: data.vultr_bare_metal_server.server: no results were found`),
			},
		},
	})
}

func testAccCheckVultrBareMetalServer(label string) string {
	return fmt.Sprintf(`data "vultr_bare_metal_server" "server" {
		filter {
		name = "label"
		values = ["%s"]
		}
		}`, label)
}
