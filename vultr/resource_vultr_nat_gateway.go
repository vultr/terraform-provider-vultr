package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrNATGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrNATGatewayCreate,
		ReadContext:   resourceVultrNATGatewayRead,
		UpdateContext: resourceVultrNATGatewayUpdate,
		DeleteContext: resourceVultrNATGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			// Optional
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_ips_v6": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"billing_charges": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"billing_monthly": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceVultrNATGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)

	req := &govultr.NATGatewayReq{
		Label: d.Get("label").(string),
		Tag:   d.Get("tag").(string),
	}

	log.Printf("[INFO] Creating NAT Gateway")
	natGateway, _, err := client.VPC.CreateNATGateway(ctx, vpcID, req)
	if err != nil {
		return diag.Errorf("error creating NAT Gateway: %v", err)
	}

	d.SetId(natGateway.ID)

	_, errWait := waitForNATGatewayAvailable(ctx, d, "active", "pending", "status", meta)
	if errWait != nil {
		return diag.Errorf("error while waiting for NAT Gateway %s to be in an active state : %s", d.Id(), err)
	}

	return resourceVultrNATGatewayRead(ctx, d, meta)
}

func resourceVultrNATGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)

	natGateway, _, err := client.VPC.GetNATGateway(ctx, vpcID, d.Id())
	if err != nil {
		return diag.Errorf("error getting NAT Gateway (%s): %v", d.Id(), err)
	}

	if err := d.Set("label", natGateway.Label); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `label` read value: %v", err)
	}

	if err := d.Set("tag", natGateway.Tag); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `tag` read value: %v", err)
	}

	if err := d.Set("date_created", natGateway.DateCreated); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `date_created` read value: %v", err)
	}

	if err := d.Set("status", natGateway.Status); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `status` read value: %v", err)
	}

	if err := d.Set("public_ips", natGateway.PublicIPs); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `public_ips` read value: %v", err)
	}

	if err := d.Set("public_ips_v6", natGateway.PublicIPsV6); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `public_ips_v6` read value: %v", err)
	}

	if err := d.Set("private_ips", natGateway.PrivateIPs); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `private_ips` read value: %v", err)
	}

	if err := d.Set("billing_charges", natGateway.Billing.Charges); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `billing_charges` read value: %v", err)
	}

	if err := d.Set("billing_monthly", natGateway.Billing.Monthly); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `billing_monthly` read value: %v", err)
	}

	return nil
}

func resourceVultrNATGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)

	if d.HasChange("label") {
		log.Printf("[INFO] Updating Label")
		_, newVal := d.GetChange("label")
		label := newVal.(string)
		req := &govultr.NATGatewayReq{
			Label: label,
		}
		if _, _, err := client.VPC.UpdateNATGateway(ctx, vpcID, d.Id(), req); err != nil {
			return diag.Errorf("error updating NAT Gateway %s : %s", d.Id(), err.Error())
		}
	}

	if d.HasChange("tag") {
		log.Printf("[INFO] Updating Tag")
		_, newVal := d.GetChange("tag")
		tag := newVal.(string)
		req := &govultr.NATGatewayReq{
			Tag: tag,
		}
		if _, _, err := client.VPC.UpdateNATGateway(ctx, vpcID, d.Id(), req); err != nil {
			return diag.Errorf("error updating NAT Gateway %s : %s", d.Id(), err.Error())
		}
	}

	return resourceVultrNATGatewayRead(ctx, d, meta)
}

func resourceVultrNATGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting NAT Gateway (%s)", d.Id())

	vpcID := d.Get("vpc_id").(string)

	if err := client.VPC.DeleteNATGateway(ctx, vpcID, d.Id()); err != nil {
		return diag.Errorf("error destroying NAT Gateway %s : %v", d.Id(), err)
	}

	return nil
}

func waitForNATGatewayAvailable(ctx context.Context, d *schema.ResourceData, target, pending string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll
	log.Printf(
		"[INFO] Waiting for NAT Gateway (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:        []string{pending},
		Target:         []string{target},
		Refresh:        newNATGatewayStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newNATGatewayStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) retry.StateRefreshFunc { //nolint:lll
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Creating NAT Gateway")
		natGateway, _, err := client.VPC.GetNATGateway(ctx, d.Get("vpc_id").(string), d.Id())

		if err != nil {
			return nil, "", fmt.Errorf("error retrieving NAT Gateway %s : %s", d.Id(), err)
		}

		if attr == "status" {
			log.Printf("[INFO] The NAT Gateway Status is %s", natGateway.Status)
			return natGateway, natGateway.Status, nil
		}

		return nil, "", nil
	}
}
