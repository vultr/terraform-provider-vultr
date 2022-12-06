package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrReverseIPV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrReverseIPV4Create,
		ReadContext:   resourceVultrReverseIPV4Read,
		DeleteContext: resourceVultrReverseIPV4Delete,

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
			"netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrReverseIPV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)
	ip := d.Get("ip").(string)
	req := &govultr.ReverseIP{
		IP:      ip,
		Reverse: d.Get("reverse").(string),
	}

	log.Printf("[INFO] Creating reverse IPv4")

	if err := client.Instance.CreateReverseIPv4(ctx, instanceID, req); err != nil {
		return diag.Errorf("error creating reverse IPv4: %v", err)
	}

	d.SetId(ip)

	if err := d.Set("instance_id", instanceID); err != nil {
		return diag.Errorf("unable to set resource reverse_ipv4 `instance_id` create value: %v", err)
	}

	return resourceVultrReverseIPV4Read(ctx, d, meta)
}

func resourceVultrReverseIPV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	var ReverseIPV4 *govultr.IPv4

	options := &govultr.ListOptions{}
	for {
		ReverseIPV4s, meta, err := client.Instance.ListIPv4(ctx, instanceID, options)
		if err != nil {
			return diag.Errorf("error getting reverse IPv4s: %v, %v", err, instanceID)
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
			return diag.Errorf("error getting reverse IPv4s: %v, %v", err, instanceID)
		}

		options.Cursor = meta.Links.Next
	}

	if err := d.Set("ip", ReverseIPV4.IP); err != nil {
		return diag.Errorf("unable to set resource reverse_ipv4 `ip` read value: %v", err)
	}
	if err := d.Set("reverse", ReverseIPV4.Reverse); err != nil {
		return diag.Errorf("unable to set resource reverse_ipv4 `reverse` read value: %v", err)
	}
	if err := d.Set("netmask", ReverseIPV4.Netmask); err != nil {
		return diag.Errorf("unable to set resource reverse_ipv4 `netmask` read value: %v", err)
	}
	if err := d.Set("gateway", ReverseIPV4.Gateway); err != nil {
		return diag.Errorf("unable to set resource reverse_ipv4 `gateway` read value: %v", err)
	}

	return nil
}

func resourceVultrReverseIPV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Deleting reverse IPv4: %s", d.Id())
	if err := client.Instance.DefaultReverseIPv4(ctx, instanceID, d.Id()); err != nil {
		return diag.Errorf("error resetting reverse IPv4 (%s): %v", d.Id(), err)
	}

	return nil
}
