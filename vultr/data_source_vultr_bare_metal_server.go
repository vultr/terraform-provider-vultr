package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrBareMetalServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrBareMetalServerRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"netmask_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"v6_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrBareMetalServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	servers, err := client.BareMetalServer.List(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting bare metal servers: %v", err)
	}

	serverList := []govultr.BareMetalServer{}

	f := buildVultrDataSourceFilter(filters.(*schema.Set))

	for _, s := range servers {
		// we need convert the a struct INTO a map so we can easily manipulate the data here
		sm, err := structToMap(s)

		if err != nil {
			return err
		}

		if filterLoop(f, sm) {
			serverList = append(serverList, s)
		}
	}

	if len(serverList) > 1 {
		return errors.New("your search returned too many results. Please refine your search to be more specific")
	}

	if len(serverList) < 1 {
		return errors.New("no results were found")
	}

	d.SetId(serverList[0].BareMetalServerID)
	d.Set("os", serverList[0].Os)
	d.Set("ram", serverList[0].RAM)
	d.Set("disk", serverList[0].Disk)
	d.Set("main_ip", serverList[0].MainIP)
	d.Set("cpu_count", serverList[0].CPUCount)
	d.Set("location", serverList[0].Location)
	d.Set("region_id", serverList[0].RegionID)
	d.Set("default_password", serverList[0].DefaultPassword)
	d.Set("date_created", serverList[0].DateCreated)
	d.Set("status", serverList[0].Status)
	d.Set("netmask_v4", serverList[0].NetmaskV4)
	d.Set("gateway_v4", serverList[0].GatewayV4)
	d.Set("plan_id", serverList[0].BareMetalPlanID)
	d.Set("label", serverList[0].Label)
	d.Set("tag", serverList[0].Tag)
	d.Set("os_id", serverList[0].OsID)
	d.Set("app_id", serverList[0].AppID)

	var ipv6s []map[string]string
	for _, net := range serverList[0].V6Networks {
		v6network := map[string]string{
			"v6_network":      net.Network,
			"v6_main_ip":      net.MainIP,
			"v6_network_size": net.NetworkSize,
		}
		ipv6s = append(ipv6s, v6network)
	}
	d.Set("v6_networks", ipv6s)

	return nil
}
