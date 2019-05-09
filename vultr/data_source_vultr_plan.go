package vultr

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrPlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrPlanRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_gb": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"price_per_month": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"windows": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"plan_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"deprecated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	plans, err := client.Plan.GetList(context.Background(), "")

	if err != nil {
		return fmt.Errorf("Error getting plans: %v", err)
	}

	planList := []govultr.Plan{}
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

	d.SetId(strconv.Itoa(planList[0].VpsID))
	d.Set("name", planList[0].Name)
	d.Set("vcpu_count", planList[0].VCpus)
	d.Set("ram", planList[0].RAM)
	d.Set("disk", planList[0].Disk)
	d.Set("bandwidth", planList[0].Bandwidth)
	d.Set("bandwidth_gb", planList[0].BandwidthGB)
	d.Set("price_per_month", planList[0].Price)
	d.Set("windows", planList[0].Windows)
	d.Set("plan_type", planList[0].PlanType)
	d.Set("available_locations", planList[0].Regions)
	d.Set("deprecated", planList[0].Deprecated)
	return nil
}
