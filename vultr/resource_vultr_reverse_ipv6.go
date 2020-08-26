package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
)

func resourceVultrReverseIPV6() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrReverseIPV6Create,
		Read:   resourceVultrReverseIPV6Read,
		Delete: resourceVultrReverseIPV6Delete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reverse": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrReverseIPV6Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	ip := d.Get("ip").(string)

	reverse := d.Get("reverse").(string)
	log.Printf("[INFO] Creating reverse IPv6")

	err := client.Server.SetReverseIPV6(context.Background(), instanceID, ip, reverse)
	if err != nil {
		return fmt.Errorf("Error creating reverse IPv6: %v", err)
	}

	d.SetId(ip)

	return resourceVultrReverseIPV6Read(d, meta)
}

func resourceVultrReverseIPV6Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	reverseIPV6s, err := client.Server.ListReverseIPV6(context.Background(), instanceID)
	if err != nil {
		return fmt.Errorf("Error getting reverse IPv6s: %v", err)
	}

	var reverseIPV6 *govultr.ReverseIPV6
	for i := range reverseIPV6s {
		if reverseIPV6s[i].IP == d.Id() {
			reverseIPV6 = &reverseIPV6s[i]
			break
		}
	}

	if reverseIPV6 == nil {
		log.Printf("[WARN] Removing reverse IPv6 (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("ip", reverseIPV6.IP)
	d.Set("reverse", reverseIPV6.Reverse)

	return nil
}

func resourceVultrReverseIPV6Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Deleting reverse IPv6: %s", d.Id())
	err := client.Server.DeleteReverseIPV6(context.Background(), instanceID, d.Id())

	if err != nil {
		return fmt.Errorf("Error destroying reverse IPv6 (%s): %v", d.Id(), err)
	}

	return nil
}
