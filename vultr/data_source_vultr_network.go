package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrNetworkRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
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
	}
}

func dataSourceVultrNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	networks, err := client.Network.List(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting networks: %v", err)
	}

	networkList := []govultr.Network{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, n := range networks {
		// we need convert the a struct INTO a map so we can easily manipulate the data here
		sm, err := structToMap(n)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			networkList = append(networkList, n)
		}
	}

	if len(networkList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(networkList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(networkList[0].NetworkID)
	d.Set("region_id", networkList[0].RegionID)
	d.Set("description", networkList[0].Description)
	d.Set("date_created", networkList[0].DateCreated)
	d.Set("cidr_block", fmt.Sprintf("%s/%d", networkList[0].V4Subnet, networkList[0].V4SubnetMask))

	return nil
}
