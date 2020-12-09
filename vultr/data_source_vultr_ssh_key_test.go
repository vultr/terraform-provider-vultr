package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrSSHKey(t *testing.T) {

	rName := fmt.Sprintf("%s-%d-terraform", acctest.RandString(3), acctest.RandInt())
	rSSH, _, err := acctest.RandSSHKeyPair("foobar")
	name := "data.vultr_ssh_key.my_key"
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrSSHKeyConfig_basic(rName, rSSH),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttrSet(name, "ssh_key"),
					resource.TestCheckResourceAttrSet(name, "date_created"),
				),
			},
		},
	})
}

func testAccCheckVultrSSHKeyConfig_basic(name, ssh string) string {
	return fmt.Sprintf(`
		resource "vultr_ssh_key" "foo" {
			name = "%s"
			ssh_key = "%s"
		}

		data "vultr_ssh_key" "my_key" {
		filter {
			name = "name"
			values = ["${vultr_ssh_key.foo.name}"]
			}
		}
		`, name, ssh)
}
