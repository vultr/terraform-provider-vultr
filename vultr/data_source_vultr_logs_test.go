package vultr

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVultrLogs(t *testing.T) {
	t.Parallel()
	name := "data.vultr_logs.test_logs"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVultrLogs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "results.#"),
					resource.TestCheckResourceAttr(name, "log_level", "debug"),
				),
			},
		},
	})
}

func testAccCheckVultrLogs() string {
	now := time.Now()
	month := now.Month()
	day := now.Day()
	config := fmt.Sprintf(`
data "vultr_logs" "test_logs" {
	start_time = "2026-%02d-%02dT00:00:00Z"
	end_time = "2026-%02d-%02dT00:00:00Z"
	log_level = "debug"
}`, month, day, month, day+1)

	fmt.Println(config)

	return config
}
