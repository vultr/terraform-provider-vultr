package vultr

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vultr/govultr"
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
			"server_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVultrDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()

	domain := d.Get("domain").(string)
	ip := d.Get("server_ip").(string)

	validIp := net.ParseIP(ip)
	if validIp == nil {
		return fmt.Errorf("The supplied IP address is invalid : %s", ip)
	}
	log.Print("[INFO] Creating DNS domain")

	err := client.DNSDomain.Create(context.Background(), domain, ip)

	if err != nil {
		return fmt.Errorf("Error while creating DNS domain : %s", err)
	}

	d.SetId(domain)

	return resourceVultrDnsDomainRead(d, meta)
}

func resourceVultrDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	domains, err := client.DNSDomain.List(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting domains ")
	}

	var domain *govultr.DNSDomain
	for i := range domains {
		if domains[i].Domain == d.Id() {
			domain = &domains[i]
			break
		}
	}

	if domain == nil {
		log.Printf("[WARN] Removing DNS domain %s because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	records, err := client.DNSRecord.List(context.Background(), d.Id())

	if err != nil {
		if strings.Contains(err.Error(), "Invalid domain") {
			log.Printf("[WARN] Removing DNS domain %s because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting DNS records for DNS domain %s: %v", d.Id(), err)
	}

	var record *govultr.DNSRecord
	for i := range records {
		if records[i].Type == "A" && records[i].Name == "" {
			record = &records[i]
			break
		}
	}

	if record == nil {
		log.Printf("[WARN] Removing DNS domain (%s) because it has no default record", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("domain", domain.Domain)
	d.Set("server_ip", record.Data)

	return nil
}

func resourceVultrDnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	records, err := client.DNSRecord.List(context.Background(), d.Id())

	if err != nil {
		if strings.Contains(err.Error(), "Invalid domain") {
			log.Printf("[WARN] Removing DNS domain %s because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting DNS records for DNS domain %s: %v", d.Id(), err)
	}

	var record *govultr.DNSRecord
	for i := range records {
		if records[i].Type == "A" && records[i].Name == "" {
			record = &records[i]
			break
		}
	}

	if record == nil {
		log.Printf("[WARN] Removing DNS domain (%s) because it has no default record", d.Id())
		d.SetId("")
		return resourceVultrDnsDomainRead(d, meta)
	}

	record.Data = d.Get("server_ip").(string)
	log.Print("[INFO] Updating DNS domain")
	err = client.DNSRecord.Update(context.Background(), d.Id(), record)
	if err != nil {
		return fmt.Errorf("Error updating the default DNS record for DNS domain %s : %v", d.Id(), err)
	}

	return resourceVultrDnsDomainRead(d, meta)
}

func resourceVultrDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting DNS domain (%s)", d.Id())
	err := client.DNSDomain.Delete(context.Background(), d.Id())

	if err != nil {
		return fmt.Errorf("Error destroying DNS domain %s: %v", d.Id(), err)
	}

	return nil
}
