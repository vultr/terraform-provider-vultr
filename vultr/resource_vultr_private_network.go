package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v2"
)

func resourceVultrPrivateNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrNetworkCreate,
		Read:   resourceVultrNetworkRead,
		Update: resourceVultrNetworkUpdate,
		Delete: resourceVultrNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"v4_subnet": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"v4_subnet_mask": {
				Type:     schema.TypeInt,
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

	networkReq := &govultr.NetworkReq{
		Region:       d.Get("region").(string),
		Description:  d.Get("description").(string),
		V4Subnet:     d.Get("v4_subnet").(string),
		V4SubnetMask: d.Get("v4_subnet_mask").(int),
	}

	network, err := client.Network.Create(context.Background(), networkReq)
	if err != nil {
		return fmt.Errorf("error creating network: %v", err)
	}

	d.SetId(network.NetworkID)
	log.Printf("[INFO] Network ID: %s", d.Id())

	return resourceVultrNetworkRead(d, meta)
}

func resourceVultrNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	network, err := client.Network.Get(context.Background(), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid private network ID") {
			log.Printf("[WARN] Vultr Private Network (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error getting private network: %v", err)
	}

	d.Set("region", network.Region)
	d.Set("description", network.Description)
	d.Set("v4_subnet", network.V4Subnet)
	d.Set("v4_subnet_mask", network.V4SubnetMask)
	d.Set("date_created", network.DateCreated)

	return nil
}

func resourceVultrNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	if err := client.Network.Update(context.Background(), d.Id(), d.Get("description").(string)); err != nil {
		return fmt.Errorf("error update private network: %v", err)
	}

	return resourceVultrNetworkRead(d, meta)
}

func resourceVultrNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Network: %s", d.Id())
	if err := client.Network.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error destroying private network (%s): %v", d.Id(), err)
	}

	return nil
}
