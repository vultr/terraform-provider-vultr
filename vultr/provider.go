package vultr

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the base Vultr terraform provider
func Provider() *schema.Provider {
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
			"vultr_account":                     dataSourceVultrAccount(),
			"vultr_application":                 dataSourceVultrApplication(),
			"vultr_backup":                      dataSourceVultrBackup(),
			"vultr_bare_metal_plan":             dataSourceVultrBareMetalPlan(),
			"vultr_bare_metal_server":           dataSourceVultrBareMetalServer(),
			"vultr_block_storage":               dataSourceVultrBlockStorage(),
			"vultr_container_registry":          dataSourceVultrContainerRegistry(),
			"vultr_database":                    dataSourceVultrDatabase(),
			"vultr_dns_domain":                  dataSourceVultrDNSDomain(),
			"vultr_firewall_group":              dataSourceVultrFirewallGroup(),
			"vultr_inference":                   dataSourceVultrInference(),
			"vultr_iso_private":                 dataSourceVultrIsoPrivate(),
			"vultr_iso_public":                  dataSourceVultrIsoPublic(),
			"vultr_kubernetes":                  dataSourceVultrKubernetes(),
			"vultr_load_balancer":               dataSourceVultrLoadBalancer(),
			"vultr_object_storage":              dataSourceVultrObjectStorage(),
			"vultr_object_storage_cluster":      dataSourceVultrObjectStorageClusters(),
			"vultr_object_storage_tier":         dataSourceVultrObjectStorageTier(),
			"vultr_os":                          dataSourceVultrOS(),
			"vultr_plan":                        dataSourceVultrPlan(),
			"vultr_region":                      dataSourceVultrRegion(),
			"vultr_reserved_ip":                 dataSourceVultrReservedIP(),
			"vultr_reverse_ipv4":                dataSourceVultrReverseIPV4(),
			"vultr_reverse_ipv6":                dataSourceVultrReverseIPV6(),
			"vultr_instance":                    dataSourceVultrInstance(),
			"vultr_instances":                   dataSourceVultrInstances(),
			"vultr_instance_ipv4":               dataSourceVultrInstanceIPV4(),
			"vultr_snapshot":                    dataSourceVultrSnapshot(),
			"vultr_ssh_key":                     dataSourceVultrSSHKey(),
			"vultr_startup_script":              dataSourceVultrStartupScript(),
			"vultr_user":                        dataSourceVultrUser(),
			"vultr_virtual_file_system_storage": dataSourceVultrVirtualFileSystemStorage(),
			"vultr_vpc":                         dataSourceVultrVPC(),
			"vultr_vpc2":                        dataSourceVultrVPC2(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vultr_bare_metal_server":           resourceVultrBareMetalServer(),
			"vultr_block_storage":               resourceVultrBlockStorage(),
			"vultr_container_registry":          resourceVultrContainerRegistry(),
			"vultr_database":                    resourceVultrDatabase(),
			"vultr_database_connection_pool":    resourceVultrDatabaseConnectionPool(),
			"vultr_database_db":                 resourceVultrDatabaseDB(),
			"vultr_database_replica":            resourceVultrDatabaseReplica(),
			"vultr_database_user":               resourceVultrDatabaseUser(),
			"vultr_database_topic":              resourceVultrDatabaseTopic(),
			"vultr_database_quota":              resourceVultrDatabaseQuota(),
			"vultr_database_connector":          resourceVultrDatabaseConnector(),
			"vultr_dns_domain":                  resourceVultrDNSDomain(),
			"vultr_dns_record":                  resourceVultrDNSRecord(),
			"vultr_firewall_group":              resourceVultrFirewallGroup(),
			"vultr_firewall_rule":               resourceVultrFirewallRule(),
			"vultr_inference":                   resourceVultrInference(),
			"vultr_iso_private":                 resourceVultrIsoPrivate(),
			"vultr_kubernetes":                  resourceVultrKubernetes(),
			"vultr_kubernetes_node_pools":       resourceVultrKubernetesNodePools(),
			"vultr_load_balancer":               resourceVultrLoadBalancer(),
			"vultr_nat_gateway":                 resourceVultrNATGateway(),
			"vultr_object_storage":              resourceVultrObjectStorage(),
			"vultr_reserved_ip":                 resourceVultrReservedIP(),
			"vultr_reverse_ipv4":                resourceVultrReverseIPV4(),
			"vultr_reverse_ipv6":                resourceVultrReverseIPV6(),
			"vultr_snapshot":                    resourceVultrSnapshot(),
			"vultr_snapshot_from_url":           resourceVultrSnapshotFromURL(),
			"vultr_instance":                    resourceVultrInstance(),
			"vultr_instance_ipv4":               resourceVultrInstanceIPV4(),
			"vultr_ssh_key":                     resourceVultrSSHKey(),
			"vultr_startup_script":              resourceVultrStartupScript(),
			"vultr_user":                        resourceVultrUsers(),
			"vultr_virtual_file_system_storage": resourceVultrVirtualFileSystemStorage(),
			"vultr_vpc":                         resourceVultrVPC(),
			"vultr_vpc2":                        resourceVultrVPC2(),
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
