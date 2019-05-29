package vultr

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func resourceVultrNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrNetworkCreate,
		Read:   resourceVultrNetworkRead,
		Delete: resourceVultrNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	regionID := d.Get("region_id").(string)

	var description string
	desc, descOk := d.GetOk("description")
	if descOk {
		description = desc.(string)
	}

	var cidrBlock string
	cidr, cidrOk := d.GetOk("cidr_block")
	if cidrOk {
		cidrBlock = cidr.(string)
		_, _, err := net.ParseCIDR(cidrBlock)
		if err != nil {
			return fmt.Errorf("Error parsing cidr_block for network: %v", err)
		}
	}

	network, err := client.Network.Create(context.Background(), regionID, description, cidrBlock)
	if err != nil {
		return fmt.Errorf("Error creating network: %v", err)
	}

	d.SetId(network.NetworkID)
	log.Printf("[INFO] Network ID: %s", d.Id())

	return resourceVultrNetworkRead(d, meta)
}

func resourceVultrNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	nets, err := client.Network.List(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting Reserved IPs: %v", err)
	}

	var network *govultr.Network
	for i := range nets {
		if nets[i].NetworkID == d.Id() {
			network = &nets[i]
			break
		}
	}

	if network == nil {
		log.Printf("[WARN] Vultr Reserved IP (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("region_id", network.RegionID)
	d.Set("description", network.Description)
	d.Set("date_created", network.DateCreated)
	d.Set("cidr_block", fmt.Sprintf("%s/%d", network.V4Subnet, network.V4SubnetMask))

	return nil
}

func resourceVultrNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Network: %s", d.Id())
	if err := client.Network.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("Error destroying Network (%s): %v", d.Id(), err)
	}

	return nil
}
