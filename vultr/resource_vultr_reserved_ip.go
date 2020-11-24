package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vultr/govultr/v2"
)

func resourceVultrReservedIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrReservedIPCreate,
		Read:   resourceVultrReservedIPRead,
		Update: resourceVultrReservedIPUpdate,
		Delete: resourceVultrReservedIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ForceNew: true,
				Default:  "",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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

func resourceVultrReservedIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	req := &govultr.ReservedIPReq{
		Region:     d.Get("region").(string),
		IPType:     d.Get("ip_type").(string),
		Label:      d.Get("label").(string),
		InstanceID: d.Get("instance_id").(string),
	}
	rip, err := client.ReservedIP.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("error creating reserved IP: %v", err)
	}

	d.SetId(rip.ID)
	log.Printf("[INFO] Reserved IP ID: %s", d.Id())

	if a, attachedOK := d.GetOk("instance_id"); attachedOK {
		if err := client.ReservedIP.Attach(context.Background(), d.Id(), a.(string)); err != nil {
			return fmt.Errorf("error attaching reserved IP: %v %v", d.Id(), a.(string))
		}
	}

	return resourceVultrReservedIPRead(d, meta)
}

func resourceVultrReservedIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	rip, err := client.ReservedIP.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting Reserved IPs: %v", err)
	}

	if rip == nil {
		log.Printf("[WARN] Vultr Reserved IP (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("region", rip.Region)
	d.Set("ip_type", rip.IPType)
	d.Set("subnet", rip.Subnet)
	d.Set("subnet_size", rip.SubnetSize)
	d.Set("label", rip.Label)
	d.Set("instance_id", rip.InstanceID)

	return nil
}

func resourceVultrReservedIPUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("instance_id") {
		client := meta.(*Client).govultrClient()

		log.Printf("[INFO] Updating Reserved IP: %s", d.Id())

		old, newVal := d.GetChange("instance_id")

		if old.(string) != "" {
			if err := client.ReservedIP.Detach(context.Background(), d.Id()); err != nil {
				return fmt.Errorf("error detaching Reserved IP (%s): %v", d.Id(), err)
			}
		}
		if newVal.(string) != "" {
			if err := client.ReservedIP.Attach(context.Background(), d.Id(), newVal.(string)); err != nil {
				return fmt.Errorf("error attaching Reserved IP (%s): %v", d.Id(), err)
			}
		}
	}

	return resourceVultrReservedIPRead(d, meta)
}

func resourceVultrReservedIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Reserved IP: %s", d.Id())
	if err := client.ReservedIP.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error destroying Reserved IP (%s): %v", d.Id(), err)
	}

	return nil
}
