package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
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
				ValidateFunc: validateRecordType,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceVultrDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	data := d.Get("data").(string)
	domain := d.Get("domain").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	priority := d.Get("priority")
	ttl := d.Get("ttl").(int)

	log.Print("[INFO] Creating DNS record")
	err := client.DNSRecord.Create(context.Background(), domain, recordType, name, data, ttl, priority.(int))

	if err != nil {
		return fmt.Errorf("Error creating DNS record : %v", err)
	}

	// Grab Unique RecordID since create does not return it
	records, err := client.DNSRecord.List(context.Background(), domain)

	if err != nil {
		return fmt.Errorf("Error getting DNS records : %v", err)
	}

	for _, v := range records {
		if data == v.Data && recordType == v.Type && name == v.Name {
			d.SetId(strconv.Itoa(v.RecordID))
			return resourceVultrDnsRecordRead(d, meta)
		}
	}

	return fmt.Errorf("Error finding DNS record: %v", err)
}
func resourceVultrDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	records, err := client.DNSRecord.List(context.Background(), d.Get("domain").(string))

	if err != nil {
		return fmt.Errorf("Error getting DNS records for DNS Domain %s: %v", d.Get("domain").(string), err)
	}

	var record *govultr.DNSRecord
	for _, v := range records {
		if strconv.Itoa(v.RecordID) == d.Id() {
			record = &v
			break
		}
	}

	if record == nil {
		log.Printf("[WARN] DNS Record %s not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("type", record.Type)
	d.Set("name", record.Name)
	d.Set("data", record.Data)
	d.Set("priority", *record.Priority)
	d.Set("ttl", record.TTL)
	return nil
}
func resourceVultrDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Updating DNS record: %s", d.Id())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retreiving DNS record ID : %s", d.Id())
	}

	p := d.Get("priority").(int)
	record := &govultr.DNSRecord{
		RecordID: id,
		Data:     d.Get("data").(string),
		Name:     d.Get("name").(string),
		TTL:      d.Get("ttl").(int),
		Priority: &p,
	}

	err = client.DNSRecord.Update(context.Background(), d.Get("domain").(string), record)

	if err != nil {
		return fmt.Errorf("Error updating DNS record %s : %v", d.Id(), err)
	}

	return resourceVultrDnsRecordRead(d, meta)
}

func resourceVultrDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting DNS record: %s", d.Id())
	err := client.DNSRecord.Delete(context.Background(), d.Get("domain").(string), d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting dns record %s : %v", d.Id(), err)
	}

	return nil
}

func validateRecordType(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	recordTypes := []string{"A", "AAAA", "CNAME", "NS", "MX", "SRV", "TXT", "CAA", "SSHFP"}

	exists := func() bool {
		for _, i := range recordTypes {
			if v == i {
				return true
			}
		}
		return false
	}

	if !exists() {
		errs = append(errs, fmt.Errorf("the value %q given for %q is invalid", v, key))
	}

	return
}

func resourceVultrDnsRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Client).govultrClient()

	importID := d.Id()
	commaIdx := strings.IndexByte(importID, ',')
	if commaIdx == -1 {
		return nil, fmt.Errorf(`invalid import format, expected "domain,resourceID"`)
	}
	domain, recordID := importID[:commaIdx], importID[commaIdx+1:]

	records, err := client.DNSRecord.List(context.Background(), domain)
	if err != nil {
		return nil, fmt.Errorf("error getting DNS records for DNS Domain %s: %v", domain, err)
	}

	var record *govultr.DNSRecord
	for _, v := range records {
		if strconv.Itoa(v.RecordID) == recordID {
			record = &v
			break
		}
	}
	if record == nil {
		return nil, fmt.Errorf("DNS record %s not found for domain %s", recordID, domain)
	}

	d.SetId(recordID)
	d.Set("domain", domain)
	return []*schema.ResourceData{d}, nil
}
