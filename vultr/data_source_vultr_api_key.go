package vultr

//
//import (
//	"context"
//	"fmt"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
//)
//
//func dataSourceVultrApi() *schema.Resource {
//	return &schema.Resource{
//		Read: dataSourceVultrApiRead,
//		Schema: map[string]*schema.Schema{
//			"name": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"email": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"acl": {
//				Type:     schema.TypeList,
//				Elem:     &schema.Schema{Type: schema.TypeString},
//				Computed: true,
//			},
//		},
//	}
//}
//
//func dataSourceVultrApiRead(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	api, err := client.API.GetInfo(context.Background())
//
//	if err != nil {
//		return fmt.Errorf("Error getting api information: %v", err)
//	}
//
//	d.SetId(api.Email)
//	d.Set("name", api.Name)
//	d.Set("email", api.Email)
//	d.Set("acl", api.ACL)
//	return nil
//}
