package vultr

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrInferenceBasic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-inf-rs")

	name := "vultr_inference.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVultrInferenceDestroy,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrInferenceBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
				),
			},
		},
	})
}

func TestAccVultrInferenceUpdate(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-inf-rs-up")
	newName := rName + "-updated"

	name := "vultr_inference.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrInferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrInferenceBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rName),
				),
			},
			{
				Config: testAccVultrInferenceBaseUpdatedLabel(newName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", newName),
				),
			},
		},
	})
}

func testAccCheckVultrInferenceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_inference" {
			continue
		}

		client := testAccProvider.Meta().(*Client).govultrClient()
		_, _, err := client.Inference.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Not a valid Vultr Serverless Inference subscription") {
				return nil
			}
			return fmt.Errorf("error getting inference: %s", err)
		}

		return fmt.Errorf("inference %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccVultrInferenceBase(name string) string {
	return fmt.Sprintf(`
		resource "vultr_inference" "test" {
			label = "%s"
		} `, name)
}

func testAccVultrInferenceBaseUpdatedLabel(name string) string {
	return fmt.Sprintf(`
		resource "vultr_inference" "test" {
			label = "%s"
		} `, name)
}
