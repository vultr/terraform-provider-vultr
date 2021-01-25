package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/vultr/terraform-provider-vultr/vultr"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vultr.Provider,
	})
}
