package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrPlan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrPlanRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vcpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"monthly_cost": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gpu_vram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gpu_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"disk_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrPlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	planList := []govultr.Plan{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		plans, meta,_, err := client.Plan.List(ctx, "", options)
		if err != nil {
			return diag.Errorf("Error getting plans: %v", err)
		}

		for _, a := range plans {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(a)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				planList = append(planList, a)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(planList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(planList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(planList[0].ID)
	if err := d.Set("vcpu_count", planList[0].VCPUCount); err != nil {
		return diag.Errorf("unable to set plan `vcpu_count` read value: %v", err)
	}
	if err := d.Set("ram", planList[0].RAM); err != nil {
		return diag.Errorf("unable to set plan `ram` read value: %v", err)
	}
	if err := d.Set("disk", planList[0].Disk); err != nil {
		return diag.Errorf("unable to set plan `disk` read value: %v", err)
	}
	if err := d.Set("bandwidth", planList[0].Bandwidth); err != nil {
		return diag.Errorf("unable to set plan `bandwidth` read value: %v", err)
	}
	if err := d.Set("monthly_cost", planList[0].MonthlyCost); err != nil {
		return diag.Errorf("unable to set plan `monthly_cost` read value: %v", err)
	}
	if err := d.Set("disk_count", planList[0].DiskCount); err != nil {
		return diag.Errorf("unable to set plan `disk_count` read value: %v", err)
	}
	if err := d.Set("type", planList[0].Type); err != nil {
		return diag.Errorf("unable to set plan `type` read value: %v", err)
	}
	if err := d.Set("gpu_vram", planList[0].GPUVRAM); err != nil {
		return diag.Errorf("unable to set plan `gpu_vram` read value: %v", err)
	}
	if err := d.Set("gpu_type", planList[0].GPUType); err != nil {
		return diag.Errorf("unable to set plan `gpu_type` read value: %v", err)
	}
	if err := d.Set("locations", planList[0].Locations); err != nil {
		return diag.Errorf("unable to set plan `available_locations` read value: %v", err)
	}
	return nil
}
