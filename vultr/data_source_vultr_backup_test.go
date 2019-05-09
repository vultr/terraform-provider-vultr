package vultr

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
)

func TestAccVultrBackup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBackup_read("auto-backup 63.209.32.248"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.vultr_backup.backs", "size", "860174523"),
					resource.TestCheckResourceAttr(
						"data.vultr_backup.backs", "status", "complete"),
					resource.TestCheckResourceAttr(
						"data.vultr_backup.backs", "description", "auto-backup 63.209.32.248 "),
					resource.TestCheckResourceAttr(
						"data.vultr_backup.backs", "date_created", "2019-05-09 00:09:57"),
					resource.TestCheckResourceAttr(
						"data.vultr_backup.backs", "id", "bca5cd36fd57c"),
				),
			},
			{
				Config: testAccVultrBackup_noResults("auto-backup 63.209.32.244"),
				ExpectError: regexp.MustCompile(`.* data.vultr_backup.backs: data.vultr_backup.backs: no results were found`),
			},
		},
	})
}

func testAccVultrBackup_read(description string) string {
	return fmt.Sprintf(`data "vultr_backup" "backs" {
  filter {
    name = "description"
    values = ["%s "]
  }
}`, description)
}

func testAccVultrBackup_noResults(description string) string {
	return fmt.Sprintf(`data "vultr_backup" "backs" {
  filter {
    name = "description"
    values = ["%s "]
  }
}`, description)
}
