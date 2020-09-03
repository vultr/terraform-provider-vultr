package vultr

//
//import (
//	"context"
//	"errors"
//	"fmt"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
//	"github.com/vultr/govultr/v2"
//)
//
//func dataSourceVultrReservedIP() *schema.Resource {
//	return &schema.Resource{
//		Read: dataSourceVultrReservedIPRead,
//		Schema: map[string]*schema.Schema{
//			"filter": dataSourceFiltersSchema(),
//			"region_id": {
//				Type:     schema.TypeInt,
//				Computed: true,
//			},
//			"ip_type": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"subnet": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"subnet_size": {
//				Type:     schema.TypeInt,
//				Computed: true,
//			},
//			"label": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"attached_to_vps": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//		},
//	}
//}
//
//func dataSourceVultrReservedIPRead(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	filters, filtersOk := d.GetOk("filter")
//
//	if !filtersOk {
//		return fmt.Errorf("issue with filter: %v", filtersOk)
//	}
//
//	ips, err := client.ReservedIP.List(context.Background())
//	if err != nil {
//		return fmt.Errorf("Error getting applications: %v", err)
//	}
//
//	ipList := []govultr.ReservedIP{}
//
//	f := buildVultrDataSourceFilter(filters.(*schema.Set))
//
//	for _, i := range ips {
//		// we need convert the a struct INTO a map so we can easily manipulate the data here
//		sm, err := structToMap(i)
//
//		if err != nil {
//			return err
//		}
//
//		if filterLoop(f, sm) {
//			ipList = append(ipList, i)
//		}
//	}
//
//	if len(ipList) > 1 {
//		return errors.New("your search returned too many results. Please refine your search to be more specific")
//	}
//
//	if len(ipList) < 1 {
//		return errors.New("no results were found")
//	}
//
//	d.SetId(ipList[0].ReservedIPID)
//	d.Set("region_id", ipList[0].RegionID)
//	d.Set("ip_type", ipList[0].IPType)
//	d.Set("subnet", ipList[0].Subnet)
//	d.Set("subnet_size", ipList[0].SubnetSize)
//	d.Set("label", ipList[0].Label)
//	d.Set("attached_to_vps", ipList[0].AttachedID)
//	return nil
//}
