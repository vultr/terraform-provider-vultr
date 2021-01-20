package vultr

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVultrIsoBase(t *testing.T) {
	t.Parallel()
	url := "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.3-x86_64.iso"
	updateURL := "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.4-x86_64.iso"
	name := "vultr_iso_private.alpine"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrIsoScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrIsoBase(url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "size"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
					resource.TestCheckResourceAttr(name, "filename", "alpine-virt-3.9.3-x86_64.iso"),
				),
			},
			{
				Config: testAccVultrIsoBase(updateURL),
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

		if _, err := client.ISO.Get(context.Background(), rs.Primary.ID); err == nil {
			return fmt.Errorf("ISO still exists : %s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccVultrIsoBase(url string) string {
	return fmt.Sprintf(`resource "vultr_iso_private" "alpine" {
		url = "%s"
		}`, url)
}
