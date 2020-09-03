package vultr

//
//import (
//	"context"
//	"fmt"
//	"log"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
//	"github.com/vultr/govultr/v2"
//)
//
//func resourceVultrReverseIPV4() *schema.Resource {
//	return &schema.Resource{
//		Create: resourceVultrReverseIPV4Create,
//		Read:   resourceVultrReverseIPV4Read,
//		Delete: resourceVultrReverseIPV4Delete,
//
//		Schema: map[string]*schema.Schema{
//			"instance_id": {
//				Type:     schema.TypeString,
//				Required: true,
//				ForceNew: true,
//			},
//			"ip": {
//				Type:     schema.TypeString,
//				Required: true,
//				ForceNew: true,
//			},
//			"reverse": {
//				Type:     schema.TypeString,
//				Required: true,
//				ForceNew: true,
//			},
//		},
//	}
//}
//
//func resourceVultrReverseIPV4Create(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	instanceID := d.Get("instance_id").(string)
//	ip := d.Get("ip").(string)
//	reverse := d.Get("reverse").(string)
//
//	log.Printf("[INFO] Creating reverse IPv4")
//
//	if err := client.Server.SetReverseIPV4(context.Background(), instanceID, ip, reverse); err != nil {
//		return fmt.Errorf("error creating reverse IPv4: %v", err)
//	}
//
//	d.SetId(ip)
//
//	return resourceVultrReverseIPV4Read(d, meta)
//}
//
//func resourceVultrReverseIPV4Read(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	instanceID := d.Get("instance_id").(string)
//
//	ReverseIPV4s, err := client.Server.IPV4Info(context.Background(), instanceID, true)
//	if err != nil {
//		return fmt.Errorf("error getting reverse IPv4s: %v, %v", err, instanceID)
//	}
//
//	var ReverseIPV4 *govultr.IPV4
//	for i := range ReverseIPV4s {
//		if ReverseIPV4s[i].IP == d.Id() {
//			ReverseIPV4 = &ReverseIPV4s[i]
//			break
//		}
//	}
//
//	if ReverseIPV4 == nil {
//		log.Printf("[WARN] Removing reverse IPv4 (%s) because it is gone", d.Id())
//		d.SetId("")
//		return nil
//	}
//
//	d.Set("ip", ReverseIPV4.IP)
//	d.Set("reverse", ReverseIPV4.Reverse)
//
//	return nil
//}
//
//func resourceVultrReverseIPV4Delete(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	instanceID := d.Get("instance_id").(string)
//
//	log.Printf("[INFO] Deleting reverse IPv4: %s", d.Id())
//	if err := client.Server.SetDefaultReverseIPV4(context.Background(), instanceID, d.Id()); err != nil {
//		return fmt.Errorf("error resetting reverse IPv4 (%s): %v", d.Id(), err)
//	}
//
//	return nil
//}
