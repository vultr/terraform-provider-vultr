package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrObjectStorage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrObjectStorageRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"s3_secret_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceVultrObjectStorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOK := d.GetOk("filter")
	if !filtersOK {
		return diag.Errorf("issue with filter: %v", filtersOK)
	}

	objStoreList := []govultr.ObjectStorage{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		objectStorages, meta, err := client.ObjectStorage.List(context.Background(), options)
		if err != nil {
			return diag.Errorf("error getting object storage list: %v", filtersOK)
		}

		for _, n := range objectStorages {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(n)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				objStoreList = append(objStoreList, n)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}
	if len(objStoreList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(objStoreList) < 1 {
		return diag.Errorf("no results were found")
	}
	d.SetId(objStoreList[0].ID)
	d.Set("date_created", objStoreList[0].DateCreated)
	d.Set("cluster_id", objStoreList[0].ObjectStoreClusterID)
	d.Set("label", objStoreList[0].Label)
	d.Set("region", objStoreList[0].Region)
	d.Set("location", objStoreList[0].Location)
	d.Set("status", objStoreList[0].Status)
	d.Set("s3_hostname", objStoreList[0].S3Hostname)
	d.Set("s3_access_key", objStoreList[0].S3AccessKey)
	d.Set("s3_secret_key", objStoreList[0].S3SecretKey)
	return nil
}
