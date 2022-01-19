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
		ips, meta, err := client.ReservedIP.List(context.Background(), options)
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
	d.Set("region", ipList[0].Region)
	d.Set("ip_type", ipList[0].IPType)
	d.Set("subnet", ipList[0].Subnet)
	d.Set("subnet_size", ipList[0].SubnetSize)
	d.Set("label", ipList[0].Label)
	d.Set("instance_id", ipList[0].InstanceID)
	return nil
}
