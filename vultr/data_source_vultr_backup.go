package vultr

//
//import (
//	"context"
//	"errors"
//	"fmt"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
//)
//
//func dataSourceVultrBackup() *schema.Resource {
//	return &schema.Resource{
//		Read: dataSourceVultrBackupRead,
//		Schema: map[string]*schema.Schema{
//			"filter": dataSourceFiltersSchema(),
//			"backups": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Elem:     &schema.Schema{Type: schema.TypeMap},
//			},
//		},
//	}
//}
//
//func dataSourceVultrBackupRead(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	filters, filtersOk := d.GetOk("filter")
//
//	if !filtersOk {
//		return fmt.Errorf("issue with filter: %v", filtersOk)
//	}
//
//	backups, err := client.Backup.List(context.Background())
//	if err != nil {
//		return fmt.Errorf("Error getting applications: %v", err)
//	}
//
//	var backupList []map[string]interface{}
//
//	f := buildVultrDataSourceFilter(filters.(*schema.Set))
//	for _, b := range backups {
//		// we need convert the a struct INTO a map so we can easily manipulate the data here
//		sm, err := structToMap(b)
//		if err != nil {
//			return err
//		}
//
//		if filterLoop(f, sm) {
//			backupList = append(backupList, sm)
//		}
//	}
//
//	if len(backupList) < 1 {
//		return errors.New("no results were found")
//	}
//
//	//d.SetId(backupList[0]["BACKUPID"].(string))
//	d.SetId(backupList[0]["description"].(string))
//	if err := d.Set("backups", backupList); err != nil {
//		return fmt.Errorf("Error setting `backups`: %#v", err)
//	}
//
//	return nil
//}
