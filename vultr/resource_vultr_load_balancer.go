package vultr

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"balancing_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check": {
				Type:     schema.TypeMap,
				Optional: true,
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

	// Health
	healthCheck := &govultr.HealthCheck{}
	healthCheckData, healthCheckOk := d.GetOk("health_check")
	if healthCheckOk {
		healthCheck = generateHealthCheck(healthCheckData)
	} else {
		healthCheck = nil
	}

	fwMap := []govultr.ForwardingRule{}
	fr, frOk := d.GetOk("forwarding_rules")
	if frOk {
		for _, value := range fr.([]interface{}) {
			rule := generateRule(value.(map[string]interface{}))
			fwMap = append(fwMap, rule)
		}
	}

	if len(fwMap) == 0 {
		fwMap = nil
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
			idInt := value.(string)
			instance, err := client.Server.GetServer(context.Background(), idInt)
			if err != nil || instance.Status != "active" {
				return fmt.Errorf("Could not attach requested instance  %v to load balancer: %v", idInt, err)
			}
		}

		for _, value := range attachInstances.([]interface{}) {
			attachId, _ := strconv.Atoi(value.(string))
			instanceList.InstanceList = append(instanceList.InstanceList, attachId)
		}
	} else {
		instanceList = nil
	}

	// return fmt.Errorf("Error creating load balancer: REGION: %v LABEL: %v GENERIC: %v HEALTH: %v FR: %v SSL: %v INSTANCES: %v ", regionID, label, genericInfo, healthCheck, fwMap, ssl, instanceList)

	lb, err := client.LoadBalancer.Create(context.Background(), regionID, label, genericInfo, healthCheck, fwMap, ssl, instanceList)
	if err != nil {
		return fmt.Errorf("Error creating load balancer: %v", err)
	}
	id := strconv.Itoa(lb.ID)
	d.SetId(id)

	_, err = waitForLBAvailable(d, "active", []string{"pending", "installing"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error while waiting for load balancer %s to be completed: %s", d.Id(), err)
	}

	log.Printf("[INFO] load balancer ID: %s", d.Id())

	return nil
}

func resourceVultrLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	id, _ := strconv.Atoi(d.Id())
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

	// tk fix rix rule_id
	frList, err := client.LoadBalancer.ListForwardingRules(context.Background(), id)
	if err != nil {
		return fmt.Errorf("Error getting forwarding rules for load balancer (%v): %v", id, err)
	}

	var rulesList []map[string]interface{}
	for _, rules := range frList.ForwardRuleList {
		rule := map[string]interface{}{
			// "rule_id":            rules.RuleID,
			"frontend_protocol": rules.FrontendProtocol,
			"frontend_port":     strconv.Itoa(rules.FrontendPort),
			"backend_protocol":  rules.BackendProtocol,
			"backend_port":      strconv.Itoa(rules.BackendPort),
		}
		rulesList = append(rulesList, rule)
	}

	if err := d.Set("forwarding_rules", rulesList); err != nil {
		return fmt.Errorf("Error setting `forwarding_rules`: %v", err)
	}

	genericInfo, err := client.LoadBalancer.GetGenericInfo(context.Background(), id)
	if err != nil {
		return fmt.Errorf("Error getting generic info for loadbalancer (%v): %v", id, err)
	}

	instanceList, err := client.LoadBalancer.AttachedInstances(context.Background(), id)
	if err != nil {
		return fmt.Errorf("Error getting attached instance list for loadbalancer (%v): %v", id, err)
	}

	healthCheck, err := client.LoadBalancer.GetHealthCheck(context.Background(), id)
	if err != nil {
		return fmt.Errorf("Error getting health check info for loadbalancer (%v): %v", id, err)
	}

	hcInfo := map[string]interface{}{
		"protocol":            healthCheck.Protocol,
		"port":                strconv.Itoa(healthCheck.Port),
		"path":                healthCheck.Path,
		"check_interval":      strconv.Itoa(healthCheck.CheckInterval),
		"response_timeout":    strconv.Itoa(healthCheck.ResponseTimeout),
		"unhealthy_threshold": strconv.Itoa(healthCheck.UnhealthyThreshold),
		"healthy_threshold":   strconv.Itoa(healthCheck.HealthyThreshold),
	}

	if err := d.Set("health_check", hcInfo); err != nil {
		return fmt.Errorf("Error setting `health_check`: %#v", err)
	}

	ssl, err := client.LoadBalancer.HasSSL(context.Background(), id)
	if err != nil {
		return fmt.Errorf("Error getting ssl info for loadbalancer (%v): %v", id, err)
	}

	d.Set("has_ssl", ssl.SSLInfo)

	d.Set("attached_instances", instanceList.InstanceList)
	d.Set("balancing_algorithm", genericInfo.BalancingAlgorithm)
	d.Set("ssl_redirect", genericInfo.SSLRedirect)
	d.Set("proxy_protocol", genericInfo.ProxyProtocol)
	d.Set("cookie_name", genericInfo.StickySessions.CookieName)

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

	// Update GenericInfo
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

	// Health Check
	if d.HasChange("health_check") {
		_, newHealthCheck := d.GetChange("health_check")
		healthCheck := generateHealthCheck(newHealthCheck)

		log.Printf(`[INFO] Updating load balancer health info (%v)`, id)
		err := client.LoadBalancer.SetHealthCheck(context.Background(), id, healthCheck)
		if err != nil {
			return fmt.Errorf("Error updating load balancer health info (%v): %v", id, err)
		}
	}

	// SSL
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

		var oldIDs []interface{}
		for _, v := range oldFR.([]interface{}) {
			oldIDs = append(oldIDs, v.(interface{}))
		}

		var newIDs []interface{}
		for _, v := range newFR.([]interface{}) {
			newIDs = append(newIDs, v.(interface{}))
		}

		diff := func(in, out []interface{}) []interface{} {
			var diff []interface{}

			for _, v := range in {
				exists := false
				for _, v2 := range out {
					if reflect.DeepEqual(v, v2) {
						exists = true
					}
				}

				if !exists {
					diff = append(diff, v)
				}
			}

			return diff
		}

		// delete
		// formattedRules := make(map[string]interface{})
		// for _, v := range frList.ForwardRuleList {
		// 	newRule := make(map[string]interface{})
		// 	newRule["frontend_protocol"] = v.FrontendProtocol
		// 	newRule["frontend_port"] = v.FrontendPort
		// 	newRule["backend_protocol"] = v.BackendProtocol
		// 	newRule["backend_port"] = v.BackendPort
		// 	formattedRules[v.RuleID] = newRule
		// }

		// for _, value := range diff(oldIDs, newIDs) {
		// 	for _, val := range value.(map[string]interface{}) {

		// 		// if match {
		// 		// 	err := client.LoadBalancer.DeleteForwardingRule(context.Background(), id, curRuleID)

		// 		// 	if err != nil {
		// 		// 		return fmt.Errorf("Error updating forwarding rules for loadbalancer RULEID:%v --- %v: %v", curRuleID, id, err)
		// 		// 	}
		// 		// }

		// 	}

		// }

		// add
		for _, value := range diff(newIDs, oldIDs) {
			return fmt.Errorf("rules for loadbalancer %v: \n %v \n %v", value, oldFR, newFR)

			// rule := generateRule(value.(map[string]interface{}))
			// _, err := client.LoadBalancer.CreateForwardingRule(context.Background(), id, &rule)

			// if err != nil {
			// 	return fmt.Errorf("Error updating forwarding rules for loadbalancer %v: %v", value, err)
			// }
		}
	}

	if d.HasChange("attached_instances") {
		oldInstances, newInstances := d.GetChange("attached_instances")
		var oldIDs []string
		for _, v := range oldInstances.([]interface{}) {
			oldIDs = append(oldIDs, v.(string))
		}

		var newIDs []string
		for _, v := range newInstances.([]interface{}) {
			newIDs = append(newIDs, v.(string))
		}

		diff := func(in, out []string) []string {
			var diff []string

			b := map[string]string{}
			for i := range in {
				b[in[i]] = ""
			}

			for i := range out {
				if _, ok := b[out[i]]; !ok {
					diff = append(diff, out[i])
				}
			}

			return diff
		}

		for _, v := range diff(newIDs, oldIDs) {
			detachID, _ := strconv.Atoi(v)
			err := client.LoadBalancer.DetachInstance(context.Background(), id, detachID)

			if err != nil {
				return fmt.Errorf("Error detaching instance id %v from LoadBalancer %v : %v", v, id, diff(oldIDs, newIDs))
			}
		}

		for _, v := range diff(oldIDs, newIDs) {
			attachID, _ := strconv.Atoi(v)
			err := client.LoadBalancer.AttachInstance(context.Background(), id, attachID)

			if err != nil {
				return fmt.Errorf("Error attaching instance id %v to loadbalancer %v : %v", v, id, err)
			}
		}
	}

	return resourceVultrLoadBalancerRead(d, meta)
}

func resourceVultrLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting load balancer: %s", d.Id())

	id, _ := strconv.Atoi(d.Id())
	err := client.LoadBalancer.Delete(context.Background(), id)

	if err != nil {
		return fmt.Errorf("Error deleting load balancer %s : %v", d.Id(), err)
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
			return nil, "", fmt.Errorf("Error retrieving load balancer %s : %s", d.Id(), err)
		}

		if attr == "status" {
			log.Printf("[INFO] The load balancer Status is %s", lb.Status)
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
			num, _ := strconv.Atoi(v.(string))
			r.FrontendPort = num
		case "frontend_protocol":
			r.FrontendProtocol = v.(string)
		case "backend_port":
			num, _ := strconv.Atoi(v.(string))
			r.BackendPort = num
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
			num, _ := strconv.Atoi(v.(string))
			healthCheck.Port = num
		case "path":
			healthCheck.Path = v.(string)
		case "check_interval":
			num, _ := strconv.Atoi(v.(string))
			healthCheck.CheckInterval = num
		case "response_timeout":
			num, _ := strconv.Atoi(v.(string))
			healthCheck.ResponseTimeout = num
		case "unhealthy_threshold":
			num, _ := strconv.Atoi(v.(string))
			healthCheck.UnhealthyThreshold = num
		case "healthy_threshold":
			num, _ := strconv.Atoi(v.(string))
			healthCheck.HealthyThreshold = num
		}
	}

	return healthCheck
}
