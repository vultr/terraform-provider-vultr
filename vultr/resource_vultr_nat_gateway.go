package vultr

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
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
			"port_forwarding_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required
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
				},
			},
			"firewall_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
						},
						"subnet": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPv4Address,
						},
						"subnet_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 32),
						},
						"port": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validatePortOrPortRange,
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
				},
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

	var portForwardingRulesFlat []map[string]interface{}
	if portForwardingRules, portForwardingRulesOK := d.GetOk("port_forwarding_rules"); portForwardingRulesOK {
		portForwardingRulesList := portForwardingRules.(*schema.Set).List()
		for i := range portForwardingRulesList {
			portForwardingRulesFlat = append(portForwardingRulesFlat, portForwardingRulesList[i].(map[string]interface{}))
		}
	}

	var firewallRulesFlat []map[string]interface{}
	if firewallRules, firewallRulesOK := d.GetOk("firewall_rules"); firewallRulesOK {
		firewallRulesRulesList := firewallRules.(*schema.Set).List()
		for i := range firewallRulesRulesList {
			firewallRulesFlat = append(firewallRulesFlat, firewallRulesRulesList[i].(map[string]interface{}))
		}
	}

	rulesetErr := validateCompatibleRuleset(portForwardingRulesFlat, firewallRulesFlat)
	if rulesetErr != nil {
		return diag.Errorf("error creating NAT Gateway: %v", rulesetErr)
	}

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

	// Handle port forwarding rules after create
	if len(portForwardingRulesFlat) > 0 {
		log.Printf("[INFO] Adding NAT Gateway port forwarding rules")
		for i := range portForwardingRulesFlat {
			req2 := &govultr.NATGatewayPortForwardingRuleReq{
				Name:         portForwardingRulesFlat[i]["name"].(string),
				Description:  portForwardingRulesFlat[i]["description"].(string),
				InternalIP:   portForwardingRulesFlat[i]["internal_ip"].(string),
				Protocol:     portForwardingRulesFlat[i]["protocol"].(string),
				InternalPort: portForwardingRulesFlat[i]["internal_port"].(int),
				ExternalPort: portForwardingRulesFlat[i]["external_port"].(int),
				Enabled:      govultr.BoolToBoolPtr(portForwardingRulesFlat[i]["enabled"].(bool)),
			}

			_, _, pfwRuleErr := client.VPC.CreateNATGatewayPortForwardingRule(ctx, vpcID, d.Id(), req2)
			if pfwRuleErr != nil {
				return diag.Errorf("error creating NAT Gateway port forwarding rule: %v", pfwRuleErr)
			}
		}
	}

	// Handle firewall rules rules last
	if len(firewallRulesFlat) > 0 {
		log.Printf("[INFO] Adding NAT Gateway firewall rules")
		for i := range firewallRulesFlat {
			req3 := &govultr.NATGatewayFirewallRuleCreateReq{
				Protocol:   firewallRulesFlat[i]["protocol"].(string),
				Port:       firewallRulesFlat[i]["port"].(string),
				Subnet:     firewallRulesFlat[i]["subnet"].(string),
				SubnetSize: firewallRulesFlat[i]["subnet_size"].(int),
				Notes:      firewallRulesFlat[i]["notes"].(string),
			}

			_, _, fwRuleErr := client.VPC.CreateNATGatewayFirewallRule(ctx, vpcID, d.Id(), req3)
			if fwRuleErr != nil {
				return diag.Errorf("error creating NAT Gateway firewall rule: %v", fwRuleErr)
			}
		}
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

	portForwardingRules, _, _, err := client.VPC.ListNATGatewayPortForwardingRules(ctx, vpcID, d.Id(), nil)
	if err != nil {
		return diag.Errorf("error getting NAT Gateway port forwarding rules (%s): %v", d.Id(), err)
	}

	firewallRules, _, _, err := client.VPC.ListNATGatewayFirewallRules(ctx, vpcID, d.Id(), nil)
	if err != nil {
		return diag.Errorf("error getting NAT Gateway firewall rules (%s): %v", d.Id(), err)
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

	var portForwardingRulesFlat []map[string]interface{}
	for i := range portForwardingRules {
		portForwardingRulesFlat = append(portForwardingRulesFlat, map[string]interface{}{
			"id":            portForwardingRules[i].ID,
			"name":          portForwardingRules[i].Name,
			"protocol":      portForwardingRules[i].Protocol,
			"external_port": portForwardingRules[i].ExternalPort,
			"internal_ip":   portForwardingRules[i].InternalIP,
			"internal_port": portForwardingRules[i].InternalPort,
			"enabled":       portForwardingRules[i].Enabled,
			"description":   portForwardingRules[i].Description,
			"date_created":  portForwardingRules[i].DateCreated,
			"date_updated":  portForwardingRules[i].DateUpdated,
		})
	}

	if err := d.Set("port_forwarding_rules", portForwardingRulesFlat); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `port_forwarding_rules` read value: %v", err)
	}

	var firewallRulesFlat []map[string]interface{}
	for i := range firewallRules {
		firewallRulesFlat = append(firewallRulesFlat, map[string]interface{}{
			"id":          firewallRules[i].ID,
			"action":      firewallRules[i].Action,
			"protocol":    firewallRules[i].Protocol,
			"subnet":      firewallRules[i].Subnet,
			"subnet_size": firewallRules[i].SubnetSize,
			"port":        firewallRules[i].Port,
			"notes":       firewallRules[i].Notes,
		})
	}

	if err := d.Set("firewall_rules", firewallRulesFlat); err != nil {
		return diag.Errorf("unable to set resource NAT Gateway `firewall_rules` read value: %v", err)
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

func validateCompatibleRuleset(portForwardingRules, firewallRules []map[string]interface{}) (err error) {
	if len(portForwardingRules) == 0 {
		if len(firewallRules) > 0 {
			return fmt.Errorf("Firewall rules must have corresponding port forwarding rules for the same ports and protocols")
		}
		return nil
	}

	type Ruleset struct {
		tcp []int
		udp []int
	}

	var ruleset Ruleset
	for i := range portForwardingRules {
		protocol := portForwardingRules[i]["protocol"].(string)
		port := portForwardingRules[i]["external_port"].(int)
		enabled := portForwardingRules[i]["enabled"].(bool)

		if enabled && (protocol == "tcp" || protocol == "both") {
			ruleset.tcp = append(ruleset.tcp, port)
		}

		if enabled && (protocol == "tcp" || protocol == "both") {
			ruleset.udp = append(ruleset.udp, port)
		}
	}

	for i := range firewallRules {
		protocol := firewallRules[i]["protocol"].(string)
		parts := strings.Split(firewallRules[i]["port"].(string), ":")

		// Single port
		if len(parts) == 1 {
			fwPort, err2 := strconv.Atoi(parts[0])
			if err2 != nil {
				return fmt.Errorf("invalid port %v", fwPort)
			}

			if (protocol == "tcp" && !slices.Contains(ruleset.tcp, fwPort)) ||
				(protocol == "udp" && !slices.Contains(ruleset.udp, fwPort)) {
				return fmt.Errorf("Firewall rule for port %v and protocol %v requires a corresponding enabled forwarding rule", fwPort, protocol)
			}

			continue
		}

		// Port range
		startPort, err2 := strconv.Atoi(parts[0])
		if err2 != nil {
			return fmt.Errorf("invalid port range %v", firewallRules[i]["port"].(string))
		}

		endPort, err3 := strconv.Atoi(parts[1])
		if err3 != nil {
			return fmt.Errorf("invalid port range %v", firewallRules[i]["port"].(string))
		}

		for j := startPort; j <= endPort; j++ {
			if (protocol == "tcp" && !slices.Contains(ruleset.tcp, j)) ||
				(protocol == "udp" && !slices.Contains(ruleset.udp, j)) {
				return fmt.Errorf("Firewall rule for port %v and protocol %v requires a corresponding enabled forwarding rule", j, protocol)
			}
		}
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
