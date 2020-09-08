package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrReverseIPV4() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrReverseIPV4Create,
		Read:   resourceVultrReverseIPV4Read,
		Delete: resourceVultrReverseIPV4Delete,

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

func resourceVultrReverseIPV4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)
	ip := d.Get("ip").(string)
	req := &govultr.ReverseIP{
		IP:      ip,
		Reverse: d.Get("reverse").(string),
	}

	log.Printf("[INFO] Creating reverse IPv4")

	if err := client.Instance.CreateReverseIPv4(context.Background(), instanceID, req); err != nil {
		return fmt.Errorf("error creating reverse IPv4: %v", err)
	}

	d.SetId(ip)

	return resourceVultrReverseIPV4Read(d, meta)
}

func resourceVultrReverseIPV4Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	var ReverseIPV4 *govultr.IPv4

	options := &govultr.ListOptions{}
	for {
		ReverseIPV4s, meta, err := client.Instance.ListIPv4(context.Background(), instanceID, options)
		if err != nil {
			return fmt.Errorf("error getting reverse IPv4s: %v, %v", err, instanceID)
		}

		for _, v := range ReverseIPV4s {

			if v.IP == d.Id() {
				ReverseIPV4 = &v
				break
			}
		}

		if ReverseIPV4 != nil {
			break
		}

		if meta.Links.Next == "" {
			return fmt.Errorf("error getting reverse IPv4s: %v, %v", err, instanceID)
		}

		options.Cursor = meta.Links.Next
	}

	d.Set("ip", ReverseIPV4.IP)
	d.Set("reverse", ReverseIPV4.Reverse)

	return nil
}

func resourceVultrReverseIPV4Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Deleting reverse IPv4: %s", d.Id())
	if err := client.Instance.DefaultReverseIPv4(context.Background(), instanceID, d.Id()); err != nil {
		return fmt.Errorf("error resetting reverse IPv4 (%s): %v", d.Id(), err)
	}

	return nil
}
