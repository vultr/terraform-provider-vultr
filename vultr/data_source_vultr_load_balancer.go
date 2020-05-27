package vultr

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr"
	"strconv"
)

func dataSourceVultrLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forwarding_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"balancing_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"ssl": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"ssl_redirect": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"attached_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cookie_name": {
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

func dataSourceVultrLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	lbs, err := client.LoadBalancer.List(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting load balancer: %v", err)
	}

	lbList := []govultr.LoadBalancers{}

	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, b := range lbs {
		sm, err := structToMap(b)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			lbList = append(lbList, b)
		}
	}

	if len(lbList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(lbList) < 1 {
		return errors.New("no results were found")
	}

	id := lbList[0].ID

	lb, err := client.LoadBalancer.GetFullConfig(context.Background(), id)
	if err != nil {
		return fmt.Errorf("error retrieving load balancer configuration (%d): %v", id, err)

	}

	d.SetId(strconv.Itoa(id))
	d.Set("has_ssl", lb.SSLInfo)
	d.Set("attached_instances", lb.InstanceList)
	d.Set("balancing_algorithm", lb.BalancingAlgorithm)
	d.Set("ssl_redirect", lb.SSLRedirect)
	d.Set("proxy_protocol", lb.ProxyProtocol)
	d.Set("cookie_name", lb.StickySessions.CookieName)
	d.Set("date_created", lbList[0].DateCreated)
	d.Set("status", lbList[0].Status)
	d.Set("region_id", lbList[0].RegionID)
	d.Set("label", lbList[0].Label)
	d.Set("ipv4", lbList[0].IPV4)
	d.Set("ipv6", lbList[0].IPV6)

	var rulesList []map[string]interface{}
	for _, rules := range lb.ForwardRuleList {
		rule := map[string]interface{}{
			"rule_id":           rules.RuleID,
			"frontend_protocol": rules.FrontendProtocol,
			"frontend_port":     strconv.Itoa(rules.FrontendPort),
			"backend_protocol":  rules.BackendProtocol,
			"backend_port":      strconv.Itoa(rules.BackendPort),
		}
		rulesList = append(rulesList, rule)
	}

	if err := d.Set("forwarding_rules", rulesList); err != nil {
		return fmt.Errorf("error setting `forwarding_rules`: %#v", err)
	}

	hcInfo := map[string]interface{}{
		"protocol":            lb.HealthCheck.Protocol,
		"port":                strconv.Itoa(lb.HealthCheck.Port),
		"path":                lb.HealthCheck.Path,
		"check_interval":      strconv.Itoa(lb.HealthCheck.CheckInterval),
		"response_timeout":    strconv.Itoa(lb.HealthCheck.ResponseTimeout),
		"unhealthy_threshold": strconv.Itoa(lb.HealthCheck.UnhealthyThreshold),
		"healthy_threshold":   strconv.Itoa(lb.HealthCheck.HealthyThreshold),
	}

	if err := d.Set("health_check", hcInfo); err != nil {
		return fmt.Errorf("error setting `health_check`: %#v", err)
	}

	return nil
}
