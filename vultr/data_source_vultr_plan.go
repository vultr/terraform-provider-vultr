package vultr

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrPlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrPlanRead,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	planList := []govultr.Plan{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		plans, meta, err := client.Plan.List(context.Background(), "", options)
		if err != nil {
			return fmt.Errorf("Error getting plans: %v", err)
		}

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

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(planList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(planList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(planList[0].ID)
	d.Set("vcpu_count", planList[0].VCPUCount)
	d.Set("ram", planList[0].Ram)
	d.Set("disk", planList[0].Disk)
	d.Set("bandwidth", planList[0].Bandwidth)
	d.Set("monthly_cost", planList[0].MonthlyCost)
	d.Set("type", planList[0].Type)
	if err := d.Set("locations", planList[0].Locations); err != nil {
		return fmt.Errorf("error setting `available_locations`: %#v", err)
	}
	return nil
}
