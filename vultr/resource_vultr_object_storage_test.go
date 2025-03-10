package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrObjectStorageBasic(t *testing.T) {
	t.Parallel()

	rLabel := acctest.RandomWithPrefix("tf-s3")
	name := "vultr_object_storage.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrObjectStorageBase(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
				),
			},
		},
	})
}

func TestAccVultrObjectStorageUpdateLabel(t *testing.T) {
	t.Parallel()

	rLabel := acctest.RandomWithPrefix("tf-s3")
	updatedLabel := acctest.RandomWithPrefix("tf-s3")
	name := "vultr_object_storage.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrObjectStorageBase(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", rLabel),
				),
			},
			{
				Config: testAccVultrObjectStorageUpdated(updatedLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "label", updatedLabel),
				),
			},
		},
	})
}

func testAccVultrObjectStorageBase(label string) string {
	return fmt.Sprintf(`
		resource "vultr_object_storage" "test" {
			cluster_id = 9
			tier_id = 5
			label = "%s"
		}`, label)
}

func testAccVultrObjectStorageUpdated(label string) string {
	return fmt.Sprintf(`
		resource "vultr_object_storage" "test" {
			cluster_id = 9
			tier_id = 5
			label = "%s"
		}`, label)
}
