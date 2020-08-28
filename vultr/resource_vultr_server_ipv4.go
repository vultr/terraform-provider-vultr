package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
)

func resourceVultrServerIPV4() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrServerIPV4Create,
		Read:   resourceVultrServerIPV4Read,
		Delete: resourceVultrServerIPV4Delete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reverse": {
				Type: schema.TypeString,
				Computed: true,
			},
			"reboot": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrServerIPV4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Creating IPv4")

	reboot := "yes"
	if d.Get("reboot").(bool) == false {
		reboot = "no"
	}

	ip, err := client.Server.AddIPV4(context.Background(), instanceID, reboot)
	if err != nil {
		return fmt.Errorf("error creating IPv4: %v", err)
	}

	d.SetId(ip.IPv4)
	d.Set("instance_id", instanceID)
	d.Set("reboot", reboot)

	return resourceVultrServerIPV4Read(d, meta)
}

func resourceVultrServerIPV4Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	ips, err := client.Server.IPV4Info(context.Background(), instanceID, true)
	if err != nil {
		return fmt.Errorf("error getting IPv4s: %v", err)
	}

	var ipv4s *govultr.IPV4
	for i := range ips {
		if ips[i].IP == d.Id() {
			ipv4s = &ips[i]
			break
		}
	}

	if ipv4s == nil {
		log.Printf("[WARN] Removing IPv4 (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("ip", ipv4s.IP)
	d.Set("instance_id", instanceID)
	d.Set("reverse", ipv4s.Reverse)
	d.Set("reboot", d.Get("reboot").(bool))

	return nil
}

func resourceVultrServerIPV4Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Deleting IPv4: %s", d.Id())
	err := client.Server.DestroyIPV4(context.Background(), instanceID, d.Id())

	if err != nil {
		return fmt.Errorf("error destroying IPv4 (%s): %v", d.Id(), err)
	}

	return nil
}
