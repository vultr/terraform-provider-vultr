package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrPrivateNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrPrivateNetworkRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v4_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v4_subnet_mask": {
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
	}
}

func dataSourceVultrPrivateNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	var networkList []govultr.Network
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		networks, meta, err := client.Network.List(context.Background(), options)
		if err != nil {
			return fmt.Errorf("error getting networks: %v", err)
		}

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

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(networkList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(networkList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(networkList[0].NetworkID)
	d.Set("region", networkList[0].Region)
	d.Set("description", networkList[0].Description)
	d.Set("date_created", networkList[0].DateCreated)
	d.Set("v4_subnet", networkList[0].V4Subnet)
	d.Set("v4_subnet_mask", networkList[0].V4SubnetMask)

	return nil
}
