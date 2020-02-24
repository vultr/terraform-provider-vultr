package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVultrReservedIP_IPv4(t *testing.T) {
	rServerLabel := acctest.RandomWithPrefix("tf-vps-rip4")
	rLabel := acctest.RandomWithPrefix("tf-rip4-rs")
	ipType := "v4"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
			{
				Config: testAccVultrReservedIPConfig_attach(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "attached_id"),
				),
			},
			{
				// test detach by unsetting the attached_id
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
		},
	})
}

func TestAccVultrReservedIP_IPv6(t *testing.T) {
	rServerLabel := acctest.RandomWithPrefix("tf-vps-rip6")
	rLabel := acctest.RandomWithPrefix("tf-rip6-rs")
	ipType := "v6"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
			{
				Config: testAccVultrReservedIPConfig_attach(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "attached_id"),
				),
			},
			{
				// test detach by unsetting the attached_id
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region_id"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
		},
	})
}

func testAccCheckVultrReservedIPDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_reserved_ip" {
			continue
		}

		ripID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		rips, err := client.ReservedIP.List(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting reserved IPs: %s", err)
		}

		exists := false
		for i := range rips {
			if rips[i].ReservedIPID == ripID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("Reserved IP still exists: %s", ripID)
		}
	}
	return nil
}

func testAccCheckVultrReservedIPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Reserved IP ID is not set")
		}

		ripID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		rips, err := client.ReservedIP.List(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting reserved IPs: %s", err)
		}

		exists := false
		for i := range rips {
			if rips[i].ReservedIPID == ripID {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Reserved IP does not exist: %s", ripID)
		}

		return nil
	}
}

func testAccVultrReservedIPConfig(rServerLabel, label, ipType string) string {
	return fmt.Sprintf(`
	resource "vultr_server" "ip" {
        label = "%s"
        region_id = 6
        plan_id = 201
        os_id = 147
		enable_ipv6 = true
    }
    resource "vultr_reserved_ip" "foo" {
        label       = "%s"
        region_id   = 6
        ip_type        = "%s"
    }
   `, rServerLabel, label, ipType)
}

func testAccVultrReservedIPConfig_attach(rServerLabel, label, ipType string) string {
	return fmt.Sprintf(`
	resource "vultr_server" "ip" {
        label = "%s"
        region_id = 6
        plan_id = 201
        os_id = 147
		enable_ipv6 = true
    }
    resource "vultr_reserved_ip" "foo" {
        label       = "%s"
        region_id   = 6
        ip_type        = "%s"
        attached_id = "${vultr_server.ip.id}"
    }
   `, rServerLabel, label, ipType)
}
