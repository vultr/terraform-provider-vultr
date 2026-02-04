package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrNATGatewayPortForwardingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrNATGatewayPortForwardingRuleCreate,
		ReadContext:   resourceVultrNATGatewayPortForwardingRuleRead,
		UpdateContext: resourceVultrNATGatewayPortForwardingRuleUpdate,
		DeleteContext: resourceVultrNATGatewayPortForwardingRuleDelete,
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
			"nat_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "both"}, false),
			},
			"external_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"internal_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"internal_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			// Optional
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrNATGatewayPortForwardingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)
	natGatewayID := d.Get("nat_gateway_id").(string)

	req := &govultr.NATGatewayPortForwardingRuleReq{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		InternalIP:   d.Get("internal_ip").(string),
		Protocol:     d.Get("protocol").(string),
		ExternalPort: d.Get("external_port").(int),
		InternalPort: d.Get("internal_port").(int),
		Enabled:      govultr.BoolToBoolPtr(d.Get("enabled").(bool)),
	}

	log.Printf("[INFO] Creating NAT Gateway port forwarding rule")
	NATGatewayPortForwardingRule, _, err := client.VPC.CreateNATGatewayPortForwardingRule(ctx, vpcID, natGatewayID, req)
	if err != nil {
		return diag.Errorf("error creating NAT Gateway port forwarding rule: %v", err)
	}

	d.SetId(NATGatewayPortForwardingRule.ID)

	return resourceVultrNATGatewayPortForwardingRuleRead(ctx, d, meta)
}

func resourceVultrNATGatewayPortForwardingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)
	natGatewayID := d.Get("nat_gateway_id").(string)

	natGatewayPortForwardingRule, _, err := client.VPC.GetNATGatewayPortForwardingRule(ctx, vpcID, natGatewayID, d.Id())
	if err != nil {
		return diag.Errorf("error getting NAT Gateway port forwarding rule (%s): %v", d.Id(), err)
	}

	if err := d.Set("name", natGatewayPortForwardingRule.Name); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `name` read value: %v", err)
	}

	if err := d.Set("description", natGatewayPortForwardingRule.Description); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `description` read value: %v", err)
	}

	if err := d.Set("internal_ip", natGatewayPortForwardingRule.InternalIP); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `internal_ip` read value: %v", err)
	}

	if err := d.Set("protocol", natGatewayPortForwardingRule.Protocol); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `protocol` read value: %v", err)
	}

	if err := d.Set("external_port", natGatewayPortForwardingRule.ExternalPort); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `external_port` read value: %v", err)
	}

	if err := d.Set("internal_port", natGatewayPortForwardingRule.InternalPort); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `internal_port` read value: %v", err)
	}

	if err := d.Set("enabled", natGatewayPortForwardingRule.Enabled); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `enabled` read value: %v", err)
	}

	if err := d.Set("date_created", natGatewayPortForwardingRule.DateCreated); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `date_created` read value: %v", err)
	}

	if err := d.Set("date_updated", natGatewayPortForwardingRule.DateUpdated); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway port forwarding rule `date_updated` read value: %v", err)
	}

	return nil
}

func resourceVultrNATGatewayPortForwardingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)
	natGatewayID := d.Get("nat_gateway_id").(string)

	req := &govultr.NATGatewayPortForwardingRuleReq{}

	if d.HasChange("name") {
		log.Printf("[INFO] Updating Name")
		_, newVal := d.GetChange("name")
		req.Name = newVal.(string)
	}

	if d.HasChange("description") {
		log.Printf("[INFO] Updating Description")
		_, newVal := d.GetChange("description")
		req.Description = newVal.(string)
	}

	if d.HasChange("internal_ip") {
		log.Printf("[INFO] Updating Internal IP")
		_, newVal := d.GetChange("internal_ip")
		req.InternalIP = newVal.(string)
	}

	if d.HasChange("protocol") {
		log.Printf("[INFO] Updating Protocol")
		_, newVal := d.GetChange("protocol")
		req.Protocol = newVal.(string)
	}

	if d.HasChange("external_port") {
		log.Printf("[INFO] Updating External Port")
		_, newVal := d.GetChange("external_port")
		req.ExternalPort = newVal.(int)
	}

	if d.HasChange("internal_port") {
		log.Printf("[INFO] Updating Internal Port")
		_, newVal := d.GetChange("internal_port")
		req.InternalPort = newVal.(int)
	}

	if d.HasChange("enabled") {
		log.Printf("[INFO] Updating Enabled")
		_, newVal := d.GetChange("enabled")
		req.Enabled = govultr.BoolToBoolPtr(newVal.(bool))
	}

	if _, _, err := client.VPC.UpdateNATGatewayPortForwardingRule(ctx, vpcID, natGatewayID, d.Id(), req); err != nil {
		return diag.Errorf("error updating NAT Gateway port forwarding rule %s : %s", d.Id(), err.Error())
	}

	return resourceVultrNATGatewayPortForwardingRuleRead(ctx, d, meta)
}

func resourceVultrNATGatewayPortForwardingRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting NAT Gateway port forwarding rule (%s)", d.Id())

	vpcID := d.Get("vpc_id").(string)
	natGatewayID := d.Get("nat_gateway_id").(string)

	if err := client.VPC.DeleteNATGatewayPortForwardingRule(ctx, vpcID, natGatewayID, d.Id()); err != nil {
		return diag.Errorf("error destroying NAT Gateway port forwarding rule %s : %v", d.Id(), err)
	}

	return nil
}
