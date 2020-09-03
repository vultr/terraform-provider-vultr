package vultr

//
//import (
//	"context"
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
//	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
//	"github.com/vultr/govultr/v2"
//	"log"
//	"net"
//	"strconv"
//)
//
//func resourceVultrDnsDomain() *schema.Resource {
//	return &schema.Resource{
//		Create: resourceVultrDnsDomainCreate,
//		Read:   resourceVultrDnsDomainRead,
//		Delete: resourceVultrDnsDomainDelete,
//		Importer: &schema.ResourceImporter{
//			State: schema.ImportStatePassthrough,
//		},
//		Schema: map[string]*schema.Schema{
//			"domain": {
//				Type:         schema.TypeString,
//				Required:     true,
//				ForceNew:     true,
//				ValidateFunc: validation.NoZeroValues,
//			},
//			"server_ip": {
//				Type:     schema.TypeString,
//				Optional: true,
//				ForceNew: true,
//			},
//		},
//	}
//}
//
//func resourceVultrDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	domain := d.Get("domain").(string)
//	ip, ipOk := d.GetOk("server_ip")
//
//	if ipOk {
//		validIp := net.ParseIP(ip.(string))
//		if validIp == nil {
//			return fmt.Errorf("the supplied IP address is invalid : %s", ip)
//		}
//	} else {
//		ip = "169.254.1.1"
//	}
//	log.Print("[INFO] Creating DNS domain")
//
//	err := client.DNSDomain.Create(context.Background(), domain, ip.(string))
//	if err != nil {
//		return fmt.Errorf("error while creating DNS domain : %s", err)
//	}
//
//	d.Set("server_ip", ip)
//	if ip == "169.254.1.1" {
//		d.Set("server_ip", nil)
//		log.Print("[INFO] DNS Domain : destroying default records")
//		records, err := client.DNSRecord.List(context.Background(), domain)
//		if err != nil {
//			return fmt.Errorf("error while retrieving records for delete : %s", err)
//		}
//
//		// clear out all records except NS
//		for _, v := range records {
//			if (v.Type == "A" && v.Name == "") || (v.Type == "CNAME" && v.Name == "*") || (v.Type == "MX" && v.Name == "") {
//				if err := client.DNSRecord.Delete(context.Background(), domain, strconv.Itoa(v.RecordID)); err != nil {
//					return fmt.Errorf("error while delete record %d : %s", v.RecordID, err)
//				}
//			}
//		}
//	}
//
//	d.SetId(domain)
//
//	return resourceVultrDnsDomainRead(d, meta)
//}
//
//func resourceVultrDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	domains, err := client.DNSDomain.List(context.Background())
//
//	if err != nil {
//		return fmt.Errorf("error getting domains : %v", err)
//	}
//
//	var domain *govultr.DNSDomain
//	for _, v := range domains {
//		if v.Domain == d.Id() {
//			domain = &v
//			break
//		}
//	}
//
//	if domain == nil {
//		log.Printf("[WARN] Removing DNS domain %s because it is gone", d.Id())
//		d.SetId("")
//		return nil
//	}
//
//	d.Set("domain", domain.Domain)
//
//	return nil
//}
//
//func resourceVultrDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	log.Printf("[INFO] Deleting DNS domain (%s)", d.Id())
//	err := client.DNSDomain.Delete(context.Background(), d.Id())
//
//	if err != nil {
//		return fmt.Errorf("error destroying DNS domain %s: %v", d.Id(), err)
//	}
//
//	return nil
//}
