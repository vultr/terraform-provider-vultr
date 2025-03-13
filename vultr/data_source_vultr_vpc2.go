package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrVPC2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrVPC2Read,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		DeprecationMessage: "VPC2 is deprecated and will not be supported in a future release.  Use VPC instead",
	}
}

func dataSourceVultrVPC2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var vpcList []govultr.VPC2
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		vpcs, meta, _, err := client.VPC2.List(ctx, options) //nolint:staticcheck
		if err != nil {
			return diag.Errorf("error getting VPCs 2.0: %v", err)
		}

		for _, n := range vpcs {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(n)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				vpcList = append(vpcList, n)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(vpcList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(vpcList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(vpcList[0].ID)
	if err := d.Set("region", vpcList[0].Region); err != nil {
		return diag.Errorf("unable to set vpc2 `region` read value: %v", err)
	}
	if err := d.Set("description", vpcList[0].Description); err != nil {
		return diag.Errorf("unable to set vpc2 `description` read value: %v", err)
	}
	if err := d.Set("date_created", vpcList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set vpc2 `date_created` read value: %v", err)
	}
	if err := d.Set("ip_block", vpcList[0].IPBlock); err != nil {
		return diag.Errorf("unable to set vpc2 `ip_block` read value: %v", err)
	}
	if err := d.Set("prefix_length", vpcList[0].PrefixLength); err != nil {
		return diag.Errorf("unable to set vpc2 `prefix_length` read value: %v", err)
	}

	return nil
}
