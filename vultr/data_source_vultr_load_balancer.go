package vultr

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
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
			"region": {
				Type:     schema.TypeString,
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
			"has_ssl": {
				Type:     schema.TypeBool,
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
	var lbList []govultr.LoadBalancer
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		lbs, meta, err := client.LoadBalancer.List(context.Background(), options)
		if err != nil {
			return fmt.Errorf("error getting load balancer: %v", err)
		}

		for _, b := range lbs {
			sm, err := structToMap(b)

			if err != nil {
				return err
			}

			if filterLoop(f, sm) {
				lbList = append(lbList, b)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}
	if len(lbList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(lbList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(lbList[0].ID)
	d.Set("has_ssl", lbList[0].SSLInfo)
	d.Set("attached_instances", lbList[0].Instances)
	d.Set("balancing_algorithm", lbList[0].GenericInfo.BalancingAlgorithm)
	d.Set("ssl_redirect", lbList[0].GenericInfo.SSLRedirect)
	d.Set("proxy_protocol", lbList[0].GenericInfo.ProxyProtocol)
	d.Set("cookie_name", lbList[0].GenericInfo.StickySessions.CookieName)
	d.Set("date_created", lbList[0].DateCreated)
	d.Set("status", lbList[0].Status)
	d.Set("region", lbList[0].Region)
	d.Set("label", lbList[0].Label)
	d.Set("ipv4", lbList[0].IPV4)
	d.Set("ipv6", lbList[0].IPV6)

	var rulesList []map[string]interface{}
	for _, rules := range lbList[0].ForwardingRules {
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
		"protocol":            lbList[0].HealthCheck.Protocol,
		"port":                strconv.Itoa(lbList[0].HealthCheck.Port),
		"path":                lbList[0].HealthCheck.Path,
		"check_interval":      strconv.Itoa(lbList[0].HealthCheck.CheckInterval),
		"response_timeout":    strconv.Itoa(lbList[0].HealthCheck.ResponseTimeout),
		"unhealthy_threshold": strconv.Itoa(lbList[0].HealthCheck.UnhealthyThreshold),
		"healthy_threshold":   strconv.Itoa(lbList[0].HealthCheck.HealthyThreshold),
	}

	if err := d.Set("health_check", hcInfo); err != nil {
		return fmt.Errorf("error setting `health_check`: %#v", err)
	}

	return nil
}
