package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
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
			"region_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forwarding_rules": {
				Type:     schema.TypeSet,
				Required: true,
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
			"balancing_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp"}, false),
						},
						"path": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
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
			"ssl": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"has_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl_redirect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"attached_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cookie_name": {
				Type:     schema.TypeString,
				Optional: true,
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
	regionID := d.Get("region_id").(int)
	label := d.Get("label").(string)
	_, proxyProtocolOk := d.GetOk("proxy_protocol")
	_, sslRedirectOk := d.GetOk("ssl_redirect")
	cookieName, cookieOk := d.GetOk("cookie_name")
	balancingAlgorithm, balancingAlgorithmOk := d.GetOk("balancing_algorithm")

	genericInfo := &govultr.GenericInfo{}
	stickySessions := &govultr.StickySessions{}
	if cookieOk {
		stickySessions.StickySessionsEnabled = "on"
		stickySessions.CookieName = cookieName.(string)
		genericInfo.StickySessions = stickySessions
	} else {
		genericInfo.StickySessions = nil
	}

	if balancingAlgorithmOk {
		genericInfo.BalancingAlgorithm = balancingAlgorithm.(string)
	}

	if proxyProtocolOk {
		proxyProtocol := d.Get("proxy_protocol").(bool)
		genericInfo.ProxyProtocol = &proxyProtocol
	}

	if sslRedirectOk {
		sslRedirect := d.Get("ssl_redirect").(bool)
		genericInfo.SSLRedirect = &sslRedirect
	}

	if !proxyProtocolOk && !balancingAlgorithmOk && !cookieOk && !sslRedirectOk {
		genericInfo = nil
	}

	healthCheck := &govultr.HealthCheck{}
	healthCheckData, healthCheckOk := d.GetOk("health_check")
	if healthCheckOk {
		for _, value := range healthCheckData.(*schema.Set).List() {
			healthCheck = generateHealthCheck(value)
			break
		}

		if healthCheck.Protocol == "tcp" && healthCheck.Path != "" {
			return fmt.Errorf("Error creating load balancer. Cannot set health check path when protocol is TCP")
		}
	} else {
		healthCheck = nil
	}

	fwMap := []govultr.ForwardingRule{}
	fr, frOk := d.GetOk("forwarding_rules")
	if frOk {
		for _, value := range fr.(*schema.Set).List() {
			rule := generateRule(value.(map[string]interface{}))
			fwMap = append(fwMap, rule)
		}
	}

	ssl := &govultr.SSL{}
	sslData, sslOk := d.GetOk("ssl")
	if sslOk {
		for k, v := range sslData.(map[string]interface{}) {
			switch k {
			case "private_key":
				ssl.PrivateKey = v.(string)
			case "certificate":
				ssl.Certificate = v.(string)
			case "chain":
				ssl.Chain = v.(string)
			}
		}
	} else {
		ssl = nil
	}

	instanceList := &govultr.InstanceList{}
	attachInstances, attachInstancesOk := d.GetOk("attached_instances")
	if attachInstancesOk {
		for _, value := range attachInstances.([]interface{}) {
			idInt := value.(int)
			instance, err := client.Server.GetServer(context.Background(), strconv.Itoa(idInt))
			if err != nil || instance.Status != "active" {
				return fmt.Errorf("Could not attach requested instance  %v to load balancer: %v", idInt, err)
			}
		}

		for _, value := range attachInstances.([]interface{}) {
			instanceList.InstanceList = append(instanceList.InstanceList, value.(int))
		}
	} else {
		instanceList = nil
	}

	lb, err := client.LoadBalancer.Create(context.Background(), regionID, label, genericInfo, healthCheck, fwMap, ssl, instanceList)
	if err != nil {
		return fmt.Errorf("Error creating load balancer: %v", err)
	}
	id := strconv.Itoa(lb.ID)
	d.SetId(id)

	_, err = waitForLBAvailable(d, "active", []string{"pending", "installing"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error while waiting for load balancer %v to be completed: %v", id, err)
	}

	log.Printf("[INFO] load balancer ID: %v", id)

	return resourceVultrLoadBalancerRead(d, meta)
}

func resourceVultrLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	id, _ := strconv.Atoi(d.Id())

	lbConfig, err := client.LoadBalancer.GetFullConfig(context.Background(), id)
	if err != nil {
		return fmt.Errorf("Error getting load balancer: %v", err)
	}

	lbs, err := client.LoadBalancer.List(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting load balancer: %v", err)
	}

	var lb *govultr.LoadBalancers
	for i := range lbs {
		if lbs[i].ID == id {
			lb = &lbs[i]
			break
		}
	}

	if lb == nil {
		log.Printf("[WARN] Vultr load balancer (%v) not found", id)
		d.SetId("")
		return nil
	}

	var rulesList []map[string]interface{}
	for _, rules := range lbConfig.ForwardRuleList {
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
		return fmt.Errorf("Error setting `forwarding_rules`: %v", err)
	}

	var hc []map[string]interface{}
	hcInfo := map[string]interface{}{
		"protocol":            lbConfig.HealthCheck.Protocol,
		"port":                lbConfig.HealthCheck.Port,
		"path":                lbConfig.HealthCheck.Path,
		"check_interval":      lbConfig.HealthCheck.CheckInterval,
		"response_timeout":    lbConfig.HealthCheck.ResponseTimeout,
		"unhealthy_threshold": lbConfig.HealthCheck.UnhealthyThreshold,
		"healthy_threshold":   lbConfig.HealthCheck.HealthyThreshold,
	}
	hc = append(hc, hcInfo)

	if err := d.Set("health_check", hc); err != nil {
		return fmt.Errorf("Error setting `health_check`: %v", err)
	}

	d.Set("has_ssl", lbConfig.SSLInfo)

	d.Set("attached_instances", lbConfig.InstanceList)
	d.Set("balancing_algorithm", lbConfig.BalancingAlgorithm)
	d.Set("ssl_redirect", lbConfig.SSLRedirect)
	d.Set("proxy_protocol", lbConfig.ProxyProtocol)
	d.Set("cookie_name", lbConfig.StickySessions.CookieName)

	d.Set("date_created", lb.DateCreated)
	d.Set("region_id", lb.RegionID)
	d.Set("location", lb.Location)
	d.Set("label", lb.Label)
	d.Set("status", lb.Status)
	d.Set("ipv4", lb.IPV4)
	d.Set("ipv6", lb.IPV6)

	return nil
}

func resourceVultrLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()
	id, _ := strconv.Atoi(d.Id())

	sslRedirect := d.Get("ssl_redirect").(bool)
	cookieName := d.Get("cookie_name").(string)

	stickySessions := govultr.StickySessions{
		StickySessionsEnabled: "on",
		CookieName:            cookieName,
	}

	if cookieName == "" {
		stickySessions.StickySessionsEnabled = "off"
	}

	genericInfo := govultr.GenericInfo{
		BalancingAlgorithm: d.Get("balancing_algorithm").(string),
		SSLRedirect:        &sslRedirect,
		StickySessions:     &stickySessions,
	}

	log.Printf(`[INFO] Updating load balancer generic info (%v)`, id)
	err := client.LoadBalancer.UpdateGenericInfo(context.Background(), id, d.Get("label").(string), &genericInfo)
	if err != nil {
		return fmt.Errorf("Error updating load balancer generic info (%v): %v", id, err)
	}

	if d.HasChange("health_check") {
		_, newHealthCheck := d.GetChange("health_check")
		healthCheck := &govultr.HealthCheck{}
		hcList := newHealthCheck.(*schema.Set).List()
		for _, value := range hcList {
			healthCheck = generateHealthCheck(value)
			break
		}

		if healthCheck.Protocol == "tcp" && healthCheck.Path != "" {
			return fmt.Errorf("Error updating load balancer (%v) health check. Cannot set health check path when protocol is TCP", id)
		}

		if len(hcList) == 0 {
			healthCheck = nil
		}

		log.Printf(`[INFO] Updating load balancer health info (%v)`, id)
		err := client.LoadBalancer.SetHealthCheck(context.Background(), id, healthCheck)
		if err != nil {
			return fmt.Errorf("Error updating load balancer health info (%v): %v", id, err)
		}
	}

	if d.HasChange("ssl") {
		ssl := govultr.SSL{}
		for k, v := range d.Get("ssl").(map[string]interface{}) {
			switch k {
			case "private_key":
				ssl.PrivateKey = v.(string)
			case "certificate":
				ssl.Certificate = v.(string)
			case "chain":
				ssl.Chain = v.(string)
			}
		}

		if ssl.PrivateKey == "" && ssl.Certificate == "" && ssl.Chain == "" {
			log.Printf(`[INFO] Removing load balancer SSL certificate (%v)`, id)
			err := client.LoadBalancer.RemoveSSL(context.Background(), id)
			if err != nil {
				return fmt.Errorf("Error removing SSL certificate for load balancer (%v): %v", id, err)
			}
		} else {
			log.Printf(`[INFO] Adding load balancer SSL certificate (%v)`, id)
			err := client.LoadBalancer.AddSSL(context.Background(), id, &ssl)
			if err != nil {
				return fmt.Errorf("Error adding SSL certificate for load balancer (%v): %v", id, err)
			}
		}
	}

	if d.HasChange("forwarding_rules") {
		oldFR, newFR := d.GetChange("forwarding_rules")

		oldFRList := oldFR.(*schema.Set).Difference(newFR.(*schema.Set))
		newFRList := newFR.(*schema.Set).Difference(oldFR.(*schema.Set))

		for _, value := range oldFRList.List() {
			for key, val := range value.(map[string]interface{}) {
				if key == "rule_id" {
					err := client.LoadBalancer.DeleteForwardingRule(context.Background(), id, val.(string))

					if err != nil {
						return fmt.Errorf("Error removing forwarding rules for load balancer (%v): %v", id, err)
					}
				}
			}
		}

		for _, value := range newFRList.List() {
			rule := generateRule(value.(map[string]interface{}))
			_, err := client.LoadBalancer.CreateForwardingRule(context.Background(), id, &rule)

			if err != nil {
				return fmt.Errorf("Error adding forwarding rules for load balancer (%v): %v", id, err)
			}
		}
	}

	if d.HasChange("attached_instances") {
		oldInstances, newInstances := d.GetChange("attached_instances")

		var oldIDs []int
		for _, v := range oldInstances.([]interface{}) {
			oldIDs = append(oldIDs, v.(int))
		}

		var newIDs []int
		for _, v := range newInstances.([]interface{}) {
			newIDs = append(newIDs, v.(int))
		}

		diff := func(in, out []int) []int {
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

		for _, v := range diff(newIDs, oldIDs) {
			err := client.LoadBalancer.DetachInstance(context.Background(), id, v)

			if err != nil {
				return fmt.Errorf("Error detaching instance id %v from load balancer (%v): %v", v, id, err)
			}
		}

		for _, v := range diff(oldIDs, newIDs) {
			err := client.LoadBalancer.AttachInstance(context.Background(), id, v)

			if err != nil {
				return fmt.Errorf("Error attaching instance id %v to load balancer (%v): %v", v, id, err)
			}
		}
	}

	return resourceVultrLoadBalancerRead(d, meta)
}

func resourceVultrLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting load balancer: %v", d.Id())

	id, _ := strconv.Atoi(d.Id())
	err := client.LoadBalancer.Delete(context.Background(), id)

	if err != nil {
		return fmt.Errorf("Error deleting load balancer %v : %v", d.Id(), err)
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

		id, _ := strconv.Atoi(d.Id())
		lbs, err := client.LoadBalancer.List(context.Background())
		if err != nil {
			return nil, "", fmt.Errorf("Error getting load balancer: %v", err)
		}

		var lb *govultr.LoadBalancers
		for i := range lbs {
			if lbs[i].ID == id {
				lb = &lbs[i]
				break
			}
		}

		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving load balancer %v : %v", d.Id(), err)
		}

		if attr == "status" {
			log.Printf("[INFO] The load balancer Status is %v", lb.Status)
			return lb, lb.Status, nil
		} else {
			return nil, "", nil
		}
	}
}

func generateRule(rule map[string]interface{}) govultr.ForwardingRule {
	r := govultr.ForwardingRule{}
	for k, v := range rule {
		switch k {
		case "frontend_port":
			r.FrontendPort = v.(int)
		case "frontend_protocol":
			r.FrontendProtocol = v.(string)
		case "backend_port":
			r.BackendPort = v.(int)
		case "backend_protocol":
			r.BackendProtocol = v.(string)
		}
	}
	return r
}

func generateHealthCheck(params interface{}) *govultr.HealthCheck {
	healthCheck := &govultr.HealthCheck{}
	for k, v := range params.(map[string]interface{}) {
		switch k {
		case "protocol":
			healthCheck.Protocol = v.(string)
		case "port":
			healthCheck.Port = v.(int)
		case "path":
			healthCheck.Path = v.(string)
		case "check_interval":
			healthCheck.CheckInterval = v.(int)
		case "response_timeout":
			healthCheck.ResponseTimeout = v.(int)
		case "unhealthy_threshold":
			healthCheck.UnhealthyThreshold = v.(int)
		case "healthy_threshold":
			healthCheck.HealthyThreshold = v.(int)
		}
	}

	return healthCheck
}
