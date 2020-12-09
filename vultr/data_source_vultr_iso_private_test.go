package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrIsoPrivate(t *testing.T) {
	t.Parallel()
	url := "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.2-x86_64.iso"
	name := "data.vultr_iso_private.alpine"
	fileName := "alpine-virt-3.9.2-x86_64.iso"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrIsoScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrIsoPrivate_read(url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "size"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttrSet(name, "md5sum"),
					resource.TestCheckResourceAttrSet(name, "sha512sum"),
					resource.TestCheckResourceAttr(name, "filename", fileName),
				),
			},
		},
	})
}

func testAccVultrIsoPrivate_read(description string) string {
	return fmt.Sprintf(`
		resource "vultr_iso_private" "alpine" {
			url = "%s"
		}

		data "vultr_iso_private" "alpine" {
			filter {
				name = "filename"
				values = ["${vultr_iso_private.alpine.filename}"]
			}
		}`, description)
}
