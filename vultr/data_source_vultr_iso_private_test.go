package vultr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrIsoPrivate(t *testing.T) {

	url := "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.4-x86_64.iso"
	name := "data.vultr_iso_private.alpine"
	fileName := "alpine-virt-3.9.4-x86_64.iso"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrIsoScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrIso_base(url),
			},
			{
				Config: testAccVultrIso_base(url) + testAccVultrIsoPrivate_read(fileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "size"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttr(name, "filename", fileName),
				),
			},
			{
				Config:      testAccVultrIsoPrivate_noResults(fileName),
				ExpectError: regexp.MustCompile(`.* data.vultr_iso_private.alpine: data.vultr_iso_private.alpine: no results were found`),
			},
		},
	})
}

func testAccVultrIsoPrivate_read(description string) string {
	return fmt.Sprintf(`
		data "vultr_iso_private" "alpine" {
  			filter {
    			name = "filename"
    			values = ["%s"]
  			}
		}`, description)
}

func testAccVultrIsoPrivate_noResults(name string) string {
	return fmt.Sprintf(`
		data "vultr_iso_private" "alpine" {
  			filter {
    			name = "filename"
    			values = ["%s"]
  			}
		}`, name)
}
