package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrOIDCToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOIDCTokenCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"grant_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redirect_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"refresh_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_seconds": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrOIDCTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Print("[INFO] Creating oidc token")

	token, _, err := client.OIDC.CreateOIDCToken(ctx, &govultr.OIDCTokenReq{
		GrantType:    d.Get("grant_type").(string),
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		Code:         d.Get("code").(string),
		RedirectURI:  d.Get("redirect_uri").(string),
		RefreshToken: d.Get("refresh_token").(string),
	})
	if err != nil {
		return diag.Errorf("error while creating oidc token : %s", err)
	}

	d.SetId(token.AccessToken)

	if err := d.Set("access_token", token.AccessToken); err != nil {
		return diag.Errorf("unable to set instance `access_token` create value: %v", err)
	}
	if err := d.Set("token_type", token.TokenType); err != nil {
		return diag.Errorf("unable to set instance `token_type` create value: %v", err)
	}
	if err := d.Set("expires_seconds", token.ExpiresSeconds); err != nil {
		return diag.Errorf("unable to set instance `expires_seconds` create value: %v", err)
	}
	if err := d.Set("id_token", token.IDToken); err != nil {
		return diag.Errorf("unable to set instance `id_token` create value: %v", err)
	}
	if err := d.Set("scope", token.Scope); err != nil {
		return diag.Errorf("unable to set instance `scope` create value: %v", err)
	}

	return nil
}
