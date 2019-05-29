package vultr

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrIsoPrivate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrIsoPrivateRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"md5sum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sha512sum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrIsoPrivateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOK := d.GetOk("filter")

	if !filtersOK {
		return fmt.Errorf("issue with filter: %v", filtersOK)
	}

	iso, err := client.ISO.List(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting applications: %v", err)
	}

	isoList := []govultr.ISO{}
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
	d.Set("date_created", isoList[0].DateCreated)
	d.Set("filename", isoList[0].FileName)
	d.Set("size", isoList[0].Size)
	d.Set("md5sum", isoList[0].MD5Sum)
	d.Set("sha512sum", isoList[0].SHA512Sum)
	d.Set("status", isoList[0].Status)
	return nil
}
