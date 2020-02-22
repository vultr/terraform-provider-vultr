package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrBareMetalPlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrBareMetalPlanRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_tb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"price_per_month": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deprecated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"available_locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceVultrBareMetalPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	plans, err := client.Plan.GetBareMetalList(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting bare metal plans: %v", err)
	}

	planList := []govultr.BareMetalPlan{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, a := range plans {
		// we need convert the a struct INTO a map so we can easily manipulate the data here
		sm, err := structToMap(a)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			planList = append(planList, a)
		}
	}

	if len(planList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(planList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(planList[0].PlanID)
	d.Set("name", planList[0].Name)
	d.Set("cpu_count", planList[0].CPUs)
	d.Set("cpu_model", planList[0].CPUModel)
	d.Set("ram", planList[0].RAM)
	d.Set("disk", planList[0].Disk)
	d.Set("bandwidth_tb", planList[0].BandwidthTB)
	d.Set("price_per_month", planList[0].Price)
	d.Set("plan_type", planList[0].PlanType)

	if err := d.Set("available_locations", planList[0].Regions); err != nil {
		return fmt.Errorf("Error setting `available_locations`: %#v", err)
	}

	d.Set("deprecated", planList[0].Deprecated)
	return nil
}
