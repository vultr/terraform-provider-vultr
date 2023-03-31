package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrReservedIPIPv4(t *testing.T) {
	rServerLabel := acctest.RandomWithPrefix("tf-vps-rip4")
	rLabel := acctest.RandomWithPrefix("tf-rip4-rs")
	rLabelUpdated := rLabel + "_updated"
	ipType := "v4"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
			{
				Config: testAccVultrReservedIPConfigAttach(rServerLabel, rLabelUpdated, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabelUpdated),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "instance_id"),
				),
			},
			{
				// test detach by unsetting the attached_id
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
		},
	})
}

func TestAccVultrReservedIPIPv6(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
			{
				Config: testAccVultrReservedIPConfigAttach(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "instance_id"),
				),
			},
			{
				// test detach by unsetting the attached_id
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
		},
	})
}

func TestAccVultrReservedIPLabelUpdate(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-rip-rs")
	rLabelUpdated := rLabel + "_updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfigLabel(rLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
				),
			},
			{
				Config: testAccVultrReservedIPConfigLabel(rLabelUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabelUpdated),
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

		_, _, err := client.ReservedIP.Get(context.Background(), ripID)
		if err == nil {
			return fmt.Errorf("reserved IP still exists: %s", ripID)
		}
	}
	return nil
}

func testAccCheckVultrReservedIPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("reserved IP ID is not set")
		}

		ripID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		if _, _, err := client.ReservedIP.Get(context.Background(), ripID); err != nil {
			return fmt.Errorf("reserved IP does not exist: %s", ripID)
		}

		return nil
	}
}

func testAccVultrReservedIPConfigLabel(label string) string {
	return fmt.Sprintf(`

   resource "vultr_reserved_ip" "foo" {
       label     = "%s"
       region    = "ewr"
       ip_type   = "v4"
   }`, label)
}

func testAccVultrReservedIPConfig(rServerLabel, label, ipType string) string {
	return fmt.Sprintf(`
	resource "vultr_instance" "ip" {
       label       = "%s"
       region      = "ewr"
       plan        = "vc2-1c-1gb"
       os_id       = 167
	   enable_ipv6 = true
   }
   resource "vultr_reserved_ip" "foo" {
       label     = "%s"
       region    = "ewr"
       ip_type   = "%s"
   }`, rServerLabel, label, ipType)
}

func testAccVultrReservedIPConfigAttach(rServerLabel, label, ipType string) string {
	return fmt.Sprintf(`
	resource "vultr_instance" "ip" {
       label = "%s"
       region = "ewr"
       plan = "vc2-1c-1gb"
       os_id = 167
       enable_ipv6 = true
   }
   resource "vultr_reserved_ip" "foo" {
       label       = "%s"
       region      = "ewr"
       ip_type     = "%s"
       instance_id = "${vultr_instance.ip.id}"
   }`, rServerLabel, label, ipType)
}
