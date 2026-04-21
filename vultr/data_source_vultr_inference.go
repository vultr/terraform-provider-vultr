package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrInference() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrInferenceRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrInferenceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var inferenceList []govultr.Inference
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	inferenceSubs, _, err := client.Inference.List(ctx)
	if err != nil {
		return diag.Errorf("error getting inference subscriptions: %v", err)
	}

	for s := range inferenceSubs {
		// we need convert the a struct INTO a map so we can easily manipulate the data here
		sm, err := structToMap(inferenceSubs[s])

		if err != nil {
			return diag.FromErr(err)
		}

		if filterLoop(f, sm) {
			inferenceList = append(inferenceList, inferenceSubs[s])
		}
	}

	if len(inferenceList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(inferenceList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(inferenceList[0].ID)

	if err := d.Set("date_created", inferenceList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set inference `date_created` read value: %v", err)
	}

	if err := d.Set("label", inferenceList[0].Label); err != nil {
		return diag.Errorf("unable to set inference `label` read value: %v", err)
	}

	if err := d.Set("api_key", inferenceList[0].APIKey); err != nil {
		return diag.Errorf("unable to set inference `api_key` read value: %v", err)
	}

	return nil
}
