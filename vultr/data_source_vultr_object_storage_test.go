package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrObjectStorage(t *testing.T) {
	t.Parallel()
	rLabel := acctest.RandomWithPrefix("tf-test-s3")
	name := "data.vultr_object_storage.s3"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrObjectStorage(rLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "label"),
					resource.TestCheckResourceAttrSet(name, "cluster_id"),
					resource.TestCheckResourceAttrSet(name, "region"),
					resource.TestCheckResourceAttrSet(name, "location"),
					resource.TestCheckResourceAttrSet(name, "s3_access_key"),
					resource.TestCheckResourceAttrSet(name, "s3_hostname"),
					resource.TestCheckResourceAttrSet(name, "s3_secret_key"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
				),
			},
		},
	})
}

func testAccCheckVultrObjectStorage(label string) string {
	return fmt.Sprintf(`
		resource "vultr_object_storage" "tf" {
			cluster_id = 2
			label = "%s"
		}

		data "vultr_object_storage" "s3" {
			filter {
				name = "label"
				values = ["${vultr_object_storage.tf.label}"]
			}
		}`, label)
}
