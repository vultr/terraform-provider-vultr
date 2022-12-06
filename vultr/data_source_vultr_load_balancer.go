package vultr

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrLoadBalancerRead,
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
			"proxy_protocol": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"firewall_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"private_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}
	var lbList []govultr.LoadBalancer
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		lbs, meta, err := client.LoadBalancer.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting load balancer: %v", err)
		}

		for _, b := range lbs {
			sm, err := structToMap(b)

			if err != nil {
				return diag.FromErr(err)
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
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(lbList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(lbList[0].ID)
	if err := d.Set("has_ssl", lbList[0].SSLInfo); err != nil {
		return diag.Errorf("unable to set load_balancer `has_ssl` read value: %v", err)
	}
	if err := d.Set("attached_instances", lbList[0].Instances); err != nil {
		return diag.Errorf("unable to set load_balancer `attached_instances` read value: %v", err)
	}
	if err := d.Set("balancing_algorithm", lbList[0].GenericInfo.BalancingAlgorithm); err != nil {
		return diag.Errorf("unable to set load_balancer `balancing_algorithm` read value: %v", err)
	}
	if err := d.Set("ssl_redirect", lbList[0].GenericInfo.SSLRedirect); err != nil {
		return diag.Errorf("unable to set load_balancer `ssl_redirect` read value: %v", err)
	}
	if err := d.Set("proxy_protocol", lbList[0].GenericInfo.ProxyProtocol); err != nil {
		return diag.Errorf("unable to set load_balancer `proxy_protocol` read value: %v", err)
	}
	if err := d.Set("cookie_name", lbList[0].GenericInfo.StickySessions.CookieName); err != nil {
		return diag.Errorf("unable to set load_balancer `cookie_name` read value: %v", err)
	}
	if err := d.Set("date_created", lbList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set load_balancer `date_created` read value: %v", err)
	}
	if err := d.Set("status", lbList[0].Status); err != nil {
		return diag.Errorf("unable to set load_balancer `status` read value: %v", err)
	}
	if err := d.Set("region", lbList[0].Region); err != nil {
		return diag.Errorf("unable to set load_balancer `region` read value: %v", err)
	}
	if err := d.Set("label", lbList[0].Label); err != nil {
		return diag.Errorf("unable to set load_balancer `label` read value: %v", err)
	}
	if err := d.Set("ipv4", lbList[0].IPV4); err != nil {
		return diag.Errorf("unable to set load_balancer `ipv4` read value: %v", err)
	}
	if err := d.Set("ipv6", lbList[0].IPV6); err != nil {
		return diag.Errorf("unable to set load_balancer `ipv6` read value: %v", err)
	}
	if err := d.Set("private_network", lbList[0].GenericInfo.PrivateNetwork); err != nil {
		return diag.Errorf("unable to set load_balancer `private_network` read value: %v", err)
	}

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
		return diag.Errorf("unable to set load_balancer `forwarding_rules` read value: %v", err)
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
		return diag.Errorf("unable to set load_balancer `health_check` read value: %v", err)
	}

	var fwrRules []map[string]interface{}
	for _, rules := range lbList[0].FirewallRules {
		rule := map[string]interface{}{
			"id":      rules.RuleID,
			"ip_type": rules.IPType,
			"port":    strconv.Itoa(rules.Port),
			"source":  rules.Source,
		}
		fwrRules = append(fwrRules, rule)
	}

	if err := d.Set("firewall_rules", fwrRules); err != nil {
		return diag.Errorf("unable to set load_balancer `firewall_rules` read value: %v", err)
	}

	return nil
}
