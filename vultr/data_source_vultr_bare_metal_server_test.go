package vultr

//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
//)
//
//func TestAccDataSourceVultrBareMetalServer(t *testing.T) {
//	t.Parallel()
//	rName := acctest.RandomWithPrefix("tf-bms-ds")
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:  func() { testAccPreCheck(t) },
//		Providers: testAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCheckVultrBareMetalServer(rName),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "os"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "ram"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "disk"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "main_ip"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "cpu_count"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "location"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "region_id"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "default_password"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "netmask_v4"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "gateway_v4"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "plan_id"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "label"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "tag"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "os_id"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "app_id"),
//					resource.TestCheckResourceAttrSet("data.vultr_bare_metal_server.server", "v6_networks.#"),
//				),
//			},
//		},
//	})
//}
//
//func testAccCheckVultrBareMetalServer(label string) string {
//	return fmt.Sprintf(`
//		resource "vultr_bare_metal_server" "foo" {
//			region_id 		  = 5
//			os_id 			  = 270
//			plan_id           = 100
//			enable_ipv6       = true
//			notify_activate   = false
//			label             = "%s"
//			tag 			  = "bms-tag"
//		}
//
//		data "vultr_bare_metal_server" "server" {
//			filter {
//				name = "label"
//				values = ["${vultr_bare_metal_server.foo.label}"]
//			}
//		}`, label)
//}
