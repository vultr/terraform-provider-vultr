package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceVultrSnapshot(t *testing.T) {
	t.Parallel()
	rDesc := acctest.RandomWithPrefix("tf-snap-ds")
	rLabel := acctest.RandomWithPrefix("tf-test-vps")
	name := "data.vultr_snapshot.my_snapshot"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrSnapshotDestroy,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceVultrSnapshotBase(rLabel, rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", rDesc),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "os_id"),
					resource.TestCheckResourceAttrSet(name, "app_id"),
				),
			},
		},
	})
}

func testAccDataSourceVultrSnapshotBase(vpsLabel, desc string) string {
	return fmt.Sprintf(`
	resource "vultr_instance" "test" {
			plan = "vc2-1c-1gb"
			region = "ewr"
			os_id = "167"
			label = "%s"
			hostname = "testing-the-hostname"
			tag = "even better tag"
		}

		resource "vultr_snapshot" "foo" {
			instance_id  = "${vultr_instance.test.id}"
			description  = "%s"
		}

		data "vultr_snapshot" "my_snapshot" {
		filter {
			name = "description"
			values = ["${vultr_snapshot.foo.description}"]
			}
		}
		`, vpsLabel, desc)
}
