package vultr

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrIso_base(t *testing.T) {

	url := "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.3-x86_64.iso"
	updateUrl := "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.4-x86_64.iso"
	name := "vultr_iso_private.alpine"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrIsoScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrIso_base(url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "size"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttr(name, "filename", "alpine-virt-3.9.3-x86_64.iso"),
				),
			},
			{
				Config: testAccVultrIso_base(updateUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "size"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttr(name, "filename", "alpine-virt-3.9.4-x86_64.iso"),
				),
			},
		},
	})
}

func testAccCheckVultrIsoScriptDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).govultrClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_iso_private" {
			continue
		}

		isoList, err := client.Iso.GetList(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting list of ISO : %s", err)
		}

		exists := false
		for i := range isoList {
			if strconv.Itoa(isoList[i].IsoID) == rs.Primary.ID {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("ISO still exists : %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccVultrIso_base(url string) string {
	return fmt.Sprintf(`resource "vultr_iso_private" "alpine" {
		url = "%s"
		}`, url)
}
