package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrLoadBalancerCreate,
		Read:   resourceVultrLoadBalancerRead,
		Update: resourceVultrLoadBalancerUpdate,
		Delete: resourceVultrLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
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

			"has_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"attached_instances": {
				Type:     schema.TypeList,
				Optional: true,
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

func resourceVultrLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
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

	ssl := &govultr.SSL{}
	if sslData, sslOk := d.GetOk("ssl"); sslOk {
		ssl = generateSSL(sslData)
	} else {
		ssl = nil
	}

	cookieName, cookieOk := d.GetOk("cookie_name")
	stickySessions := &govultr.StickySessions{}
	if cookieOk {
		stickySessions.StickySessionsEnabled = "on"
		stickySessions.CookieName = cookieName.(string)
	} else {
		stickySessions = nil
	}

	req := &govultr.LoadBalancerReq{
		Region:             d.Get("region").(string),
		Label:              d.Get("label").(string),
		Instances:          instanceList,
		HealthCheck:        healthCheck,
		StickySessions:     stickySessions,
		ForwardingRules:    fwMap,
		SSL:                ssl,
		SSLRedirect:        d.Get("ssl_redirect").(bool),
		ProxyProtocol:      d.Get("proxy_protocol").(bool),
		BalancingAlgorithm: d.Get("balancing_algorithm").(string),
	}

	lb, err := client.LoadBalancer.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("error creating load balancer: %v", err)
	}
	d.SetId(lb.ID)

	_, err = waitForLBAvailable(d, "active", []string{"pending", "installing"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"error while waiting for load balancer %v to be completed: %v", lb.ID, err)
	}

	log.Printf("[INFO] load balancer ID: %v", lb.ID)

	return resourceVultrLoadBalancerRead(d, meta)
}

func resourceVultrLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	lb, err := client.LoadBalancer.Get(context.Background(), d.Id())
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
		return fmt.Errorf("error setting `forwarding_rules`: %v", err)
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
		return fmt.Errorf("error setting `health_check`: %v", err)
	}

	d.Set("has_ssl", lb.SSLInfo)

	d.Set("attached_instances", lb.Instances)
	d.Set("balancing_algorithm", lb.GenericInfo.BalancingAlgorithm)
	d.Set("proxy_protocol", lb.GenericInfo.ProxyProtocol)
	d.Set("cookie_name", lb.GenericInfo.StickySessions.CookieName)
	d.Set("date_created", lb.DateCreated)
	d.Set("label", lb.Label)
	d.Set("status", lb.Status)
	d.Set("ipv4", lb.IPV4)
	d.Set("ipv6", lb.IPV6)

	return nil
}

func resourceVultrLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	req := &govultr.LoadBalancerReq{
		Region:             d.Get("region").(string),
		Label:              d.Get("label").(string),
		SSLRedirect:        d.Get("ssl_redirect").(bool),
		ProxyProtocol:      d.Get("proxy_protocol").(bool),
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

		stickySessions.StickySessionsEnabled = "on"
		stickySessions.CookieName = cookieName.(string)
		req.StickySessions = stickySessions
	}

	if err := client.LoadBalancer.Update(context.Background(), d.Id(), req); err != nil {
		return fmt.Errorf("error updating load balancer generic info (%v): %v", d.Id(), err)
	}

	return resourceVultrLoadBalancerRead(d, meta)
}

func resourceVultrLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting load balancer: %v", d.Id())

	if err := client.LoadBalancer.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error deleting load balancer %v : %v", d.Id(), err)
	}

	return nil
}

func waitForLBAvailable(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for load balancer (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newLBStateRefresh(d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     5 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForState()
}

func newLBStateRefresh(d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating load balancer")

		lb, err := client.LoadBalancer.Get(context.Background(), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving lb %s ", d.Id())
		}

		if attr == "status" {
			log.Printf("[INFO] The load balancer Status is %v", lb.Status)
			return lb, lb.Status, nil
		} else {
			return nil, "", nil
		}
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

func diff(in, out []int) []int {
	var diff []int

	b := map[int]int{}
	for i := range in {
		b[in[i]] = 0
	}

	for i := range out {
		if _, ok := b[out[i]]; !ok {
			diff = append(diff, out[i])
		}
	}
	return diff
}
