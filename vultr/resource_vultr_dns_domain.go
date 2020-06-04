package vultr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vultr/govultr"
	"log"
	"net"
	"strconv"
)

func resourceVultrDnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrDnsDomainCreate,
		Read:   resourceVultrDnsDomainRead,
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
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	domain := d.Get("domain").(string)
	ip, ipOk := d.GetOk("server_ip")

	if ipOk {
		validIp := net.ParseIP(ip.(string))
		if validIp == nil {
			return fmt.Errorf("The supplied IP address is invalid : %s", ip)
		}
	} else {
		ip = "169.254.1.1"
	}
	log.Print("[INFO] Creating DNS domain")

	err := client.DNSDomain.Create(context.Background(), domain, ip.(string))
	if err != nil {
		return fmt.Errorf("Error while creating DNS domain : %s", err)
	}

	d.Set("server_ip", ip)
	if ip == "169.254.1.1" {
		d.Set("server_ip", nil)
		log.Print("[INFO] DNS Domain : destroying default records")
		records, err := client.DNSRecord.List(context.Background(), domain)
		if err != nil {
			return fmt.Errorf("error while retrieving records for delete : %s", err)
		}

		// clear out all records expect NS
		for i := range records {
			if (records[i].Type == "A" && records[i].Name == "") || (records[i].Type == "CNAME" && records[i].Name == "*") || (records[i].Type == "MX" && records[i].Name == "") {
				if err := client.DNSRecord.Delete(context.Background(), domain, strconv.Itoa(records[i].RecordID)); err != nil {
					return fmt.Errorf("error while delete record %d : %s", records[i].RecordID, err)
				}
			}
		}
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

	d.Set("domain", domain.Domain)

	return nil
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
