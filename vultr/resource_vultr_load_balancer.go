package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
)

func resourceVultrLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrLoadBalancerCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"healthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"unhealthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"check_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frontend_protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backend_protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backend_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"balancing_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_redirect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cookie_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVultrLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	regionID := d.Get("region_id").(int)
	label := d.Get("label").(string)
	sslRedirect := d.Get("ssl_redirect").(bool)

	stickySessions := &govultr.StickySessions{
		StickySessionsEnabled: "false", // tk
		CookieName:            d.Get("cookie_name").(string),
	}

	genericInfo := &govultr.GenericInfo{
		BalancingAlgorithm: d.Get("balancing_algorithm").(string),
		SSLRedirect:        &sslRedirect,
		StickySessions:     stickySessions,
	}

	healthCheck := &govultr.HealthCheck{} // tk
	rules := []govultr.ForwardingRule{}   // tk
	ssl := &govultr.SSL{}                 // tk

	lb, err := client.LoadBalancer.Create(context.Background(), regionID, label, genericInfo, healthCheck, rules, ssl)
	if err != nil {
		return fmt.Errorf("Error creating load balancer: %v", err)
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
