package vultr

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vultr_account":           dataSourceVultrAccount(),
			"vultr_api":               dataSourceVultrApi(),
			"vultr_application":       dataSourceVultrApplication(),
			"vultr_backup":            dataSourceVultrBackup(),
			"vultr_bare_metal_plan":   dataSourceVultrBareMetalPlan(),
			"vultr_bare_metal_server": dataSourceVultrBareMetalServer(),
			"vultr_block_storage":     dataSourceVultrBlockStorage(),
			"vultr_dns_domain":        dataSourceVultrDnsDomain(),
			"vultr_firewall_group":    dataSourceVultrFirewallGroup(),
			"vultr_iso_private":       dataSourceVultrIsoPrivate(),
			"vultr_iso_public":        dataSourceVultrIsoPublic(),
			"vultr_os":                dataSourceVultrOS(),
			"vultr_plan":              dataSourceVultrPlan(),
			"vultr_region":            dataSourceVultrRegion(),
			"vultr_reserved_ip":       dataSourceVultrReservedIp(),
			"vultr_server":            dataSourceVultrServer(),
			"vultr_snapshot":          dataSourceVultrSnapshot(),
			"vultr_ssh_key":           dataSourceVultrSSHKey(),
			"vultr_startup_script":    dataSourceVultrStartupScript(),
			"vultr_user":              dataSourceVultrUser(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vultr_dns_domain":        resourceVultrDnsDomain(),
			"vultr_firewall_group":    resourceVultrFirewallGroup(),
			"vultr_firewall_rule":     resourceVultrFirewallRule(),
			"vultr_iso_private":       resourceVultrIsoPrivate(),
			"vultr_reserved_ip":       resourceVultrReservedIP(),
			"vultr_snapshot":          resourceVultrSnapshot(),
			"vultr_snapshot_from_url": resourceVultrSnapshotFromURL(),
			"vultr_ssh_key":           resourceVultrSSHKey(),
			"vultr_startup_script":    resourceVultrStartupScript(),
			"vultr_user":              resourceVultrUsers(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey: d.Get("api_key").(string),
	}

	return config.Client()
}
