package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrBackup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBackupRead("auto-backup 45.77.152.237 tf-backup"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vultr_backup.backs", "backups.0.size"),
					resource.TestCheckResourceAttrSet("data.vultr_backup.backs", "backups.0.date_created"),
					resource.TestCheckResourceAttr("data.vultr_backup.backs", "backups.0.status", "complete"),
				),
			},
		},
	})
}

func testAccVultrBackupRead(description string) string {
	return fmt.Sprintf(`
		data "vultr_backup" "backs" {
			filter {
				name = "description"
				values = ["%s"]
			}
		}`, description)
}
