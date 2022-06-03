package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
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
		plans, meta, err := client.Plan.List(ctx, "", options)
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
	d.Set("vcpu_count", planList[0].VCPUCount)
	d.Set("ram", planList[0].RAM)
	d.Set("disk", planList[0].Disk)
	d.Set("bandwidth", planList[0].Bandwidth)
	d.Set("monthly_cost", planList[0].MonthlyCost)
	d.Set("disk_count", planList[0].DiskCount)
	d.Set("type", planList[0].Type)
	d.Set("gpu_vram", planList[0].GPUVRAM)
	d.Set("gpu_type", planList[0].GPUType)
	if err := d.Set("locations", planList[0].Locations); err != nil {
		return diag.Errorf("error setting `available_locations`: %#v", err)
	}
	return nil
}
