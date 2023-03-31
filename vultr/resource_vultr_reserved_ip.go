package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrReservedIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrReservedIPCreate,
		ReadContext:   resourceVultrReservedIPRead,
		UpdateContext: resourceVultrReservedIPUpdate,
		DeleteContext: resourceVultrReservedIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: IgnoreCase,
			},
			"ip_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"v4", "v6"}, false),
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceVultrReservedIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.ReservedIPReq{
		Region:     d.Get("region").(string),
		IPType:     d.Get("ip_type").(string),
		Label:      d.Get("label").(string),
		InstanceID: d.Get("instance_id").(string),
	}
	rip, _, err := client.ReservedIP.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating reserved IP: %v", err)
	}

	d.SetId(rip.ID)
	log.Printf("[INFO] Reserved IP ID: %s", d.Id())

	if a, attachedOK := d.GetOk("instance_id"); attachedOK {
		if err := client.ReservedIP.Attach(ctx, d.Id(), a.(string)); err != nil {
			return diag.Errorf("error attaching reserved IP: %v %v : %v", d.Id(), a.(string), err)
		}
	}

	return resourceVultrReservedIPRead(ctx, d, meta)
}

func resourceVultrReservedIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	rip, _, err := client.ReservedIP.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid reserved-ip ID") {
			tflog.Warn(ctx, fmt.Sprintf("Removing reserved-ip (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting Reserved IPs: %v", err)
	}

	if rip == nil {
		log.Printf("[WARN] Vultr Reserved IP (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("region", rip.Region); err != nil {
		return diag.Errorf("unable to set resource reserved_ip `region` read value: %v", err)
	}
	if err := d.Set("ip_type", rip.IPType); err != nil {
		return diag.Errorf("unable to set resource reserved_ip `ip_type` read value: %v", err)
	}
	if err := d.Set("subnet", rip.Subnet); err != nil {
		return diag.Errorf("unable to set resource reserved_ip `subnet` read value: %v", err)
	}
	if err := d.Set("subnet_size", rip.SubnetSize); err != nil {
		return diag.Errorf("unable to set resource reserved_ip `subnet_size` read value: %v", err)
	}
	if err := d.Set("label", rip.Label); err != nil {
		return diag.Errorf("unable to set resource reserved_ip `label` read value: %v", err)
	}
	if err := d.Set("instance_id", rip.InstanceID); err != nil {
		return diag.Errorf("unable to set resource reserved_ip `instance_id` read value: %v", err)
	}

	return nil
}

func resourceVultrReservedIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	if d.HasChange("instance_id") {
		log.Printf("[INFO] Updating Reserved IP instance: %s", d.Id())

		old, newVal := d.GetChange("instance_id")

		if old.(string) != "" {
			if err := client.ReservedIP.Detach(ctx, d.Id()); err != nil {
				return diag.Errorf("error detaching Reserved IP (%s): %v", d.Id(), err)
			}
		}
		if newVal.(string) != "" {
			if err := client.ReservedIP.Attach(ctx, d.Id(), newVal.(string)); err != nil {
				return diag.Errorf("error attaching Reserved IP (%s): %v", d.Id(), err)
			}
		}
	}

	if d.HasChange("label") {
		log.Printf("[INFO] Updating Reserved IP label: %s", d.Id())

		req := &govultr.ReservedIPUpdateReq{
			Label: govultr.StringToStringPtr(d.Get("label").(string)),
		}

		if _, _, err := client.ReservedIP.Update(ctx, d.Id(), req); err != nil {
			return diag.Errorf("error updating reserved IP %s : %s", d.Id(), err.Error())
		}
	}

	return resourceVultrReservedIPRead(ctx, d, meta)
}

func resourceVultrReservedIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Reserved IP: %s", d.Id())
	if err := client.ReservedIP.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying Reserved IP (%s): %v", d.Id(), err)
	}

	return nil
}
