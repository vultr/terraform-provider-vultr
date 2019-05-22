package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

func dataSourceVultrServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVultrServerRead,
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
			"vps_cpu_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pending_charges": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost_per_month": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_bandwidth": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"allowed_bandwidth": {
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
			"power_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_id": {
				Type:     schema.TypeString,
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
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kvm_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_backups": {
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
			"firewall_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("issue with filter: %v", filtersOk)
	}

	servers, err := client.Server.GetList(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting servers: %v", err)
	}

	serverList := []govultr.Server{}

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

	d.SetId(serverList[0].VpsID)
	d.Set("os", serverList[0].Os)
	d.Set("ram", serverList[0].RAM)
	d.Set("disk", serverList[0].Disk)
	d.Set("main_ip", serverList[0].MainIP)
	d.Set("vps_cpu_count", serverList[0].VPSCpus)
	d.Set("location", serverList[0].Location)
	d.Set("region_id", serverList[0].RegionID)
	d.Set("date_created", serverList[0].Created)
	d.Set("pending_charges", serverList[0].PendingCharges)
	d.Set("cost_per_month", serverList[0].Cost)
	d.Set("current_bandwidth", serverList[0].CurrentBandwidth)
	d.Set("allowed_bandwidth", serverList[0].AllowedBandwidth)
	d.Set("netmask_v4", serverList[0].NetmaskV4)
	d.Set("gateway_v4", serverList[0].GatewayV4)
	d.Set("power_status", serverList[0].PowerStatus)
	d.Set("server_status", serverList[0].ServerState)
	d.Set("plan_id", serverList[0].PlanID)
	d.Set("label", serverList[0].Label)
	d.Set("internal_ip", serverList[0].InternalIP)
	d.Set("kvm_url", serverList[0].KVMUrl)
	d.Set("auto_backups", serverList[0].AutoBackups)
	d.Set("tag", serverList[0].Tag)
	d.Set("os_id", serverList[0].OsID)
	d.Set("app_id", serverList[0].AppID)
	d.Set("firewall_group_id", serverList[0].FirewallGroupID)

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
