package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
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
		return errors.New("Error getting filter")
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
	if len(instanceIDs) == 0 {
		servers, err := client.Server.List(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting servers: %v", err)
		}

		for _, server := range servers {
			// Consider servers with at least one assigned IPv6 subnet
			if len(server.V6Networks) > 0 {
				instanceIDs = append(instanceIDs, server.InstanceID)
			}
		}
	}

	filter := buildVultrDataSourceFilter(filters.(*schema.Set))

	var result *govultr.ReverseIPV6
	resultInstanceID := ""

	for _, instanceID := range instanceIDs {
		reverseIPV6s, err := client.Server.ListReverseIPV6(context.Background(), instanceID)
		if err != nil {
			return fmt.Errorf("Error getting reverse IPv6s: %v", err)
		}

		for _, reverseIPV6 := range reverseIPV6s {
			m, err := structToMap(reverseIPV6)
			if err != nil {
				return err
			}

			if filterLoop(filter, m) {
				if result != nil {
					return errors.New("Your search returned too many results. " +
						"Please refine your search to be more specific")
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
