package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vultr/govultr/v2"
)

func resourceVultrDnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrDnsDomainCreate,
		Read:   resourceVultrDnsDomainRead,
		Update: resourceVultrDnsDomainUpdate,
		Delete: resourceVultrDnsDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceVultrDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	domainReq := &govultr.DomainReq{
		Domain: d.Get("domain").(string),
		DNSSec: d.Get("dns_sec").(string),
	}

	if ip, ok := d.GetOk("ip"); ok {
		domainReq.IP = ip.(string)
	}

	log.Print("[INFO] Creating domain")

	domain, err := client.Domain.Create(context.Background(), domainReq)
	if err != nil {
		return fmt.Errorf("error while creating domain : %s", err)
	}

	d.SetId(domain.Domain)

	return resourceVultrDnsDomainRead(d, meta)
}

func resourceVultrDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	domain, err := client.Domain.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting domains : %v", err)
	}

	d.Set("domain", domain.Domain)
	d.Set("date_created", domain.DateCreated)

	return nil
}

func resourceVultrDnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updated domain (%s)", d.Id())
	if err := client.Domain.Update(context.Background(), d.Id(), d.Get("dns_sec").(string)); err != nil {
		return fmt.Errorf("error updating domain %s: %v", d.Id(), err)
	}

	return resourceVultrDnsDomainRead(d, meta)
}

func resourceVultrDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting domain (%s)", d.Id())
	if err := client.Domain.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error destroying domain %s: %v", d.Id(), err)

	}

	return nil
}
