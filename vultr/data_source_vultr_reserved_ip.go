package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrReservedIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrReservedIPRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrReservedIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	ipList := []govultr.ReservedIP{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		ips, meta, err := client.ReservedIP.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting list of reserved ips: %v", err)
		}

		for _, i := range ips {
			sm, err := structToMap(i)
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				ipList = append(ipList, i)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(ipList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(ipList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(ipList[0].ID)
	if err := d.Set("region", ipList[0].Region); err != nil {
		return diag.Errorf("unable to set reserved_ip `region` read value: %v", err)
	}
	if err := d.Set("ip_type", ipList[0].IPType); err != nil {
		return diag.Errorf("unable to set reserved_ip `ip_type` read value: %v", err)
	}
	if err := d.Set("subnet", ipList[0].Subnet); err != nil {
		return diag.Errorf("unable to set reserved_ip `subnet` read value: %v", err)
	}
	if err := d.Set("subnet_size", ipList[0].SubnetSize); err != nil {
		return diag.Errorf("unable to set reserved_ip `subnet_size` read value: %v", err)
	}
	if err := d.Set("label", ipList[0].Label); err != nil {
		return diag.Errorf("unable to set reserved_ip `label` read value: %v", err)
	}
	if err := d.Set("instance_id", ipList[0].InstanceID); err != nil {
		return diag.Errorf("unable to set reserved_ip `instance_id` read value: %v", err)
	}
	return nil
}
