package vultr

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v2"
)

func resourceVultrVPC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrVPCCreate,
		ReadContext:   resourceVultrVPCRead,
		UpdateContext: resourceVultrVPCUpdate,
		DeleteContext: resourceVultrVPCDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceVultrVPCCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcReq := &govultr.VPCReq{
		Region:       d.Get("region").(string),
		Description:  d.Get("description").(string),
		V4Subnet:     d.Get("v4_subnet").(string),
		V4SubnetMask: d.Get("v4_subnet_mask").(int),
	}

	vpc, err := client.VPC.Create(ctx, vpcReq)
	if err != nil {
		return diag.Errorf("error creating VPC: %v", err)
	}

	d.SetId(vpc.ID)
	log.Printf("[INFO] VPC ID: %s", d.Id())

	return resourceVultrVPCRead(ctx, d, meta)
}

func resourceVultrVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpc, err := client.VPC.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid VPC ID") {
			log.Printf("[WARN] Vultr VPC (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting VPC: %v", err)
	}

	if err := d.Set("region", vpc.Region); err != nil {
		return diag.Errorf("unable to set resource vpc `region` read value: %v", err)
	}
	if err := d.Set("description", vpc.Description); err != nil {
		return diag.Errorf("unable to set resource vpc `description` read value: %v", err)
	}
	if err := d.Set("v4_subnet", vpc.V4Subnet); err != nil {
		return diag.Errorf("unable to set resource vpc `v4_subnet` read value: %v", err)
	}
	if err := d.Set("v4_subnet_mask", vpc.V4SubnetMask); err != nil {
		return diag.Errorf("unable to set resource vpc `v4_subnet_mask` read value: %v", err)
	}
	if err := d.Set("date_created", vpc.DateCreated); err != nil {
		return diag.Errorf("unable to set resource vpc `date_created` read value: %v", err)
	}

	return nil
}

func resourceVultrVPCUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	if err := client.VPC.Update(ctx, d.Id(), d.Get("description").(string)); err != nil {
		return diag.Errorf("error updating VPC: %v", err)
	}

	return resourceVultrVPCRead(ctx, d, meta)
}

func resourceVultrVPCDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting VPC: %s", d.Id())
	if err := client.VPC.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying VPC (%s): %v", d.Id(), err)
	}

	return nil
}
