package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVultrBackup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBackup_read("auto-backup 45.77.152.237 tf-backup"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_backup.backs", "backups.0.size"),
					resource.TestCheckResourceAttrSet("data.vultr_backup.backs", "backups.0.date_created"),
					resource.TestCheckResourceAttr("data.vultr_backup.backs", "backups.0.status", "complete"),
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
