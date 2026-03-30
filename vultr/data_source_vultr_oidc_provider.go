package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrOIDCProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrOIDCProviderRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrOIDCProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var provList []govultr.OIDCProvider
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	provs, _, err := client.OIDC.ListOIDCProviders(ctx)
	if err != nil {
		return diag.Errorf("error getting oidc providers: %v", err)
	}

	for i := range provs {
		sm, err := structToMap(provs[i])

		if err != nil {
			return diag.FromErr(err)
		}

		if filterLoop(f, sm) {
			provList = append(provList, provs[i])
		}
	}

	if len(provList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(provList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(provList[0].ID)

	if err := d.Set("issuer_id", provList[0].IssuerID); err != nil {
		return diag.Errorf("unable to set oidc provider `issuer_id` read value: %v", err)
	}
	if err := d.Set("name", provList[0].Name); err != nil {
		return diag.Errorf("unable to set oidc provider `name` read value: %v", err)
	}

	return nil
}
