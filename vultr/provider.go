package vultr

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VULTR_API_KEY", nil),
				Description: "The API Key that allows interaction with the API",
			},
			"rate_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Allows users to set the speed of API calls to work with the Vultr Rate Limit",
			},
			"retry_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Allows users to set the maximum number of retries allowed for a failed API call.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			//"vultr_account":           dataSourceVultrAccount(),
			//"vultr_api_key":           dataSourceVultrApi(),
			//"vultr_application":       dataSourceVultrApplication(),
			//"vultr_backup":            dataSourceVultrBackup(),
			//"vultr_bare_metal_plan":   dataSourceVultrBareMetalPlan(),
			//"vultr_bare_metal_server": dataSourceVultrBareMetalServer(),
			//"vultr_block_storage":     dataSourceVultrBlockStorage(),
			"vultr_dns_domain": dataSourceVultrDnsDomain(),
			//"vultr_firewall_group":    dataSourceVultrFirewallGroup(),
			//"vultr_iso_private":       dataSourceVultrIsoPrivate(),
			//"vultr_iso_public":        dataSourceVultrIsoPublic(),
			//"vultr_load_balancer":     dataSourceVultrLoadBalancer(),
			//"vultr_network":           dataSourceVultrNetwork(),
			//"vultr_object_storage":    dataSourceVultrObjectStorage(),
			"vultr_os": dataSourceVultrOS(),
			//"vultr_plan":              dataSourceVultrPlan(),
			//"vultr_region":            dataSourceVultrRegion(),
			//"vultr_reserved_ip":       dataSourceVultrReservedIP(),
			//"vultr_reverse_ipv4":      dataSourceVultrReverseIPV4(),
			//"vultr_reverse_ipv6":      dataSourceVultrReverseIPV6(),
			//"vultr_server":            dataSourceVultrServer(),
			//"vultr_server_ipv4":       dataSourceVultrServerIPV4(),
			//"vultr_snapshot":          dataSourceVultrSnapshot(),
			//"vultr_ssh_key":           dataSourceVultrSSHKey(),
			//"vultr_startup_script":    dataSourceVultrStartupScript(),
			"vultr_user": dataSourceVultrUser(),
		},

		ResourcesMap: map[string]*schema.Resource{
			//"vultr_bare_metal_server": resourceVultrBareMetalServer(),
			"vultr_block_storage":  resourceVultrBlockStorage(),
			"vultr_dns_domain":     resourceVultrDnsDomain(),
			"vultr_dns_record":     resourceVultrDnsRecord(),
			"vultr_firewall_group": resourceVultrFirewallGroup(),
			"vultr_firewall_rule":  resourceVultrFirewallRule(),
			//"vultr_iso_private":       resourceVultrIsoPrivate(),
			"vultr_load_balancer":     resourceVultrLoadBalancer(),
			"vultr_private_network":   resourceVultrPrivateNetwork(),
			"vultr_object_storage":    resourceVultrObjectStorage(),
			"vultr_reserved_ip":       resourceVultrReservedIP(),
			"vultr_reverse_ipv4":      resourceVultrReverseIPV4(),
			"vultr_reverse_ipv6":      resourceVultrReverseIPV6(),
			"vultr_snapshot":          resourceVultrSnapshot(),
			"vultr_snapshot_from_url": resourceVultrSnapshotFromURL(),
			//"vultr_server":            resourceVultrServer(),
			//"vultr_server_ipv4":       resourceVultrServerIPV4(),
			"vultr_ssh_key":        resourceVultrSSHKey(),
			"vultr_startup_script": resourceVultrStartupScript(),
			"vultr_user":           resourceVultrUsers(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey:     d.Get("api_key").(string),
		RateLimit:  d.Get("rate_limit").(int),
		RetryLimit: d.Get("retry_limit").(int),
	}

	return config.Client()
}
