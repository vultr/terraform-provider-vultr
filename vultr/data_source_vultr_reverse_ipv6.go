package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrReverseIPV6() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrReverseIPV6Read,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reverse": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrReverseIPV6Read(d *schema.ResourceData, meta interface{}) error {
	filters, ok := d.GetOk("filter")
	if !ok {
		return fmt.Errorf("error getting filter")
	}

	var instanceIDs []string

	for _, filter := range filters.(*schema.Set).List() {
		filterMap := filter.(map[string]interface{})

		name := filterMap["name"]
		values := filterMap["values"].([]interface{})

		if name == "instance_id" {
			for _, value := range values {
				instanceIDs = append(instanceIDs, value.(string))
			}
		}

		if name == "ip" {
			for i, value := range values {
				values[i] = value.(string)
			}
		}
	}

	client := meta.(*Client).govultrClient()

	// If the data source is not being filtered by `instance_id`, consider all
	// servers
	options := &govultr.ListOptions{}
	if len(instanceIDs) == 0 {
		for {
			servers, meta, err := client.Instance.List(context.Background(), options)
			if err != nil {
				return fmt.Errorf("Error getting servers: %v", err)
			}

			for _, server := range servers {
				// Consider servers with at least one assigned IPv6 subnet
				if server.V6MainIP != "" {
					instanceIDs = append(instanceIDs, server.ID)
				}
			}
			if meta.Links.Next == "" {
				break
			} else {
				options.Cursor = meta.Links.Next
				continue
			}
		}
	}

	filter := buildVultrDataSourceFilter(filters.(*schema.Set))
	var result *govultr.ReverseIP
	resultInstanceID := ""

	for _, instanceID := range instanceIDs {
		reverseIPV6s, err := client.Instance.ListReverseIPv6(context.Background(), instanceID)
		if err != nil {
			return fmt.Errorf("error getting reverse IPv6s: %v", err)
		}

		for _, reverseIPV6 := range reverseIPV6s {
			m, err := structToMap(reverseIPV6)
			if err != nil {
				return err
			}

			if filterLoop(filter, m) {
				if result != nil {
					return fmt.Errorf("your search returned too many results - please refine your search to be more specific")
				}

				result = &reverseIPV6
				resultInstanceID = instanceID
			}
		}
	}

	if result == nil {
		return errors.New("No results were found")
	}

	d.SetId(result.IP)
	d.Set("instance_id", resultInstanceID)
	d.Set("ip", result.IP)
	d.Set("reverse", result.Reverse)

	return nil
}
