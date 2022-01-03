package vultr

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrBackup(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping testing in Non-CI environment")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBackupRead("auto-backup 100.68.58.70 TF-BACKUPS-DND"),
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
