package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrBareMetalPlan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrBareMetalPlanRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_threads": {
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
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"disk_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrBareMetalPlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var planList []govultr.BareMetalPlan
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		plans, meta, err := client.Plan.ListBareMetal(ctx, options)

		if err != nil {
			return diag.Errorf("Error getting bare metal plans: %v", err)
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
	if err := d.Set("cpu_count", planList[0].CPUCount); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `cpu_count` read value: %v", err)
	}
	if err := d.Set("cpu_model", planList[0].CPUModel); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `cpu_model` read value: %v", err)
	}
	if err := d.Set("cpu_threads", planList[0].CPUThreads); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `cpu_threads` read value: %v", err)
	}
	if err := d.Set("ram", planList[0].RAM); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `ram` read value: %v", err)
	}
	if err := d.Set("disk", planList[0].Disk); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `disk` read value: %v", err)
	}
	if err := d.Set("bandwidth", planList[0].Bandwidth); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `bandwidth` read value: %v", err)
	}
	if err := d.Set("monthly_cost", planList[0].MonthlyCost); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `monthly_cost` read value: %v", err)
	}
	if err := d.Set("type", planList[0].Type); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `type` read value: %v", err)
	}
	if err := d.Set("disk_count", planList[0].DiskCount); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `disk_count` read value: %v", err)
	}
	if err := d.Set("locations", planList[0].Locations); err != nil {
		return diag.Errorf("unable to set bare_metal_plan `locations` read value: %#v", err)
	}

	return nil
}
