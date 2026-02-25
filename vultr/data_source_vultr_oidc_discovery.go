package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVultrOIDCDiscovery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrOIDCDiscoveryRead,
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorize_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jwks_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_info_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"response_types_supported": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subject_types_supported": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"id_token_values_supported": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"scopes_supported": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"claims_supported": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"grant_types_supported": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"token_endpoint_auth_methods_supported": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVultrOIDCDiscoveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	doc, _, err := client.OIDC.DiscoveryOIDC(ctx, d.Get("provider_id").(string))
	if err != nil {
		return diag.Errorf("error getting oidc discovery document: %v", err)
	}

	if err := d.Set("issuer", doc.Issuer); err != nil {
		return diag.Errorf("unable to set oidc discovery `issuer` read value: %v", err)
	}
	if err := d.Set("authorize_endpoint", doc.AuthorizeEndpoint); err != nil {
		return diag.Errorf("unable to set oidc discovery `authorize_endpoint` read value: %v", err)
	}
	if err := d.Set("token_endpoint", doc.TokenEndpoint); err != nil {
		return diag.Errorf("unable to set oidc discovery `token_endpoint` read value: %v", err)
	}
	if err := d.Set("jwks_uri", doc.JWKSURI); err != nil {
		return diag.Errorf("unable to set oidc discovery `jwks_uri` read value: %v", err)
	}
	if err := d.Set("user_info_endpoint", doc.UserInfoEndpoint); err != nil {
		return diag.Errorf("unable to set oidc discovery `user_info_endpoint` read value: %v", err)
	}
	if err := d.Set("response_types_supported", doc.ResponseTypesSupported); err != nil {
		return diag.Errorf("unable to set oidc discovery `response_types_supported` read value: %v", err)
	}
	if err := d.Set("subject_types_supported", doc.SubjectTypesSupported); err != nil {
		return diag.Errorf("unable to set oidc discovery `subject_types_supported` read value: %v", err)
	}
	if err := d.Set("id_token_signing_values_supported", doc.IDTokenSigningValuesSupported); err != nil {
		return diag.Errorf("unable to set oidc discovery `id_token_signing_values_supported` read value: %v", err)
	}
	if err := d.Set("scopes_supported", doc.ScopesSupported); err != nil {
		return diag.Errorf("unable to set oidc discovery `scopes_supported` read value: %v", err)
	}
	if err := d.Set("claims_supported", doc.ClaimsSupported); err != nil {
		return diag.Errorf("unable to set oidc discovery `claims_supported` read value: %v", err)
	}
	if err := d.Set("grant_types_supported", doc.GrantTypesSupported); err != nil {
		return diag.Errorf("unable to set oidc discovery `grant_types_supported` read value: %v", err)
	}
	if err := d.Set("token_endpoint_auth_methods_supported", doc.TokenEndpointAuthMethodsSupported); err != nil {
		return diag.Errorf("unable to set oidc discovery `token_endpoint_auth_methods_supported` read value: %v", err)
	}

	return nil
}
