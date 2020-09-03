package vultr

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/terraform"
//)
//
//func TestAccVultrServerIPV4_basic(t *testing.T) {
//	t.Parallel()
//
//	name := "vultr_server_ipv4.test"
//	serverLabel := acctest.RandomWithPrefix("tf-rs-vps-server-ipv4")
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:  func() { testAccPreCheck(t) },
//		Providers: testAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccVultrServerIPV4(serverLabel),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckvultrServerIPV4Exists(name),
//					resource.TestCheckResourceAttrSet(name, "instance_id"),
//					resource.TestCheckResourceAttrSet(name, "ip"),
//					resource.TestCheckResourceAttrSet(name, "reverse"),
//				),
//			},
//		},
//	})
//}
//
//func testAccCheckvultrServerIPV4Exists(n string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("ipv4 not found: %s", n)
//		}
//
//		if rs.Primary.ID == "" {
//			return errors.New("ipv4 ID is not set")
//		}
//
//		exists, err := vultrServerIPV4Exists(rs)
//		if err != nil {
//			return err
//		}
//
//		if !exists {
//			return fmt.Errorf("ipv4 does not exist: %s", rs.Primary.ID)
//		}
//
//		return nil
//	}
//}
//
//func testAccVultrServerIPV4(serverLabel string) string {
//	return fmt.Sprintf(`
//		resource "vultr_server" "foo" {
//			plan_id = "201"
//			region_id = "6"
//			os_id = "167"
//			label = "%s"
//		}
//
//		resource "vultr_server_ipv4" "test" {
//			instance_id = "${vultr_server.foo.id}"
//			reboot = false
//		}
//	`, serverLabel)
//}
//
//func vultrServerIPV4Exists(rs *terraform.ResourceState) (bool, error) {
//	client := testAccProvider.Meta().(*Client).govultrClient()
//
//	instanceID, ok := rs.Primary.Attributes["instance_id"]
//	if !ok {
//		return false, errors.New("error getting instance ID")
//	}
//
//	ipv4s, err := client.Server.IPV4Info(context.Background(), instanceID, true)
//	if err != nil {
//		return false, fmt.Errorf("error getting IPv4s: %v", err)
//	}
//
//	ip := rs.Primary.ID
//
//	for i := range ipv4s {
//		if ipv4s[i].IP == ip {
//			return true, nil
//		}
//	}
//
//	return false, nil
//}
