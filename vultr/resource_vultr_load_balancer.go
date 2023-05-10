package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

func resourceVultrLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrLoadBalancerCreate,
		ReadContext:   resourceVultrLoadBalancerRead,
		UpdateContext: resourceVultrLoadBalancerUpdate,
		DeleteContext: resourceVultrLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: IgnoreCase,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"balancing_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"leastconn", "roundrobin"}, false),
			},
			"private_network": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "private_network is deprecated and should no longer be used. Instead, use vpc",
			},
			"vpc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_redirect": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"proxy_protocol": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cookie_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp"}, false),
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"check_interval": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 300),
						},
						"response_timeout": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 300),
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 300),
						},
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 300),
						},
					},
				},
			},

			"forwarding_rules": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frontend_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp"}, false),
						},
						"frontend_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"backend_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp"}, false),
						},
						"backend_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ssl": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_key": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.NoZeroValues,
						},
						"certificate": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"chain": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"firewall_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"ip_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"v4", "v6"}, false),
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"has_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"attached_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	var healthCheck *govultr.HealthCheck
	if health, healthOk := d.GetOk("health_check"); healthOk {
		healthCheck = generateHealthCheck(health)
	} else {
		healthCheck = nil
	}

	var fwMap []govultr.ForwardingRule
	if fr, frOk := d.GetOk("forwarding_rules"); frOk {
		i := generateRules(fr)
		fwMap = i.ForwardRuleList
	} else {
		fwMap = nil
	}

	var instanceList []string
	if attachInstances, attachInstancesOk := d.GetOk("attached_instances"); attachInstancesOk {
		for _, id := range attachInstances.([]interface{}) {
			instanceList = append(instanceList, id.(string))
		}
	} else {
		instanceList = nil
	}

	var ssl *govultr.SSL
	if sslData, sslOk := d.GetOk("ssl"); sslOk {
		ssl = generateSSL(sslData)
	} else {
		ssl = nil
	}

	cookieName, cookieOk := d.GetOk("cookie_name")
	stickySessions := &govultr.StickySessions{}
	if cookieOk {
		stickySessions.CookieName = cookieName.(string)
	} else {
		stickySessions = nil
	}

	var fwrMap []govultr.LBFirewallRule
	if firewallRules, firewallRulesOk := d.GetOk("firewall_rules"); firewallRulesOk {
		fwrMap = generateFirewallRules(firewallRules)

	} else {
		fwrMap = nil
	}

	req := &govultr.LoadBalancerReq{
		Region:             d.Get("region").(string),
		Label:              d.Get("label").(string),
		Instances:          instanceList,
		HealthCheck:        healthCheck,
		StickySessions:     stickySessions,
		ForwardingRules:    fwMap,
		SSL:                ssl,
		SSLRedirect:        govultr.BoolToBoolPtr(d.Get("ssl_redirect").(bool)),
		ProxyProtocol:      govultr.BoolToBoolPtr(d.Get("proxy_protocol").(bool)),
		BalancingAlgorithm: d.Get("balancing_algorithm").(string),
		FirewallRules:      fwrMap,
	}

	if d.Get("private_network") != "" && d.Get("vpc") != "" {
		return diag.Errorf("private_network and vpc cannot be used together. Use only vpc instead.")
	}

	if d.Get("private_network") != "" {
		req.VPC = govultr.StringToStringPtr(d.Get("private_network").(string))
	}

	if d.Get("vpc") != "" {
		req.VPC = govultr.StringToStringPtr(d.Get("vpc").(string))
	}

	lb, _, err := client.LoadBalancer.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating load balancer: %v", err)
	}
	d.SetId(lb.ID)

	_, err = waitForLBAvailable(ctx, d, "active", []string{"pending", "installing"}, "status", meta)
	if err != nil {
		return diag.Errorf(
			"error while waiting for load balancer %v to be completed: %v", lb.ID, err)
	}

	log.Printf("[INFO] load balancer ID: %v", lb.ID)

	return resourceVultrLoadBalancerRead(ctx, d, meta)
}

func resourceVultrLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	lb, _, err := client.LoadBalancer.Get(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Vultr load balancer (%v) not found", d.Id())
		d.SetId("")
		return nil
	}

	var rulesList []map[string]interface{}
	for _, rules := range lb.ForwardingRules {
		rule := map[string]interface{}{
			"rule_id":           rules.RuleID,
			"frontend_protocol": rules.FrontendProtocol,
			"frontend_port":     rules.FrontendPort,
			"backend_protocol":  rules.BackendProtocol,
			"backend_port":      rules.BackendPort,
		}
		rulesList = append(rulesList, rule)
	}

	if err := d.Set("forwarding_rules", rulesList); err != nil {
		return diag.Errorf("unable to set resource load_balancer `forwarding_rules` read value: %v", err)
	}

	var fwrList []map[string]interface{}
	for _, rules := range lb.FirewallRules {
		rule := map[string]interface{}{
			"id":      rules.RuleID,
			"source":  rules.Source,
			"ip_type": rules.IPType,
			"port":    rules.Port,
		}
		fwrList = append(fwrList, rule)
	}

	if err := d.Set("firewall_rules", fwrList); err != nil {
		return diag.Errorf("unable to set resource load_balancer `firewall_rules` read value: %v", err)
	}

	var hc []map[string]interface{}
	hcInfo := map[string]interface{}{
		"protocol":            lb.HealthCheck.Protocol,
		"port":                lb.HealthCheck.Port,
		"path":                lb.HealthCheck.Path,
		"check_interval":      lb.HealthCheck.CheckInterval,
		"response_timeout":    lb.HealthCheck.ResponseTimeout,
		"unhealthy_threshold": lb.HealthCheck.UnhealthyThreshold,
		"healthy_threshold":   lb.HealthCheck.HealthyThreshold,
	}
	hc = append(hc, hcInfo)

	if err := d.Set("health_check", hc); err != nil {
		return diag.Errorf("unable to set resource load_balancer `health_check` read value: %v", err)
	}
	if err := d.Set("has_ssl", lb.SSLInfo); err != nil {
		return diag.Errorf("unable to set resource load_balancer `has_ssl` read value: %v", err)
	}
	if err := d.Set("attached_instances", lb.Instances); err != nil {
		return diag.Errorf("unable to set resource load_balancer `attached_instances` read value: %v", err)
	}
	if err := d.Set("balancing_algorithm", lb.GenericInfo.BalancingAlgorithm); err != nil {
		return diag.Errorf("unable to set resource load_balancer `balancing_algorithm` read value: %v", err)
	}
	if err := d.Set("proxy_protocol", lb.GenericInfo.ProxyProtocol); err != nil {
		return diag.Errorf("unable to set resource load_balancer `proxy_protocol` read value: %v", err)
	}
	if err := d.Set("cookie_name", lb.GenericInfo.StickySessions.CookieName); err != nil {
		return diag.Errorf("unable to set resource load_balancer `cookie_name` read value: %v", err)
	}
	if err := d.Set("label", lb.Label); err != nil {
		return diag.Errorf("unable to set resource load_balancer `label` read value: %v", err)
	}
	if err := d.Set("status", lb.Status); err != nil {
		return diag.Errorf("unable to set resource load_balancer `status` read value: %v", err)
	}
	if err := d.Set("ipv4", lb.IPV4); err != nil {
		return diag.Errorf("unable to set resource load_balancer `ipv4` read value: %v", err)
	}
	if err := d.Set("ipv6", lb.IPV6); err != nil {
		return diag.Errorf("unable to set resource load_balancer `ipv6` read value: %v", err)
	}
	if err := d.Set("region", lb.Region); err != nil {
		return diag.Errorf("unable to set resource load_balancer `region` read value: %v", err)
	}
	if err := d.Set("ssl_redirect", lb.GenericInfo.SSLRedirect); err != nil {
		return diag.Errorf("unable to set resource load_balancer `ssl_redirect` read value: %v", err)
	}

	// Manipulate the read state so that only one of these two values is
	// returned based on which is passed in. Needed since both private_network
	// and vpc are set to the same value after creation
	if d.Get("private_network") == "" && d.Get("vpc") != "" {
		if err := d.Set("private_network", ""); err != nil {
			return diag.Errorf("unable to set resource load_balancer `private_network` read value: %v", err)
		}
		if err := d.Set("vpc", lb.GenericInfo.VPC); err != nil {
			return diag.Errorf("unable to set resource load_balancer `vpc` read value: %v", err)
		}
	} else if d.Get("private_network") != "" && d.Get("vpc") == "" {
		if err := d.Set("private_network", lb.GenericInfo.VPC); err != nil {
			return diag.Errorf("unable to set resource load_balancer `private_network` read value: %v", err)
		}
		if err := d.Set("vpc", ""); err != nil {
			return diag.Errorf("unable to set resource load_balancer `vpc` read value: %v", err)
		}
	}

	return nil
}

func resourceVultrLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.LoadBalancerReq{
		Region:             d.Get("region").(string),
		Label:              d.Get("label").(string),
		SSLRedirect:        govultr.BoolToBoolPtr(d.Get("ssl_redirect").(bool)),
		ProxyProtocol:      govultr.BoolToBoolPtr(d.Get("proxy_protocol").(bool)),
		BalancingAlgorithm: d.Get("balancing_algorithm").(string),
	}

	if d.HasChange("health_check") {
		health := d.Get("health_check")
		req.HealthCheck = generateHealthCheck(health)
	}

	if d.HasChange("ssl") {
		if sslData, sslOK := d.GetOk("ssl"); sslOK {
			ssl := generateSSL(sslData)
			req.SSL = ssl
		} else {
			log.Printf(`[INFO] Removing load balancer SSL certificate (%v)`, d.Id())
			req.SSL = nil
		}
	}

	if d.HasChange("forwarding_rules") {
		_, newFR := d.GetChange("forwarding_rules")

		var rules []govultr.ForwardingRule
		for _, val := range newFR.(*schema.Set).List() {
			t := val.(map[string]interface{})

			rule := govultr.ForwardingRule{
				FrontendProtocol: t["frontend_protocol"].(string),
				FrontendPort:     t["frontend_port"].(int),
				BackendProtocol:  t["backend_protocol"].(string),
				BackendPort:      t["backend_port"].(int),
			}
			rules = append(rules, rule)

		}
		req.ForwardingRules = rules
	}

	if d.HasChange("firewall_rules") {
		_, newFWR := d.GetChange("firewall_rules")

		fwList := newFWR.(*schema.Set).List()
		req.FirewallRules = []govultr.LBFirewallRule{}

		if len(fwList) != 0 {
			for _, val := range newFWR.(*schema.Set).List() {
				t := val.(map[string]interface{})
				rule := govultr.LBFirewallRule{
					Port:   t["port"].(int),
					Source: t["source"].(string),
					IPType: t["ip_type"].(string),
				}
				req.FirewallRules = append(req.FirewallRules, rule)
			}
		}
	}

	if d.HasChange("attached_instances") {
		_, newInstances := d.GetChange("attached_instances")
		log.Println(newInstances)

		var newIDs []string
		for _, v := range newInstances.([]interface{}) {
			newIDs = append(newIDs, v.(string))
		}

		req.Instances = newIDs
	}

	if d.HasChange("cookie_name") {
		stickySessions := &govultr.StickySessions{}
		cookieName := d.Get("cookie_name")

		stickySessions.CookieName = cookieName.(string)
		req.StickySessions = stickySessions
	}

	if d.Get("private_network") != "" && d.Get("vpc") != "" {
		return diag.Errorf("private_network and vpc cannot be used together. Use only vpc instead.")
	}

	if d.HasChange("private_network") {
		req.VPC = govultr.StringToStringPtr(d.Get("private_network").(string))
	}

	if d.HasChange("vpc") {
		req.VPC = govultr.StringToStringPtr(d.Get("vpc").(string))
	}

	if err := client.LoadBalancer.Update(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating load balancer generic info (%v): %v", d.Id(), err)
	}

	return resourceVultrLoadBalancerRead(ctx, d, meta)
}

func resourceVultrLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting load balancer: %v", d.Id())

	// items we should detach before we destroy
	// instances and firewall rules are default "null" if not present in LoadBalancerReq
	detachConfig := &govultr.LoadBalancerReq{}

	if _, vpcOK := d.GetOk("vpc"); vpcOK {
		detachConfig.VPC = govultr.StringToStringPtr("")
	}

	// send update to perform detach
	if err := client.LoadBalancer.Update(ctx, d.Id(), detachConfig); err != nil {
		return diag.Errorf("error detaching VPC from load balancer before deletion (%v): %v", d.Id(), err)
	}

	if err := client.LoadBalancer.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting load balancer %v : %v", d.Id(), err)
	}

	return nil
}

func waitForLBAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for load balancer (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{ // nolint:all
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newLBStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     5 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newLBStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc { // nolint:all
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating load balancer")

		lb, _, err := client.LoadBalancer.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving lb %s ", d.Id())
		}

		if attr == "status" {
			log.Printf("[INFO] The load balancer Status is %v", lb.Status)
			return lb, lb.Status, nil
		}

		return nil, "", nil
	}
}

func generateRules(rules interface{}) *govultr.ForwardingRules {
	fwMap := &govultr.ForwardingRules{}
	for _, rule := range rules.(*schema.Set).List() {
		r := rule.(map[string]interface{})
		t := govultr.ForwardingRule{
			FrontendProtocol: r["frontend_protocol"].(string),
			FrontendPort:     r["frontend_port"].(int),
			BackendProtocol:  r["backend_protocol"].(string),
			BackendPort:      r["backend_port"].(int),
		}
		fwMap.ForwardRuleList = append(fwMap.ForwardRuleList, t)
	}
	return fwMap
}

func generateFirewallRules(rules interface{}) []govultr.LBFirewallRule {
	var fwrMap []govultr.LBFirewallRule
	for _, rule := range rules.(*schema.Set).List() {
		r := rule.(map[string]interface{})
		t := govultr.LBFirewallRule{
			Port:   r["port"].(int),
			Source: r["source"].(string),
			IPType: r["ip_type"].(string),
		}
		fwrMap = append(fwrMap, t)
	}
	return fwrMap
}

func generateHealthCheck(health interface{}) *govultr.HealthCheck {
	k := health.([]interface{})
	config := k[0].(map[string]interface{})

	return &govultr.HealthCheck{
		Protocol:           config["protocol"].(string),
		Port:               config["port"].(int),
		Path:               config["path"].(string),
		CheckInterval:      config["check_interval"].(int),
		ResponseTimeout:    config["response_timeout"].(int),
		UnhealthyThreshold: config["unhealthy_threshold"].(int),
		HealthyThreshold:   config["healthy_threshold"].(int),
	}
}

func generateSSL(sslData interface{}) *govultr.SSL {
	k := sslData.(*schema.Set).List()
	config := k[0].(map[string]interface{})

	return &govultr.SSL{
		PrivateKey:  config["private_key"].(string),
		Certificate: config["certificate"].(string),
		Chain:       config["chain"].(string),
	}
}
