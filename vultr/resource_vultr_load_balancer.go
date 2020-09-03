package vultr

//
//import (
//	"context"
//	"fmt"
//	"log"
//	"strconv"
//	"time"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
//	"github.com/vultr/govultr/v2"
//)
//
//func resourceVultrLoadBalancer() *schema.Resource {
//	return &schema.Resource{
//		Create: resourceVultrLoadBalancerCreate,
//		Read:   resourceVultrLoadBalancerRead,
//		Update: resourceVultrLoadBalancerUpdate,
//		Delete: resourceVultrLoadBalancerDelete,
//		Importer: &schema.ResourceImporter{
//			State: schema.ImportStatePassthrough,
//		},
//		Schema: map[string]*schema.Schema{
//			"region_id": {
//				Type:     schema.TypeInt,
//				ForceNew: true,
//				Required: true,
//			},
//			"label": {
//				Type:     schema.TypeString,
//				Optional: true,
//			},
//			"balancing_algorithm": {
//				Type:         schema.TypeString,
//				Optional:     true,
//				Computed:     true,
//				ValidateFunc: validation.StringInSlice([]string{"leastconn", "roundrobin"}, false),
//			},
//			"ssl_redirect": {
//				Type:     schema.TypeBool,
//				Optional: true,
//				Computed: true,
//			},
//			"proxy_protocol": {
//				Type:     schema.TypeBool,
//				Optional: true,
//				Computed: true,
//			},
//			"cookie_name": {
//				Type:     schema.TypeString,
//				Optional: true,
//			},
//			"health_check": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Optional: true,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"protocol": {
//							Type:         schema.TypeString,
//							Required:     true,
//							ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp"}, false),
//						},
//						"path": {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						"port": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(1, 65535),
//						},
//						"check_interval": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(1, 300),
//						},
//						"response_timeout": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(1, 300),
//						},
//						"unhealthy_threshold": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(1, 300),
//						},
//						"healthy_threshold": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(1, 300),
//						},
//					},
//				},
//			},
//
//			"forwarding_rules": {
//				Type:     schema.TypeSet,
//				Required: true,
//				MinItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"frontend_protocol": {
//							Type:         schema.TypeString,
//							Required:     true,
//							ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp"}, false),
//						},
//						"frontend_port": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(1, 65535),
//						},
//						"backend_protocol": {
//							Type:         schema.TypeString,
//							Required:     true,
//							ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp"}, false),
//						},
//						"backend_port": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(1, 65535),
//						},
//						"rule_id": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//					},
//				},
//			},
//
//			"ssl": {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"private_key": {
//							Type:         schema.TypeString,
//							Required:     true,
//							Sensitive:    true,
//							ValidateFunc: validation.NoZeroValues,
//						},
//						"certificate": {
//							Type:         schema.TypeString,
//							Required:     true,
//							ValidateFunc: validation.NoZeroValues,
//						},
//						"chain": {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//					},
//				},
//			},
//
//			"has_ssl": {
//				Type:     schema.TypeBool,
//				Computed: true,
//			},
//			"attached_instances": {
//				Type:     schema.TypeList,
//				Optional: true,
//				Elem:     &schema.Schema{Type: schema.TypeInt},
//			},
//			"status": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"ipv4": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"ipv6": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//		},
//	}
//}
//
//func resourceVultrLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//	regionID := d.Get("region_id").(int)
//	label := d.Get("label").(string)
//
//	proxy, proxyProtocolOk := d.GetOk("proxy_protocol")
//	sslRedirect, sslRedirectOk := d.GetOk("ssl_redirect")
//	cookieName, cookieOk := d.GetOk("cookie_name")
//	balancingAlgorithm, balancingAlgorithmOk := d.GetOk("balancing_algorithm")
//
//	genericInfo := &govultr.GenericInfo{}
//	if proxyProtocolOk || balancingAlgorithmOk || cookieOk || sslRedirectOk {
//		if proxyProtocolOk {
//			proxyProtocol := proxy.(bool)
//			genericInfo.ProxyProtocol = &proxyProtocol
//		}
//
//		if balancingAlgorithmOk {
//			genericInfo.BalancingAlgorithm = balancingAlgorithm.(string)
//		}
//
//		if sslRedirectOk {
//			sslRedirect := sslRedirect.(bool)
//			genericInfo.SSLRedirect = &sslRedirect
//		}
//
//		stickySessions := &govultr.StickySessions{}
//		if cookieOk {
//			stickySessions.StickySessionsEnabled = "on"
//			stickySessions.CookieName = cookieName.(string)
//			genericInfo.StickySessions = stickySessions
//		} else {
//			genericInfo.StickySessions = nil
//		}
//	} else {
//		genericInfo = nil
//	}
//
//	var healthCheck *govultr.HealthCheck
//	if health, healthOk := d.GetOk("health_check"); healthOk {
//		healthCheck = generateHealthCheck(health)
//	} else {
//		healthCheck = nil
//	}
//
//	fwMap := []govultr.ForwardingRule{}
//	if fr, frOk := d.GetOk("forwarding_rules"); frOk {
//		i := generateRules(fr)
//		fwMap = i.ForwardRuleList
//	} else {
//		fwMap = nil
//	}
//
//	ssl := &govultr.SSL{}
//	if sslData, sslOk := d.GetOk("ssl"); sslOk {
//		ssl = generateSSL(sslData)
//	} else {
//		ssl = nil
//	}
//
//	instanceList := &govultr.InstanceList{}
//	if attachInstances, attachInstancesOk := d.GetOk("attached_instances"); attachInstancesOk {
//		for _, id := range attachInstances.([]interface{}) {
//			instanceList.InstanceList = append(instanceList.InstanceList, id.(int))
//		}
//	} else {
//		instanceList = nil
//	}
//
//	lb, err := client.LoadBalancer.Create(context.Background(), regionID, label, genericInfo, healthCheck, fwMap, ssl, instanceList)
//	if err != nil {
//		return fmt.Errorf("Error creating load balancer: %v", err)
//	}
//	id := strconv.Itoa(lb.ID)
//	d.SetId(id)
//	d.Set("region_id", regionID)
//
//	_, err = waitForLBAvailable(d, "active", []string{"pending", "installing"}, "status", meta)
//	if err != nil {
//		return fmt.Errorf(
//			"Error while waiting for load balancer %v to be completed: %v", id, err)
//	}
//
//	log.Printf("[INFO] load balancer ID: %v", id)
//
//	return resourceVultrLoadBalancerRead(d, meta)
//}
//
//func resourceVultrLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//	id, _ := strconv.Atoi(d.Id())
//
//	lbConfig, err := client.LoadBalancer.GetFullConfig(context.Background(), id)
//	if err != nil {
//		log.Printf("[WARN] Vultr load balancer (%v) not found", id)
//		d.SetId("")
//		return nil
//	}
//
//	lbs, err := client.LoadBalancer.List(context.Background())
//	if err != nil {
//		return fmt.Errorf("Error getting load balancer: %v", err)
//	}
//
//	var lb *govultr.LoadBalancers
//	for i := range lbs {
//		if lbs[i].ID == id {
//			lb = &lbs[i]
//			break
//		}
//	}
//
//	if lb == nil {
//		log.Printf("[WARN] Vultr load balancer (%v) not found", id)
//		d.SetId("")
//		return nil
//	}
//
//	var rulesList []map[string]interface{}
//	for _, rules := range lbConfig.ForwardRuleList {
//		rule := map[string]interface{}{
//			"rule_id":           rules.RuleID,
//			"frontend_protocol": rules.FrontendProtocol,
//			"frontend_port":     rules.FrontendPort,
//			"backend_protocol":  rules.BackendProtocol,
//			"backend_port":      rules.BackendPort,
//		}
//		rulesList = append(rulesList, rule)
//	}
//
//	if err := d.Set("forwarding_rules", rulesList); err != nil {
//		return fmt.Errorf("Error setting `forwarding_rules`: %v", err)
//	}
//
//	var hc []map[string]interface{}
//	hcInfo := map[string]interface{}{
//		"protocol":            lbConfig.HealthCheck.Protocol,
//		"port":                lbConfig.HealthCheck.Port,
//		"path":                lbConfig.HealthCheck.Path,
//		"check_interval":      lbConfig.HealthCheck.CheckInterval,
//		"response_timeout":    lbConfig.HealthCheck.ResponseTimeout,
//		"unhealthy_threshold": lbConfig.HealthCheck.UnhealthyThreshold,
//		"healthy_threshold":   lbConfig.HealthCheck.HealthyThreshold,
//	}
//	hc = append(hc, hcInfo)
//
//	if err := d.Set("health_check", hc); err != nil {
//		return fmt.Errorf("Error setting `health_check`: %v", err)
//	}
//
//	d.Set("has_ssl", lbConfig.SSLInfo)
//
//	d.Set("attached_instances", lbConfig.InstanceList)
//	d.Set("balancing_algorithm", lbConfig.BalancingAlgorithm)
//	d.Set("ssl_redirect", lbConfig.SSLRedirect)
//	d.Set("proxy_protocol", lbConfig.ProxyProtocol)
//	d.Set("cookie_name", lbConfig.StickySessions.CookieName)
//	d.Set("date_created", lb.DateCreated)
//	d.Set("label", lb.Label)
//	d.Set("status", lb.Status)
//	d.Set("ipv4", lb.IPV4)
//	d.Set("ipv6", lb.IPV6)
//
//	return nil
//}
//
//func resourceVultrLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//	id, _ := strconv.Atoi(d.Id())
//
//	if d.HasChanges("balancing_algorithm", "proxy_protocol", "ssl_redirect", "cookie_name") {
//		genericInfo := govultr.GenericInfo{}
//		if d.HasChange("balancing_algorithm") {
//			genericInfo.BalancingAlgorithm = d.Get("balancing_algorithm").(string)
//		}
//
//		if d.HasChange("proxy_protocol") {
//			proxyProtocol := d.Get("proxy_protocol").(bool)
//			genericInfo.ProxyProtocol = &proxyProtocol
//		}
//
//		if d.HasChange("ssl_redirect") {
//			sslRedirect := d.Get("ssl_redirect").(bool)
//			genericInfo.SSLRedirect = &sslRedirect
//		}
//
//		if d.HasChange("cookie_name") {
//			cookieName := d.Get("cookie_name").(string)
//			s := &govultr.StickySessions{}
//			if cookieName == "" {
//				s.StickySessionsEnabled = "off"
//			} else {
//				s.StickySessionsEnabled = "on"
//				s.CookieName = cookieName
//			}
//			genericInfo.StickySessions = s
//		}
//		log.Printf(`[INFO] Updating load balancer generic info (%v)`, id)
//		err := client.LoadBalancer.UpdateGenericInfo(context.Background(), id, d.Get("label").(string), &genericInfo)
//		if err != nil {
//			return fmt.Errorf("error updating load balancer generic info (%v): %v", id, err)
//		}
//	}
//
//	if d.HasChange("health_check") {
//		health := d.Get("health_check")
//		check := generateHealthCheck(health)
//
//		log.Printf(`[INFO] Updating load balancer health info (%v)`, id)
//		err := client.LoadBalancer.SetHealthCheck(context.Background(), id, check)
//		if err != nil {
//			return fmt.Errorf("Error updating load balancer health info (%v): %v", id, err)
//		}
//	}
//
//	if d.HasChange("ssl") {
//		if sslData, sslOK := d.GetOk("ssl"); sslOK {
//			ssl := generateSSL(sslData)
//
//			if d.Get("has_ssl").(bool) {
//				log.Printf(`[INFO] Removing load balancer SSL certificate (%v)`, id)
//				err := client.LoadBalancer.RemoveSSL(context.Background(), id)
//				if err != nil {
//					return fmt.Errorf("Error removing SSL certificate for load balancer (%v): %v", id, err)
//				}
//			}
//			err := client.LoadBalancer.AddSSL(context.Background(), id, ssl)
//			if err != nil {
//				//reload the original SSL cert to prevent drift
//				original, _ := d.GetChange("ssl")
//				ssl = generateSSL(original)
//				client.LoadBalancer.AddSSL(context.Background(), id, ssl)
//
//				return fmt.Errorf("Error adding SSL certificate for load balancer (%v): %v", id, err)
//			}
//		} else {
//			log.Printf(`[INFO] Removing load balancer SSL certificate (%v)`, id)
//			err := client.LoadBalancer.RemoveSSL(context.Background(), id)
//			if err != nil {
//				return fmt.Errorf("error removing SSL certificate for load balancer (%v): %v", id, err)
//			}
//		}
//	}
//
//	if d.HasChange("forwarding_rules") {
//		oldFR, newFR := d.GetChange("forwarding_rules")
//
//		frCreate := newFR.(*schema.Set).Difference(oldFR.(*schema.Set))
//		frDelete := oldFR.(*schema.Set).Difference(newFR.(*schema.Set))
//
//		for _, val := range frDelete.List() {
//			t := val.(map[string]interface{})
//
//			err := client.LoadBalancer.DeleteForwardingRule(context.Background(), id, t["rule_id"].(string))
//			if err != nil {
//				return fmt.Errorf("error removing forwarding rules for load balancer (%v): %v", id, err)
//			}
//		}
//
//		for _, val := range frCreate.List() {
//			t := val.(map[string]interface{})
//
//			rule := &govultr.ForwardingRule{
//				FrontendProtocol: t["frontend_protocol"].(string),
//				FrontendPort:     t["frontend_port"].(int),
//				BackendProtocol:  t["backend_protocol"].(string),
//				BackendPort:      t["backend_port"].(int),
//			}
//
//			_, err := client.LoadBalancer.CreateForwardingRule(context.Background(), id, rule)
//			if err != nil {
//				return fmt.Errorf("error adding forwarding rules for load balancer (%v): %v", id, err)
//			}
//		}
//	}
//
//	if d.HasChange("attached_instances") {
//		oldInstances, newInstances := d.GetChange("attached_instances")
//
//		log.Println(oldInstances)
//		log.Println(newInstances)
//
//		var oldIDs []int
//		for _, v := range oldInstances.([]interface{}) {
//			oldIDs = append(oldIDs, v.(int))
//		}
//
//		var newIDs []int
//		for _, v := range newInstances.([]interface{}) {
//			newIDs = append(newIDs, v.(int))
//		}
//
//		if len(oldIDs) != 0 {
//			for _, v := range diff(newIDs, oldIDs) {
//				err := client.LoadBalancer.DetachInstance(context.Background(), id, v)
//				if err != nil {
//					return fmt.Errorf("error detaching instance id %v from load balancer (%v): %v", v, id, err)
//				}
//			}
//		}
//
//		if len(newIDs) != 0 {
//			for _, v := range diff(oldIDs, newIDs) {
//				err := client.LoadBalancer.AttachInstance(context.Background(), id, v)
//				if err != nil {
//					return fmt.Errorf("error attaching instance id %v to load balancer (%v): %v", v, id, err)
//				}
//			}
//		}
//
//	}
//
//	return resourceVultrLoadBalancerRead(d, meta)
//}
//
//func resourceVultrLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	log.Printf("[INFO] Deleting load balancer: %v", d.Id())
//
//	id, _ := strconv.Atoi(d.Id())
//	err := client.LoadBalancer.Delete(context.Background(), id)
//
//	if err != nil {
//		return fmt.Errorf("Error deleting load balancer %v : %v", d.Id(), err)
//	}
//
//	return nil
//}
//
//func waitForLBAvailable(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
//	log.Printf(
//		"[INFO] Waiting for load balancer (%s) to have %s of %s",
//		d.Id(), attribute, target)
//
//	stateConf := &resource.StateChangeConf{
//		Pending:        pending,
//		Target:         []string{target},
//		Refresh:        newLBStateRefresh(d, meta, attribute),
//		Timeout:        60 * time.Minute,
//		Delay:          10 * time.Second,
//		MinTimeout:     5 * time.Second,
//		NotFoundChecks: 60,
//	}
//
//	return stateConf.WaitForState()
//}
//
//func newLBStateRefresh(d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc {
//	client := meta.(*Client).govultrClient()
//	return func() (interface{}, string, error) {
//
//		log.Printf("[INFO] Creating load balancer")
//
//		id, _ := strconv.Atoi(d.Id())
//		lbs, err := client.LoadBalancer.List(context.Background())
//		if err != nil {
//			return nil, "", fmt.Errorf("Error getting load balancer: %v", err)
//		}
//
//		var lb *govultr.LoadBalancers
//		for i := range lbs {
//			if lbs[i].ID == id {
//				lb = &lbs[i]
//				break
//			}
//		}
//
//		if lb == nil {
//			return nil, "", fmt.Errorf("error retrieving lb %s ", d.Id())
//		}
//
//		if attr == "status" {
//			log.Printf("[INFO] The load balancer Status is %v", lb.Status)
//			return lb, lb.Status, nil
//		} else {
//			return nil, "", nil
//		}
//	}
//}
//
//func generateRules(rules interface{}) *govultr.ForwardingRules {
//	fwMap := &govultr.ForwardingRules{}
//	for _, rule := range rules.(*schema.Set).List() {
//		r := rule.(map[string]interface{})
//		t := govultr.ForwardingRule{
//			FrontendProtocol: r["frontend_protocol"].(string),
//			FrontendPort:     r["frontend_port"].(int),
//			BackendProtocol:  r["backend_protocol"].(string),
//			BackendPort:      r["backend_port"].(int),
//		}
//		fwMap.ForwardRuleList = append(fwMap.ForwardRuleList, t)
//	}
//	return fwMap
//}
//
//func generateHealthCheck(health interface{}) *govultr.HealthCheck {
//	k := health.([]interface{})
//	config := k[0].(map[string]interface{})
//
//	return &govultr.HealthCheck{
//		Protocol:           config["protocol"].(string),
//		Port:               config["port"].(int),
//		Path:               config["path"].(string),
//		CheckInterval:      config["check_interval"].(int),
//		ResponseTimeout:    config["response_timeout"].(int),
//		UnhealthyThreshold: config["unhealthy_threshold"].(int),
//		HealthyThreshold:   config["healthy_threshold"].(int),
//	}
//}
//
//func generateSSL(sslData interface{}) *govultr.SSL {
//	k := sslData.(*schema.Set).List()
//	config := k[0].(map[string]interface{})
//
//	return &govultr.SSL{
//		PrivateKey:  config["private_key"].(string),
//		Certificate: config["certificate"].(string),
//		Chain:       config["chain"].(string),
//	}
//}
//
//func diff(in, out []int) []int {
//	var diff []int
//
//	b := map[int]int{}
//	for i := range in {
//		b[in[i]] = 0
//	}
//
//	for i := range out {
//		if _, ok := b[out[i]]; !ok {
//			diff = append(diff, out[i])
//		}
//	}
//	return diff
//}
