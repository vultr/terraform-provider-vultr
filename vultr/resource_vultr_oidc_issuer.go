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

func resourceVultrOIDCIssuer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrOIDCIssuerCreate,
		ReadContext:   resourceVultrOIDCIssuerRead,
		DeleteContext: resourceVultrOIDCIssuerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"uri": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"kty": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"kid": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"n": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"e": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"alg": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"use": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
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

func resourceVultrOIDCIssuerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Print("[INFO] Creating oidc issuer")

	issu, _, err := client.OIDC.CreateOIDCIssuer(ctx, &govultr.OIDCIssuerReq{
		Issuer: govultr.OIDCIssuerReqDetail{
			Source:   d.Get("source").(string),
			SourceID: d.Get("source_id").(string),
			URI:      d.Get("uri").(string),
			KID:      d.Get("kid").(string),
			KTY:      d.Get("kty").(string),
			N:        d.Get("n").(string),
			E:        d.Get("e").(string),
			ALG:      d.Get("alg").(string),
			USE:      d.Get("use").(string),
		},
	})
	if err != nil {
		return diag.Errorf("error while creating oidc issuer : %s", err)
	}

	d.SetId(issu.ID)

	return resourceVultrOIDCIssuerRead(ctx, d, meta)
}

func resourceVultrOIDCIssuerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	issu, _, err := client.OIDC.GetOIDCIssuer(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Issuer Not Found") {
			tflog.Warn(ctx, fmt.Sprintf("Removing oidc issuer (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting oidc issuer : %v", err)
	}

	if err := d.Set("uri", issu.URI); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `uri` read value: %v", err)
	}
	if err := d.Set("kid", issu.KID); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `kid` read value: %v", err)
	}
	if err := d.Set("kty", issu.KTY); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `kty` read value: %v", err)
	}
	if err := d.Set("n", issu.N); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `n` read value: %v", err)
	}
	if err := d.Set("e", issu.E); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `e` read value: %v", err)
	}
	if err := d.Set("alg", issu.ALG); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `alg` read value: %v", err)
	}
	if err := d.Set("use", issu.USE); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `use` read value: %v", err)
	}
	if err := d.Set("jwks_fetched_date", issu.JWKSFetchedDate); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `jwks_fetched_date` read value: %v", err)
	}
	if err := d.Set("jwks_expiry_date", issu.JWKSExpiryDate); err != nil {
		return diag.Errorf("unable to set resource oidc issuer `jwks_expiry_date` read value: %v", err)
	}

	return nil
}

func resourceVultrOIDCIssuerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting oidc issuer (%s)", d.Id())
	if err := client.OIDC.DeleteOIDCIssuer(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting oidc issuer %s: %v", d.Id(), err)
	}

	return nil
}
