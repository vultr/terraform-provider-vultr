package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
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
	req := &govultr.ReverseIP{
		IP:      ip,
		Reverse: d.Get("reverse").(string),
	}

	if err := client.Instance.CreateReverseIPv6(context.Background(), instanceID, req); err != nil {
		return fmt.Errorf("error creating reverse IPv6: %v", err)
	}

	d.SetId(ip)

	return resourceVultrReverseIPV6Read(d, meta)
}

func resourceVultrReverseIPV6Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	reverseIPV6 := &govultr.ReverseIP{}

	reverseIPv6s, err := client.Instance.ListReverseIPv6(context.Background(), instanceID)
	if err != nil {
		return fmt.Errorf("error getting reverse IPv4s: %v, %v", err, instanceID)
	}

	for _, v := range reverseIPv6s {
		if v.IP == d.Id() {
			reverseIPV6 = &v
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
	if err := client.Instance.DeleteReverseIPv6(context.Background(), instanceID, d.Id()); err != nil {
		return fmt.Errorf("error destroying reverse IPv6 (%s): %v", d.Id(), err)
	}

	return nil
}
