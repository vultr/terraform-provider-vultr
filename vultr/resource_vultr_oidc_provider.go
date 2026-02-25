package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrOIDCProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOIDCProviderCreate,
		ReadContext:   resourceVultrOIDCProviderRead,
		DeleteContext: resourceVultrOIDCProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"issuer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrOIDCProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Print("[INFO] Creating oidc provider")

	prov, _, err := client.OIDC.CreateOIDCProvider(ctx, &govultr.OIDCProviderReq{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return diag.Errorf("error while creating oidc provider : %s", err)
	}

	d.SetId(prov.ID)

	return resourceVultrOIDCProviderRead(ctx, d, meta)
}

func resourceVultrOIDCProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	prov, _, err := client.OIDC.GetOIDCProvider(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Provider Not Found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing oidc provider (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting oidc provider : %v", err)
	}

	if err := d.Set("issuer_id", prov.IssuerID); err != nil {
		return diag.Errorf("unable to set resource oidc provider `issuer_id` read value: %v", err)
	}

	return nil
}

func resourceVultrOIDCProviderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting oidc provider (%s)", d.Id())
	if err := client.OIDC.DeleteOIDCProvider(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting oidc provider %s: %v", d.Id(), err)
	}

	return nil
}
