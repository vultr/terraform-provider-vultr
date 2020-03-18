package vultr

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
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
			"region_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateVultrReservedIPType,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"attached_id": {
				Type:     schema.TypeString,
				Optional: true,
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

	regionID := d.Get("region_id").(int)
	ipType := d.Get("ip_type").(string)

	var label string
	l, labelOk := d.GetOk("label")
	if labelOk {
		label = l.(string)
	}

	rip, err := client.ReservedIP.Create(context.Background(), regionID, ipType, label)
	if err != nil {
		return fmt.Errorf("Error creating reserved IP: %v", err)
	}

	d.SetId(rip.ReservedIPID)
	log.Printf("[INFO] Reserved IP ID: %s", d.Id())

	var attachedTo string
	a, attachedOK := d.GetOk("attached_id")
	if attachedOK {
		err := resourceVultrReservedIPRead(d, meta)

		if err != nil {
			return fmt.Errorf("Error occured while creating reservedIP : %v", err)
		}
		attachedTo = a.(string)
		err = client.ReservedIP.Attach(context.Background(), d.Get("subnet").(string), attachedTo)
		if err != nil {
			return fmt.Errorf("Error attaching reserved IP: %v", err)
		}
	}

	return resourceVultrReservedIPRead(d, meta)
}

func resourceVultrReservedIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	rips, err := client.ReservedIP.List(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting Reserved IPs: %v", err)
	}

	var rip *govultr.ReservedIP
	for i := range rips {
		if rips[i].ReservedIPID == d.Id() {
			rip = &rips[i]
			break
		}
	}

	if rip == nil {
		log.Printf("[WARN] Vultr Reserved IP (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("region_id", rip.RegionID)
	d.Set("ip_type", rip.IPType)
	d.Set("subnet", rip.Subnet)
	d.Set("subnet_size", rip.SubnetSize)
	d.Set("label", rip.Label)
	d.Set("attached_id", rip.AttachedID)

	return nil
}

func resourceVultrReservedIPUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("attached_id") {
		client := meta.(*Client).govultrClient()

		log.Printf("[INFO] Updating Reserved IP: %s", d.Id())

		ip := d.Get("subnet").(string)

		old, newVal := d.GetChange("attached_id")
		if old.(string) != "" {
			err := client.ReservedIP.Detach(context.Background(), ip, old.(string))
			if err != nil {
				return fmt.Errorf("Error detaching Reserved IP (%s): %v", d.Id(), err)
			}
		}
		if newVal.(string) != "" {
			err := client.ReservedIP.Attach(context.Background(), ip, newVal.(string))
			if err != nil {
				return fmt.Errorf("Error attaching Reserved IP (%s): %v", d.Id(), err)
			}
		}
	}

	return resourceVultrReservedIPRead(d, meta)
}

func resourceVultrReservedIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	if d.Get("attached_id").(string) != "" {
		err := client.ReservedIP.Detach(context.Background(), d.Get("subnet").(string), d.Get("attached_id").(string))
		if err != nil {
			return fmt.Errorf("error detaching rerservedIP (%s): %v", d.Id(), err)
		}
	}

	log.Printf("[INFO] Deleting Reserved IP: %s", d.Id())
	if err := client.ReservedIP.Delete(context.Background(), d.Get("subnet").(string)); err != nil {
		return fmt.Errorf("Error destroying Reserved IP (%s): %v", d.Id(), err)
	}

	return nil
}

func validateVultrReservedIPType(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if v != "v6" && v != "v4" {
		errs = append(errs, fmt.Errorf("%q must be either 'v4' or 'v6', got: %s", key, v))
	}
	return
}
