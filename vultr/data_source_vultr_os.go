package vultr

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrOS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrOSRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"family": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrOSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	osList := []govultr.OS{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		os, meta, err := client.OS.List(context.Background(), options)
		if err != nil {
			return fmt.Errorf("error getting os list: %v", err)
		}

		for _, o := range os {
			sm, err := structToMap(o)

			if err != nil {
				return err
			}

			if filterLoop(f, sm) {
				osList = append(osList, o)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(osList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(osList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(strconv.Itoa(osList[0].ID))
	d.Set("name", osList[0].Name)
	d.Set("arch", osList[0].Arch)
	d.Set("family", osList[0].Family)
	return nil
}
