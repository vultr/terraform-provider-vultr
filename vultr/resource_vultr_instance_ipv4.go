package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrInstanceIPV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrInstanceIPV4Create,
		ReadContext:   resourceVultrInstanceIPV4Read,
		DeleteContext: resourceVultrInstanceIPV4Delete,

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

func resourceVultrInstanceIPV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Creating IPv4")

	ip, _, err := client.Instance.CreateIPv4(ctx, instanceID, govultr.BoolToBoolPtr(d.Get("reboot").(bool)))
	if err != nil {
		return diag.Errorf("error creating IPv4: %v", err)
	}

	d.SetId(ip.IP)
	if err := d.Set("instance_id", instanceID); err != nil {
		return diag.Errorf("unable to set resource instance_ipv4 `instance_id` create value: %v", err)
	}

	return resourceVultrInstanceIPV4Read(ctx, d, meta)
}

func resourceVultrInstanceIPV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	var ipv4 *govultr.IPv4
	options := &govultr.ListOptions{}

	for {
		ips, meta, _, err := client.Instance.ListIPv4(ctx, instanceID, options)
		if err != nil {
			return diag.Errorf("error getting IPv4s: %v", err)
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

	if err := d.Set("ip", ipv4.IP); err != nil {
		return diag.Errorf("unable to set resource instance_ipv4 `ip` read value: %v", err)
	}
	if err := d.Set("instance_id", instanceID); err != nil {
		return diag.Errorf("unable to set resource instance_ipv4 `instance_id` read value: %v", err)
	}
	if err := d.Set("reverse", ipv4.Reverse); err != nil {
		return diag.Errorf("unable to set resource instance_ipv4 `reverse` read value: %v", err)
	}
	if err := d.Set("reboot", d.Get("reboot").(bool)); err != nil {
		return diag.Errorf("unable to set resource instance_ipv4 `reboot` read value: %v", err)
	}

	return nil
}

func resourceVultrInstanceIPV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Deleting IPv4: %s", d.Id())
	if err := client.Instance.DeleteIPv4(ctx, instanceID, d.Id()); err != nil {
		return diag.Errorf("error Deleting IPv4 (%s): %v", d.Id(), err)
	}

	return nil
}
