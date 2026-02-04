package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrNATGatewayFirewallRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrNATGatewayFirewallRuleCreate,
		ReadContext:   resourceVultrNATGatewayFirewallRuleRead,
		UpdateContext: resourceVultrNATGatewayFirewallRuleUpdate,
		DeleteContext: resourceVultrNATGatewayFirewallRuleDelete,
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
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
				ForceNew:     true,
			},
			"subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
				ForceNew:     true,
			},
			"subnet_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 32),
				ForceNew:     true,
			},
			"port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePortOrPortRange,
				ForceNew:     true,
			},
			// Optional
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrNATGatewayFirewallRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)
	natGatewayID := d.Get("nat_gateway_id").(string)

	req := &govultr.NATGatewayFirewallRuleCreateReq{
		Protocol:   d.Get("protocol").(string),
		Port:       d.Get("port").(string),
		Subnet:     d.Get("subnet").(string),
		SubnetSize: d.Get("subnet_size").(int),
		Notes:      d.Get("notes").(string),
	}

	log.Printf("[INFO] Creating NAT Gateway firewall rule")
	NATGatewayFirewallRule, _, err := client.VPC.CreateNATGatewayFirewallRule(ctx, vpcID, natGatewayID, req)
	if err != nil {
		return diag.Errorf("error creating NAT Gateway firewall rule: %v", err)
	}

	d.SetId(NATGatewayFirewallRule.ID)

	return resourceVultrNATGatewayFirewallRuleRead(ctx, d, meta)
}

func resourceVultrNATGatewayFirewallRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)
	natGatewayID := d.Get("nat_gateway_id").(string)

	natGatewayFirewallRule, _, err := client.VPC.GetNATGatewayFirewallRule(ctx, vpcID, natGatewayID, d.Id())
	if err != nil {
		return diag.Errorf("error getting NAT Gateway firewall rule (%s): %v", d.Id(), err)
	}

	if err := d.Set("action", natGatewayFirewallRule.Action); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway firewall rule `action` read value: %v", err)
	}

	if err := d.Set("protocol", natGatewayFirewallRule.Protocol); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway firewall rule `protocol` read value: %v", err)
	}

	if err := d.Set("port", natGatewayFirewallRule.Port); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway firewall rule `port` read value: %v", err)
	}

	if err := d.Set("subnet", natGatewayFirewallRule.Subnet); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway firewall rule `subnet` read value: %v", err)
	}

	if err := d.Set("subnet_size", natGatewayFirewallRule.SubnetSize); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway firewall rule `subnet_size` read value: %v", err)
	}

	if err := d.Set("notes", natGatewayFirewallRule.Notes); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway firewall rule `notes` read value: %v", err)
	}

	return nil
}

func resourceVultrNATGatewayFirewallRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpcID := d.Get("vpc_id").(string)
	natGatewayID := d.Get("nat_gateway_id").(string)

	if d.HasChange("notes") {
		log.Printf("[INFO] Updating Notes")
		_, newVal := d.GetChange("notes")
		req := &govultr.NATGatewayFirewallRuleUpdateReq{
			Notes: newVal.(string),
		}

		if _, _, err := client.VPC.UpdateNATGatewayFirewallRule(ctx, vpcID, natGatewayID, d.Id(), req); err != nil {
			return diag.Errorf("error updating NAT Gateway firewall rule %s : %s", d.Id(), err.Error())
		}
	}

	return resourceVultrNATGatewayFirewallRuleRead(ctx, d, meta)
}

func resourceVultrNATGatewayFirewallRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting NAT Gateway firewall rule (%s)", d.Id())

	vpcID := d.Get("vpc_id").(string)
	natGatewayID := d.Get("nat_gateway_id").(string)

	if err := client.VPC.DeleteNATGatewayFirewallRule(ctx, vpcID, natGatewayID, d.Id()); err != nil {
		return diag.Errorf("error destroying NAT Gateway firewall rule %s : %v", d.Id(), err)
	}

	return nil
}

func validatePortOrPortRange(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	parts := strings.Split(v, ":")
	if len(parts) > 2 {
		errs = append(errs, fmt.Errorf("%v must be a single port (1–65535) or port range (start:end), got %v", key, v))
		return
	}

	parsePort := func(p string) (int, error) {
		port, err := strconv.Atoi(p)
		if err != nil || port < 1 || port > 65535 {
			return 0, fmt.Errorf("invalid port %v", p)
		}
		return port, nil
	}

	// Single port
	if len(parts) == 1 {
		if _, err := parsePort(parts[0]); err != nil {
			errs = append(errs, fmt.Errorf("%v must be a valid port (1–65535), got %v", key, v))
		}
		return
	}

	// Port range
	start, err1 := parsePort(parts[0])
	end, err2 := parsePort(parts[1])

	if err1 != nil || err2 != nil || start > end {
		errs = append(errs, fmt.Errorf("%v must be a valid port range (start:end) with start ≤ end, got %v", key, v))
	}

	return
}
