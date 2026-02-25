package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrOIDCIssuer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrOIDCIssuerRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"n": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"e": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jwks_fetched_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jwks_expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrOIDCIssuerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var issuerList []govultr.OIDCIssuer
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	issuers, _, err := client.OIDC.ListOIDCIssuers(ctx)
	if err != nil {
		return diag.Errorf("error getting oidc issuers: %v", err)
	}

	for i := range issuers {
		sm, err := structToMap(issuers[i])

		if err != nil {
			return diag.FromErr(err)
		}

		if filterLoop(f, sm) {
			issuerList = append(issuerList, issuers[i])
		}
	}

	if len(issuerList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(issuerList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(issuerList[0].ID)

	if err := d.Set("source", issuerList[0].Source); err != nil {
		return diag.Errorf("unable to set oidc issuer `source` read value: %v", err)
	}
	if err := d.Set("uri", issuerList[0].URI); err != nil {
		return diag.Errorf("unable to set oidc issuer `uri` read value: %v", err)
	}
	if err := d.Set("n", issuerList[0].N); err != nil {
		return diag.Errorf("unable to set oidc issuer `n` read value: %v", err)
	}
	if err := d.Set("e", issuerList[0].E); err != nil {
		return diag.Errorf("unable to set oidc issuer `e` read value: %v", err)
	}
	if err := d.Set("alg", issuerList[0].ALG); err != nil {
		return diag.Errorf("unable to set oidc issuer `alg` read value: %v", err)
	}
	if err := d.Set("use", issuerList[0].USE); err != nil {
		return diag.Errorf("unable to set oidc issuer `use` read value: %v", err)
	}
	if err := d.Set("jwks_fetched_date", issuerList[0].JWKSFetchedDate); err != nil {
		return diag.Errorf("unable to set oidc issuer `jwks_fetched_date` read value: %v", err)
	}
	if err := d.Set("jwks_expiry_date", issuerList[0].JWKSExpiryDate); err != nil {
		return diag.Errorf("unable to set oidc issuer `jwks_expiry_date` read value: %v", err)
	}

	return nil
}
