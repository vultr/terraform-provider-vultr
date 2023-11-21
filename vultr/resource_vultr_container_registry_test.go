package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrContainerRegistryDBBasic(t *testing.T) {
	crName := acctest.RandomWithPrefix("tf-cr-cr")

	name := "vultr_container_registry.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrContainerRegistryDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrContainerRegistryBase(crName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", crName),
					resource.TestCheckResourceAttrSet(name, "plan"),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "public"),
				),
			},
		},
	})
}

func testAccCheckVultrContainerRegistryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_container_registry" {
			continue
		}

		crID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, _, err := client.ContainerRegistry.Get(context.Background(), crID)
		if err == nil {
			return fmt.Errorf("vpc still exists: %s", crID)
		}
	}
	return nil
}

func testAccCheckVultrContainerRegistryExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Container Registry ID is not set")
		}

		crId := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.ContainerRegistry.Get(context.Background(), crId)
		if err != nil {
			return fmt.Errorf("Container Registry does not exist: %s", crId)
		}
		return nil
	}
}

func testAccVultrContainerRegistryBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_container_registry" "foo" {
			name = "%s"
			region = "alt"
			public = true
			plan = "start_up"
		} `, name)
}
