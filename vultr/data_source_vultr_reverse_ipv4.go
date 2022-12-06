package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrReverseIPV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrReverseIPV4Read,
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
			"netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrReverseIPV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("error getting filter: %v", filtersOk)
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

	// If the data source is not being filtered by `instance_id`, consider all instances
	options := &govultr.ListOptions{}
	if len(instanceIDs) == 0 {
		for {
			servers, meta, err := client.Instance.List(ctx, options)
			if err != nil {
				return diag.Errorf("error getting servers: %v", err)
			}

			for _, server := range servers {
				instanceIDs = append(instanceIDs, server.ID)
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
	var result *govultr.IPv4
	resultInstanceID := ""

	for _, instanceID := range instanceIDs {
		ipv4s, _, err := client.Instance.ListIPv4(ctx, instanceID, nil)
		if err != nil {
			return diag.Errorf("error getting IPv4s: %v", err)
		}

		for _, ipv4 := range ipv4s {
			m, err := structToMap(ipv4)
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(filter, m) {
				if result != nil {
					return diag.Errorf("your search returned too many results - please refine your search to be more specific")
				}

				result = &ipv4
				resultInstanceID = instanceID
			}
		}
	}

	if result == nil {
		return diag.Errorf("no results were found")
	}

	d.SetId(result.IP)
	if err := d.Set("instance_id", resultInstanceID); err != nil {
		return diag.Errorf("unable to set reverse_ipv4 `instance_id` read value: %v", err)
	}
	if err := d.Set("ip", result.IP); err != nil {
		return diag.Errorf("unable to set reverse_ipv4 `ip` read value: %v", err)
	}
	if err := d.Set("reverse", result.Reverse); err != nil {
		return diag.Errorf("unable to set reverse_ipv4 `reverse` read value: %v", err)
	}
	if err := d.Set("netmask", result.Netmask); err != nil {
		return diag.Errorf("unable to set reverse_ipv4 `netmask` read value: %v", err)
	}
	if err := d.Set("gateway", result.Gateway); err != nil {
		return diag.Errorf("unable to set reverse_ipv4 `gateway` read value: %v", err)
	}

	return nil
}
