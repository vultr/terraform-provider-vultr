package vultr

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrIsoPublic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrIsoPublicRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrIsoPublicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOK := d.GetOk("filter")

	if !filtersOK {
		return fmt.Errorf("issue with filter: %v", filtersOK)
	}

	iso, err := client.ISO.GetPublicList(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting applications: %v", err)
	}

	isoList := []govultr.PublicISO{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, i := range iso {
		sm, err := structToMap(i)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			isoList = append(isoList, i)
		}
	}

	if len(isoList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(isoList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(strconv.Itoa(isoList[0].ISOID))
	d.Set("description", isoList[0].Description)
	d.Set("name", isoList[0].Name)
	return nil
}
