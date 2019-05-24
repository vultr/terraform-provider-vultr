package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVultrBackup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBackup_read("auto-backup 63.209.32.248 server-label"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_backup.backs", "size"),
					resource.TestCheckResourceAttr("data.vultr_backup.backs", "status", "complete"),
					resource.TestCheckResourceAttr("data.vultr_backup.backs", "description", "auto-backup 63.209.32.248 server-label"),
					resource.TestCheckResourceAttrSet("data.vultr_backup.backs", "date_created"),
					resource.TestCheckResourceAttrSet("data.vultr_backup.backs", "id"),
				),
			},
		},
	})
}

func testAccVultrBackup_read(description string) string {
	return fmt.Sprintf(`
		data "vultr_backup" "backs" {
  			filter {
    			name = "description"
    			values = ["%s"]
  			}
		}`, description)
}
