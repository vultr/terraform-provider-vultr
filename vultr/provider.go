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
			"vultr_account":     dataSourceVultrAccount(),
			"vultr_api":         dataSourceVultrApi(),
			"vultr_application": dataSourceVultrApplication(),
			"vultr_os":          dataSourceVultrOS(),
			"vultr_user":        dataSourceVultrUser(),
		},

		ResourcesMap:  map[string]*schema.Resource{},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey: d.Get("api_key").(string),
	}

	return config.Client()
}
