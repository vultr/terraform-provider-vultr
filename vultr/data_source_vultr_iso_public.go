package vultr

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
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
			"md5sum": {
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

	isoList := []govultr.PublicISO{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		iso, meta, err := client.ISO.ListPublic(context.Background(), options)
		if err != nil {
			return fmt.Errorf("Error getting applications: %v", err)
		}

		for _, i := range iso {
			sm, err := structToMap(i)

			if err != nil {
				return err
			}

			if filterLoop(f, sm) {
				isoList = append(isoList, i)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(isoList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(isoList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(isoList[0].ID)
	d.Set("description", isoList[0].Description)
	d.Set("name", isoList[0].Name)
	d.Set("md5sum", isoList[0].MD5Sum)
	return nil
}
