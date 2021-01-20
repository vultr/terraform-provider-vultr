package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrInstanceIPV4() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrInstanceIPV4Create,
		Read:   resourceVultrInstanceIPV4Read,
		Delete: resourceVultrInstanceIPV4Delete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reboot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reverse": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrInstanceIPV4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Creating IPv4")

	ip, err := client.Instance.CreateIPv4(context.Background(), instanceID, govultr.BoolToBoolPtr(d.Get("reboot").(bool)))
	if err != nil {
		return fmt.Errorf("error creating IPv4: %v", err)
	}

	d.SetId(ip.IP)
	d.Set("instance_id", instanceID)

	return resourceVultrInstanceIPV4Read(d, meta)
}

func resourceVultrInstanceIPV4Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	var ipv4 *govultr.IPv4
	options := &govultr.ListOptions{}

	for {
		ips, meta, err := client.Instance.ListIPv4(context.Background(), instanceID, options)
		if err != nil {
			return fmt.Errorf("error getting IPv4s: %v", err)
		}

		for i := range ips {
			if ips[i].IP == d.Id() {
				ipv4 = &ips[i]
				break
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if ipv4 == nil {
		log.Printf("[WARN] Removing IPv4 (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("ip", ipv4.IP)
	d.Set("instance_id", instanceID)
	d.Set("reverse", ipv4.Reverse)
	d.Set("reboot", d.Get("reboot").(bool))

	return nil
}

func resourceVultrInstanceIPV4Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Deleting IPv4: %s", d.Id())
	if err := client.Instance.DeleteIPv4(context.Background(), instanceID, d.Id()); err != nil {
		return fmt.Errorf("error Deleting IPv4 (%s): %v", d.Id(), err)
	}

	return nil
}
