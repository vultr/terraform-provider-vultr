package vultr

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrPrivateNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrNetworkCreate,
		ReadContext:   resourceVultrNetworkRead,
		UpdateContext: resourceVultrNetworkUpdate,
		DeleteContext: resourceVultrNetworkDelete,
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
		DeprecationMessage: "Private Networks are deprecated and will not be supported in the future. Use VPCs instead.",
	}
}

func resourceVultrNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	networkReq := &govultr.NetworkReq{
		Region:       d.Get("region").(string),
		Description:  d.Get("description").(string),
		V4Subnet:     d.Get("v4_subnet").(string),
		V4SubnetMask: d.Get("v4_subnet_mask").(int),
	}

	network, _, err := client.Network.Create(ctx, networkReq)
	if err != nil {
		return diag.Errorf("error creating network: %v", err)
	}

	d.SetId(network.NetworkID)
	log.Printf("[INFO] Network ID: %s", d.Id())

	return resourceVultrNetworkRead(ctx, d, meta)
}

func resourceVultrNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	network, _, err := client.Network.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid private network ID") {
			log.Printf("[WARN] Vultr Private Network (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting private network: %v", err)
	}

	if err := d.Set("region", network.Region); err != nil {
		return diag.Errorf("unable to set resource private_network `region` read value: %v", err)
	}
	if err := d.Set("description", network.Description); err != nil {
		return diag.Errorf("unable to set resource private_network `description` read value: %v", err)
	}
	if err := d.Set("v4_subnet", network.V4Subnet); err != nil {
		return diag.Errorf("unable to set resource private_network `v4_subnet` read value: %v", err)
	}
	if err := d.Set("v4_subnet_mask", network.V4SubnetMask); err != nil {
		return diag.Errorf("unable to set resource private_network `v4_subnet_mask` read value: %v", err)
	}
	if err := d.Set("date_created", network.DateCreated); err != nil {
		return diag.Errorf("unable to set resource private_network `date_created` read value: %v", err)
	}

	return nil
}

func resourceVultrNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	if err := client.Network.Update(ctx, d.Id(), d.Get("description").(string)); err != nil {
		return diag.Errorf("error update private network: %v", err)
	}

	return resourceVultrNetworkRead(ctx, d, meta)
}

func resourceVultrNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Network: %s", d.Id())
	if err := client.Network.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying private network (%s): %v", d.Id(), err)
	}

	return nil
}
