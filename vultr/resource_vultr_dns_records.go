package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrDnsRecordCreate,
		Read:   resourceVultrDnsRecordRead,
		Update: resourceVultrDnsRecordUpdate,
		Delete: resourceVultrDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVultrDnsRecordImport,
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

func resourceVultrDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
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
	record, err := client.DomainRecord.Create(context.Background(), d.Get("domain").(string), recordReq)
	if err != nil {
		return fmt.Errorf("error creating DNS record : %v", err)
	}

	d.SetId(record.ID)
	return resourceVultrDnsRecordRead(d, meta)
}
func resourceVultrDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	record, err := client.DomainRecord.Get(context.Background(), d.Get("domain").(string), d.Id())
	if err != nil {
		log.Printf("[WARN] DNS Record %s not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("domain", d.Get("domain").(string))
	d.Set("type", record.Type)
	d.Set("name", record.Name)
	d.Set("data", record.Data)
	d.Set("priority", record.Priority)
	d.Set("ttl", record.TTL)
	return nil
}
func resourceVultrDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating DNS record: %s", d.Id())

	p := d.Get("priority").(int)
	record := &govultr.DomainRecordReq{
		Data:     d.Get("data").(string),
		Name:     d.Get("name").(string),
		TTL:      d.Get("ttl").(int),
		Priority: &p,
	}

	if err := client.DomainRecord.Update(context.Background(), d.Get("domain").(string), d.Id(), record); err != nil {
		return fmt.Errorf("error updating DNS record %s : %v", d.Id(), err)
	}

	return resourceVultrDnsRecordRead(d, meta)
}

func resourceVultrDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting DNS record: %s", d.Id())
	if err := client.DomainRecord.Delete(context.Background(), d.Get("domain").(string), d.Id()); err != nil {
		return fmt.Errorf("error deleting dns record %s : %v", d.Id(), err)
	}

	return nil
}

func resourceVultrDnsRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Client).govultrClient()

	importID := d.Id()
	commaIdx := strings.IndexByte(importID, ',')
	if commaIdx == -1 {
		return nil, fmt.Errorf(`invalid import format, expected "domain,resourceID"`)
	}
	domain, recordID := importID[:commaIdx], importID[commaIdx+1:]

	record, err := client.DomainRecord.Get(context.Background(), domain, recordID)
	if err != nil {
		return nil, fmt.Errorf("DNS record not found for domain %s", domain)
	}

	d.SetId(record.ID)
	d.Set("domain", domain)
	return []*schema.ResourceData{d}, nil
}
