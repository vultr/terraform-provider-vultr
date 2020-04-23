package vultr

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVultrLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"cookie_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVultrLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
