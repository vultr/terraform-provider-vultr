package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrVPC2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrVPC2Create,
		ReadContext:   resourceVultrVPC2Read,
		UpdateContext: resourceVultrVPC2Update,
		DeleteContext: resourceVultrVPC2Delete,
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
			"ip_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_block": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"prefix_length": {
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

func resourceVultrVPC2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcReq := &govultr.VPC2Req{
		Region:       d.Get("region").(string),
		Description:  d.Get("description").(string),
		IPType:       d.Get("ip_type").(string),
		IPBlock:      d.Get("ip_block").(string),
		PrefixLength: d.Get("prefix_length").(int),
	}

	vpc, _, err := client.VPC2.Create(ctx, vpcReq)
	if err != nil {
		return diag.Errorf("error creating VPC 2.0: %v", err)
	}

	d.SetId(vpc.ID)
	log.Printf("[INFO] VPC 2.0 ID: %s", d.Id())

	return resourceVultrVPC2Read(ctx, d, meta)
}

func resourceVultrVPC2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpc, _, err := client.VPC2.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid VPC 2.0 ID") {
			log.Printf("[WARN] Vultr VPC 2.0 (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting VPC 2.0: %v", err)
	}

	if err := d.Set("region", vpc.Region); err != nil {
		return diag.Errorf("unable to set resource vpc2 `region` read value: %v", err)
	}
	if err := d.Set("description", vpc.Description); err != nil {
		return diag.Errorf("unable to set resource vpc2 `description` read value: %v", err)
	}
	if err := d.Set("ip_block", vpc.IPBlock); err != nil {
		return diag.Errorf("unable to set resource vpc2 `ip_block` read value: %v", err)
	}
	if err := d.Set("prefix_length", vpc.PrefixLength); err != nil {
		return diag.Errorf("unable to set resource vpc2 `prefix_length` read value: %v", err)
	}
	if err := d.Set("date_created", vpc.DateCreated); err != nil {
		return diag.Errorf("unable to set resource vpc2 `date_created` read value: %v", err)
	}

	return nil
}

func resourceVultrVPC2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	if err := client.VPC2.Update(ctx, d.Id(), d.Get("description").(string)); err != nil {
		return diag.Errorf("error updating VPC 2.0: %v", err)
	}

	return resourceVultrVPC2Read(ctx, d, meta)
}

func resourceVultrVPC2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting VPC 2.0: %s", d.Id())

	retryErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete)-time.Minute, func() *retry.RetryError {
		err := client.VPC2.Delete(ctx, d.Id())

		if err == nil {
			return nil
		}

		if strings.Contains(err.Error(), "VPC 2.0 is attached") {
			return retry.RetryableError(fmt.Errorf("cannot remove attached VPC 2.0: %s", err.Error()))
		}

		return retry.NonRetryableError(err)
	})

	if retryErr != nil {
		return diag.Errorf("error destroying VPC 2.0 (%s): %v", d.Id(), retryErr)
	}

	return nil
}
