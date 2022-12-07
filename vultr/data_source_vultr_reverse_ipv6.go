package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrReverseIPV6() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrReverseIPV6Read,
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

func dataSourceVultrReverseIPV6Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	filters, ok := d.GetOk("filter")
	if !ok {
		return diag.Errorf("error getting filter")
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
			servers, meta, err := client.Instance.List(ctx, options)
			if err != nil {
				return diag.Errorf("Error getting servers: %v", err)
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
		reverseIPV6s, err := client.Instance.ListReverseIPv6(ctx, instanceID)
		if err != nil {
			return diag.Errorf("error getting reverse IPv6s: %v", err)
		}

		for _, reverseIPV6 := range reverseIPV6s {
			m, err := structToMap(reverseIPV6)
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(filter, m) {
				if result != nil {
					return diag.Errorf("your search returned too many results - please refine your search to be more specific")
				}

				result = &reverseIPV6
				resultInstanceID = instanceID
			}
		}
	}

	if result == nil {
		return diag.Errorf("No results were found")
	}

	d.SetId(result.IP)
	if err := d.Set("instance_id", resultInstanceID); err != nil {
		return diag.Errorf("unable to set reverse_ipv6 `instance_id` read value: %v", err)
	}
	if err := d.Set("ip", result.IP); err != nil {
		return diag.Errorf("unable to set reverse_ipv6 `ip` read value: %v", err)
	}
	if err := d.Set("reverse", result.Reverse); err != nil {
		return diag.Errorf("unable to set reverse_ipv6 `reverse` read value: %v", err)
	}

	return nil
}
