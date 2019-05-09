package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrRegionRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"continent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ddos_protection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"block_storage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"regioncode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrRegionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	regions, err := client.Region.GetList(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting regions: %v", err)
	}

	regionList := []govultr.Region{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, a := range regions {
		// we need convert the a struct INTO a map so we can easily manipulate the data here
		sm, err := structToMap(a)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			regionList = append(regionList, a)
		}
	}

	if len(regionList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(regionList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(regionList[0].RegionID)
	d.Set("name", regionList[0].Name)
	d.Set("country", regionList[0].Country)
	d.Set("continent", regionList[0].Continent)
	d.Set("state", regionList[0].State)
	d.Set("ddos_protection", regionList[0].Ddos)
	d.Set("block_storage", regionList[0].BlockStorage)
	d.Set("regioncode", regionList[0].RegionCode)
	return nil
}
