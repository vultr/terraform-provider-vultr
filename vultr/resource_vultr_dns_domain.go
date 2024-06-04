package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrDNSDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDNSDomainCreate,
		ReadContext:   resourceVultrDNSDomainRead,
		UpdateContext: resourceVultrDNSDomainUpdate,
		DeleteContext: resourceVultrDNSDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"dns_sec": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "disabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrDNSDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	domainReq := &govultr.DomainReq{
		Domain: d.Get("domain").(string),
		DNSSec: d.Get("dns_sec").(string),
	}

	if ip, ok := d.GetOk("ip"); ok {
		domainReq.IP = ip.(string)
	}

	log.Print("[INFO] Creating domain")

	domain, _, err := client.Domain.Create(ctx, domainReq)
	if err != nil {
		return diag.Errorf("error while creating domain : %s", err)
	}

	d.SetId(domain.Domain)

	return resourceVultrDNSDomainRead(ctx, d, meta)
}

func resourceVultrDNSDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	domain, _, err := client.Domain.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid domain") {
			tflog.Warn(ctx, fmt.Sprintf("Removing domain (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting domains : %v", err)
	}

	if err := d.Set("domain", domain.Domain); err != nil {
		return diag.Errorf("unable to set resource dns_domain `domain` read value: %v", err)
	}
	if err := d.Set("date_created", domain.DateCreated); err != nil {
		return diag.Errorf("unable to set resource dns_domain `date_created` read value: %v", err)
	}
	if err := d.Set("dns_sec", domain.DNSSec); err != nil {
		return diag.Errorf("unable to set resource dns_domain `dns_sec` read value: %v", err)
	}

	return nil
}

func resourceVultrDNSDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updated domain (%s)", d.Id())
	if err := client.Domain.Update(ctx, d.Id(), d.Get("dns_sec").(string)); err != nil {
		return diag.Errorf("error updating domain %s: %v", d.Id(), err)
	}

	return resourceVultrDNSDomainRead(ctx, d, meta)
}

func resourceVultrDNSDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting domain (%s)", d.Id())
	if err := client.Domain.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying domain %s: %v", d.Id(), err)
	}

	return nil
}
