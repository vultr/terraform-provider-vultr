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
//func dataSourceVultrReverseIPV4() *schema.Resource {
//	return &schema.Resource{
//		Read: dataSourceVultrReverseIPV4Read,
//		Schema: map[string]*schema.Schema{
//			"filter": dataSourceFiltersSchema(),
//			"instance_id": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"ip": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"reverse": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//		},
//	}
//}
//
//func dataSourceVultrReverseIPV4Read(d *schema.ResourceData, meta interface{}) error {
//	filters, ok := d.GetOk("filter")
//	if !ok {
//		return errors.New("error getting filter")
//	}
//
//	var instanceIDs []string
//
//	for _, filter := range filters.(*schema.Set).List() {
//		filterMap := filter.(map[string]interface{})
//
//		name := filterMap["name"]
//		values := filterMap["values"].([]interface{})
//
//		if name == "instance_id" {
//			for _, value := range values {
//				instanceIDs = append(instanceIDs, value.(string))
//			}
//		}
//
//		if name == "ip" {
//			for i, value := range values {
//				values[i] = value.(string)
//			}
//		}
//	}
//
//	client := meta.(*Client).govultrClient()
//
//	// If the data source is not being filtered by `instance_id`, consider all
//	// servers
//	if len(instanceIDs) == 0 {
//		servers, err := client.Server.List(context.Background())
//		if err != nil {
//			return fmt.Errorf("error getting servers: %v", err)
//		}
//
//		for _, server := range servers {
//			instanceIDs = append(instanceIDs, server.InstanceID)
//		}
//	}
//
//	var result *govultr.IPV4
//	resultInstanceID := ""
//
//	for _, instanceID := range instanceIDs {
//		reverseIPv4s, err := client.Server.IPV4Info(context.Background(), instanceID, true)
//		if err != nil {
//			return fmt.Errorf("error getting reverse IPv4s: %v", err)
//		}
//
//		for _, reverseIPv4 := range reverseIPv4s {
//			result = &reverseIPv4
//			resultInstanceID = instanceID
//		}
//	}
//
//	if result == nil {
//		return errors.New("no results were found")
//	}
//
//	d.SetId(result.IP)
//	d.Set("instance_id", resultInstanceID)
//	d.Set("ip", result.IP)
//	d.Set("reverse", result.Reverse)
//
//	return nil
//}
