package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrReverseIPV6() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrReverseIPV6Create,
		ReadContext:   resourceVultrReverseIPV6Read,
		DeleteContext: resourceVultrReverseIPV6Delete,

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

func resourceVultrReverseIPV6Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)
	ip := d.Get("ip").(string)
	req := &govultr.ReverseIP{
		IP:      ip,
		Reverse: d.Get("reverse").(string),
	}

	if err := client.Instance.CreateReverseIPv6(ctx, instanceID, req); err != nil {
		return diag.Errorf("error creating reverse IPv6: %v", err)
	}

	d.SetId(ip)

	return resourceVultrReverseIPV6Read(ctx, d, meta)
}

func resourceVultrReverseIPV6Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	reverseIPV6 := &govultr.ReverseIP{}

	reverseIPv6s, err := client.Instance.ListReverseIPv6(ctx, instanceID)
	if err != nil {
		return diag.Errorf("error getting reverse IPv4s: %v, %v", err, instanceID)
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

	if err := d.Set("ip", reverseIPV6.IP); err != nil {
		return diag.Errorf("unable to set resource reverse_ipv6 `ip` read value: %v", err)
	}
	if err := d.Set("reverse", reverseIPV6.Reverse); err != nil {
		return diag.Errorf("unable to set resource reverse_ipv6 `reverse` read value: %v", err)
	}

	return nil
}

func resourceVultrReverseIPV6Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instanceID := d.Get("instance_id").(string)

	log.Printf("[INFO] Deleting reverse IPv6: %s", d.Id())
	if err := client.Instance.DeleteReverseIPv6(ctx, instanceID, d.Id()); err != nil {
		return diag.Errorf("error destroying reverse IPv6 (%s): %v", d.Id(), err)
	}

	return nil
}
