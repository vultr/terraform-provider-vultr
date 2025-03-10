package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrVPC2(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-vpc2-rs-nocdir")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrVPC2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPC2ConfigBase(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrVPC2Exists("vultr_vpc2.foo"),
					resource.TestCheckResourceAttr("vultr_vpc2.foo", "description", rDesc),
					resource.TestCheckResourceAttrSet("vultr_vpc2.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_vpc2.foo", "ip_block"),
				),
			},
		},
	})
}

func TestAccVultrVPC2WithSubnet(t *testing.T) {
	rDesc := acctest.RandomWithPrefix("tf-vpc2-rs-subnet")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrVPC2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrVPC2ConfigWithSubnet(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrVPC2Exists("vultr_vpc2.foo"),
					resource.TestCheckResourceAttr("vultr_vpc2.foo", "description", rDesc),
					resource.TestCheckResourceAttr("vultr_vpc2.foo", "ip_block", "10.0.0.0"),
					resource.TestCheckResourceAttrSet("vultr_vpc2.foo", "date_created"),
				),
			},
		},
	})
}

func testAccCheckVultrVPC2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_vpc2" {
			continue
		}

		vpc2ID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, _, err := client.VPC2.Get(context.Background(), vpc2ID) //nolint:staticcheck
		if err == nil {
			return fmt.Errorf("vpc 2.0 still exists: %s", vpc2ID)
		}
	}
	return nil
}

func testAccCheckVultrVPC2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("VPC 2.0 ID is not set")
		}

		vpc2ID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, _, err := client.VPC2.Get(context.Background(), vpc2ID) //nolint:staticcheck
		if err != nil {
			return fmt.Errorf("VPC 2.0 does not exist: %s", vpc2ID)
		}

		return nil
	}
}

func testAccVultrVPC2ConfigBase(rDesc string) string {
	return fmt.Sprintf(`
		resource "vultr_vpc2" "foo" {
			region      = "atl"
			description = "%s"
		}
	`, rDesc)
}

func testAccVultrVPC2ConfigWithSubnet(rDesc string) string {
	return fmt.Sprintf(`
		resource "vultr_vpc2" "foo" {
			region        = "atl"
			description   = "%s"
			ip_type 	  = "v4"
			ip_block      = "10.0.0.0"
			prefix_length = 24
		}
	`, rDesc)
}
