package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

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
			"forwarding_rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"response_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"healthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"unhealthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"check_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"balancing_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_redirect": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cookie_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_chain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_private_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVultrLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()
	regionID := d.Get("region_id").(int)
	label := d.Get("label").(string)
	sslRedirect := d.Get("ssl_redirect").(bool)
	cookieName := strings.TrimSpace(d.Get("cookie_name").(string))

	stickySessions := &govultr.StickySessions{
		StickySessionsEnabled: "false",
		CookieName:            cookieName,
	}

	if cookieName != "" {
		stickySessions.StickySessionsEnabled = "true"
	}

	genericInfo := &govultr.GenericInfo{
		BalancingAlgorithm: d.Get("balancing_algorithm").(string),
		SSLRedirect:        &sslRedirect,
		StickySessions:     stickySessions,
	}

	healthCheck := &govultr.HealthCheck{
		Protocol:           d.Get("protocol").(string),
		Port:               d.Get("port").(int),
		Path:               d.Get("path").(string),
		CheckInterval:      d.Get("check_interval").(int),
		ResponseTimeout:    d.Get("response_timeout").(int),
		UnhealthyThreshold: d.Get("unhealthy_threshold").(int),
		HealthyThreshold:   d.Get("healthy_threshold").(int),
	}

	fr, frOk := d.GetOk("forwarding_rules")
	fwMap := []govultr.ForwardingRule{}
	if frOk {
		for _, value := range fr.([]interface{}) {
			rule := govultr.ForwardingRule{}
			for k, v := range value.(map[string]interface{}) {
				if k == "frontend_port" {
					num, _ := strconv.Atoi(v.(string))
					rule.FrontendPort = num
				} else if k == "frontend_protocol" {
					rule.FrontendProtocol = v.(string)
				} else if k == "backend_port" {
					num, _ := strconv.Atoi(v.(string))
					rule.BackendPort = num
				} else if k == "backend_protocol" {
					rule.BackendProtocol = v.(string)
				}
			}
			fwMap = append(fwMap, rule)
		}
	}

	ssl := &govultr.SSL{
		PrivateKey:  d.Get("ssl_private_key").(string),
		Certificate: d.Get("ssl_certificate").(string),
		Chain:       d.Get("ssl_chain").(string),
	}

	lb, err := client.LoadBalancer.Create(context.Background(), regionID, label, genericInfo, healthCheck, fwMap, ssl)
	if err != nil {
		return fmt.Errorf("Error creating load balancer: %v, %v", err, fwMap)
	}
	id := strconv.Itoa(lb.ID)

	d.SetId(id)
	log.Printf("[INFO] Load Balancer ID: %s", d.Id())

	return resourceVultrLoadBalancerRead(d, meta)
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

	d.Set("date_created", lb.DateCreated)
	d.Set("DCID", lb.RegionID)
	d.Set("location", lb.Location)
	d.Set("label", lb.Label)
	d.Set("status", lb.Status)
	d.Set("ipv4", lb.IPV4)
	d.Set("ipv6", lb.IPV6)

	return nil
}

func resourceVultrLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVultrLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting Load Balancer: %s", d.Id())

	id, _ := strconv.Atoi(d.Id())
	err := client.LoadBalancer.Delete(context.Background(), id)

	if err != nil {
		return fmt.Errorf("Error deleting Load Balancer %s : %v", d.Id(), err)
	}

	return nil
}
