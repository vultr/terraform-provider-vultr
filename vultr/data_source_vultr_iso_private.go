package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrIsoPrivate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrIsoPrivateRead,
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

func dataSourceVultrIsoPrivateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOK := d.GetOk("filter")
	if !filtersOK {
		return diag.Errorf("issue with filter: %v", filtersOK)
	}

	var isoList []govultr.ISO
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		iso, meta,_, err := client.ISO.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting isos: %v", err)
		}

		for _, i := range iso {
			sm, err := structToMap(i)

			if err != nil {
				return diag.FromErr(err)
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
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(isoList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(isoList[0].ID)
	if err := d.Set("date_created", isoList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set iso_private `date_created` read value: %v", err)
	}
	if err := d.Set("filename", isoList[0].FileName); err != nil {
		return diag.Errorf("unable to set iso_private `filename` read value: %v", err)
	}
	if err := d.Set("size", isoList[0].Size); err != nil {
		return diag.Errorf("unable to set iso_private `size` read value: %v", err)
	}
	if err := d.Set("md5sum", isoList[0].MD5Sum); err != nil {
		return diag.Errorf("unable to set iso_private `md5sum` read value: %v", err)
	}
	if err := d.Set("sha512sum", isoList[0].SHA512Sum); err != nil {
		return diag.Errorf("unable to set iso_private `sha512sum` read value: %v", err)
	}
	if err := d.Set("status", isoList[0].Status); err != nil {
		return diag.Errorf("unable to set iso_private `status` read value: %v", err)
	}
	return nil
}
