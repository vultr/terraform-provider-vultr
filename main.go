package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/vultr/terraform-provider-vultr/vultr"
)

func main() {
	var debug bool
	var addr string

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.StringVar(&addr, "addr", "", "provider address to match when testing or debugging the provider")
	flag.Parse()

	if debug {
		if addr == "" {
			panic("cannot run in debug without a provider addr: use the -addr flag")
		}

		plugin.Serve(&plugin.ServeOpts{
			Debug:        debug,
			ProviderAddr: addr,
			ProviderFunc: vultr.Provider,
		})

	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: vultr.Provider,
		})
	}
}
