package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v2"
)

func resourceVultrDNSRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDNSRecordCreate,
		ReadContext:   resourceVultrDNSRecordRead,
		UpdateContext: resourceVultrDNSRecordUpdate,
		DeleteContext: resourceVultrDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVultrDNSRecordImport,
		},
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME", "NS", "MX", "SRV", "TXT", "CAA", "SSHFP"}, false),
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
		},
	}
}

func resourceVultrDNSRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	p := d.Get("priority").(int)
	recordReq := &govultr.DomainRecordReq{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Data:     d.Get("data").(string),
		TTL:      d.Get("ttl").(int),
		Priority: &p,
	}

	log.Print("[INFO] Creating DNS record")
	record, err := client.DomainRecord.Create(ctx, d.Get("domain").(string), recordReq)
	if err != nil {
		return diag.Errorf("error creating DNS record : %v", err)
	}

	d.SetId(record.ID)
	return resourceVultrDNSRecordRead(ctx, d, meta)
}
func resourceVultrDNSRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	record, err := client.DomainRecord.Get(ctx, d.Get("domain").(string), d.Id())
	if err != nil {
		log.Printf("[WARN] DNS Record %s not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("domain", d.Get("domain").(string)); err != nil {
		return diag.Errorf("unable to set resource dns_records `domain` read value: %v", err)
	}
	if err := d.Set("type", record.Type); err != nil {
		return diag.Errorf("unable to set resource dns_records `type` read value: %v", err)
	}
	if err := d.Set("name", record.Name); err != nil {
		return diag.Errorf("unable to set resource dns_records `name` read value: %v", err)
	}
	if err := d.Set("data", record.Data); err != nil {
		return diag.Errorf("unable to set resource dns_records `data` read value: %v", err)
	}
	if err := d.Set("priority", record.Priority); err != nil {
		return diag.Errorf("unable to set resource dns_records `priority` read value: %v", err)
	}
	if err := d.Set("ttl", record.TTL); err != nil {
		return diag.Errorf("unable to set resource dns_records `ttl` read value: %v", err)
	}
	return nil
}
func resourceVultrDNSRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating DNS record: %s", d.Id())

	p := d.Get("priority").(int)
	record := &govultr.DomainRecordReq{
		Data:     d.Get("data").(string),
		Name:     d.Get("name").(string),
		TTL:      d.Get("ttl").(int),
		Priority: &p,
	}

	if err := client.DomainRecord.Update(ctx, d.Get("domain").(string), d.Id(), record); err != nil {
		return diag.Errorf("error updating DNS record %s : %v", d.Id(), err)
	}

	return resourceVultrDNSRecordRead(ctx, d, meta)
}

func resourceVultrDNSRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting DNS record: %s", d.Id())
	if err := client.DomainRecord.Delete(ctx, d.Get("domain").(string), d.Id()); err != nil {
		return diag.Errorf("error deleting dns record %s : %v", d.Id(), err)
	}

	return nil
}

func resourceVultrDNSRecordImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Client).govultrClient()

	importID := d.Id()
	commaIdx := strings.IndexByte(importID, ',')
	if commaIdx == -1 {
		return nil, fmt.Errorf(`invalid import format, expected "domain,resourceID"`)
	}
	domain, recordID := importID[:commaIdx], importID[commaIdx+1:]

	record, err := client.DomainRecord.Get(ctx, domain, recordID)
	if err != nil {
		return nil, fmt.Errorf("DNS record not found for domain %s", domain)
	}

	d.SetId(record.ID)
	if err := d.Set("domain", domain); err != nil {
		return nil, fmt.Errorf("unable to set resource dns_records `domain` import value: %v", err)
	}
	return []*schema.ResourceData{d}, nil
}
