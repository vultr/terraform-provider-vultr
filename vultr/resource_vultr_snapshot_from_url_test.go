package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrSnapshotFromURL_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrSnapshotFromURLConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrSnapshotExists("vultr_snapshot_from_url.foo"),
					resource.TestCheckResourceAttrSet("vultr_snapshot_from_url.foo", "url"),
					resource.TestCheckResourceAttrSet("vultr_snapshot_from_url.foo", "description"),
					resource.TestCheckResourceAttrSet("vultr_snapshot_from_url.foo", "date_created"),
					resource.TestCheckResourceAttrSet("vultr_snapshot_from_url.foo", "size"),
					resource.TestCheckResourceAttrSet("vultr_snapshot_from_url.foo", "status"),
					resource.TestCheckResourceAttrSet("vultr_snapshot_from_url.foo", "os_id"),
					resource.TestCheckResourceAttrSet("vultr_snapshot_from_url.foo", "app_id"),
				),
			},
		},
	})
}

func testAccVultrSnapshotFromURLConfigBasic() string {
	return `resource "vultr_snapshot_from_url" "foo" {url = "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.1-x86_64.iso"}`
}
