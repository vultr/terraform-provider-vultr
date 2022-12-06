package vultr

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrInstanceRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"main_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allowed_bandwidth": {
				Type:     schema.TypeInt,
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
			"status": {
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
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v6_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v6_main_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v6_network_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kvm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backups": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"os_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewall_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"backups_schedule": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_network_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVultrInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var serverList []govultr.Instance
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		servers, meta, err := client.Instance.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting servers: %v", err)
		}

		for _, s := range servers {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(s)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				serverList = append(serverList, s)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(serverList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(serverList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(serverList[0].ID)
	if err := d.Set("os", serverList[0].Os); err != nil {
		return diag.Errorf("unable to set instance `os` read value: %v", err)
	}
	if err := d.Set("ram", serverList[0].RAM); err != nil {
		return diag.Errorf("unable to set instance `ram` read value: %v", err)
	}
	if err := d.Set("disk", serverList[0].Disk); err != nil {
		return diag.Errorf("unable to set instance `disk` read value: %v", err)
	}
	if err := d.Set("main_ip", serverList[0].MainIP); err != nil {
		return diag.Errorf("unable to set instance `main_ip` read value: %v", err)
	}
	if err := d.Set("vcpu_count", serverList[0].VCPUCount); err != nil {
		return diag.Errorf("unable to set instance `vcpu_count` read value: %v", err)
	}
	if err := d.Set("region", serverList[0].Region); err != nil {
		return diag.Errorf("unable to set instance `region` read value: %v", err)
	}
	if err := d.Set("date_created", serverList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set instance `date_created` read value: %v", err)
	}
	if err := d.Set("allowed_bandwidth", serverList[0].AllowedBandwidth); err != nil {
		return diag.Errorf("unable to set instance `allowed_bandwidth` read value: %v", err)
	}
	if err := d.Set("netmask_v4", serverList[0].NetmaskV4); err != nil {
		return diag.Errorf("unable to set instance `netmask_v4` read value: %v", err)
	}
	if err := d.Set("gateway_v4", serverList[0].GatewayV4); err != nil {
		return diag.Errorf("unable to set instance `gateway_v4` read value: %v", err)
	}
	if err := d.Set("status", serverList[0].Status); err != nil {
		return diag.Errorf("unable to set instance `status` read value: %v", err)
	}
	if err := d.Set("power_status", serverList[0].PowerStatus); err != nil {
		return diag.Errorf("unable to set instance `power_status` read value: %v", err)
	}
	if err := d.Set("server_status", serverList[0].ServerStatus); err != nil {
		return diag.Errorf("unable to set instance `server_status` read value: %v", err)
	}
	if err := d.Set("plan", serverList[0].Plan); err != nil {
		return diag.Errorf("unable to set instance `plan` read value: %v", err)
	}
	if err := d.Set("label", serverList[0].Label); err != nil {
		return diag.Errorf("unable to set instance `label` read value: %v", err)
	}
	if err := d.Set("internal_ip", serverList[0].InternalIP); err != nil {
		return diag.Errorf("unable to set instance `internal_ip` read value: %v", err)
	}
	if err := d.Set("kvm", serverList[0].KVM); err != nil {
		return diag.Errorf("unable to set instance `kvm` read value: %v", err)
	}
	if err := d.Set("tags", serverList[0].Tags); err != nil {
		return diag.Errorf("unable to set instance `tags` read value: %v", err)
	}
	if err := d.Set("os_id", serverList[0].OsID); err != nil {
		return diag.Errorf("unable to set instance `os_id` read value: %v", err)
	}
	if err := d.Set("app_id", serverList[0].AppID); err != nil {
		return diag.Errorf("unable to set instance `app_id` read value: %v", err)
	}
	if err := d.Set("image_id", serverList[0].ImageID); err != nil {
		return diag.Errorf("unable to set instance `image_id` read value: %v", err)
	}
	if err := d.Set("firewall_group_id", serverList[0].FirewallGroupID); err != nil {
		return diag.Errorf("unable to set instance `firewall_group_id` read value: %v", err)
	}
	if err := d.Set("v6_network", serverList[0].V6Network); err != nil {
		return diag.Errorf("unable to set instance `v6_network` read value: %v", err)
	}
	if err := d.Set("v6_main_ip", serverList[0].V6MainIP); err != nil {
		return diag.Errorf("unable to set instance `v6_main_ip` read value: %v", err)
	}
	if err := d.Set("v6_network_size", serverList[0].V6NetworkSize); err != nil {
		return diag.Errorf("unable to set instance `v6_network_size` read value: %v", err)
	}
	if err := d.Set("features", serverList[0].Features); err != nil {
		return diag.Errorf("unable to set instance `features` read value: %v", err)
	}
	if err := d.Set("hostname", serverList[0].Hostname); err != nil {
		return diag.Errorf("unable to set instance `hostname` read value: %v", err)
	}

	schedule, err := client.Instance.GetBackupSchedule(ctx, serverList[0].ID)
	if err != nil {
		return diag.Errorf("error getting backup schedule: %v", err)
	}
	if err := d.Set("backups", backupStatus(schedule.Enabled)); err != nil {
		return diag.Errorf("unable to set instance `backups` read value: %v", err)
	}

	bsInfo := map[string]interface{}{
		"type": schedule.Type,
		"hour": strconv.Itoa(schedule.Hour),
		"dom":  strconv.Itoa(schedule.Dom),
		"dow":  strconv.Itoa(schedule.Dow),
	}
	if err := d.Set("backups_schedule", bsInfo); err != nil {
		return diag.Errorf("error setting `backups_schedule`: %#v", err)
	}

	vpcs, err := getVPCs(client, d.Id())
	if err != nil {
		return diag.Errorf(err.Error())
	}

	if err := d.Set("private_network_ids", vpcs); err != nil {
		return diag.Errorf("unable to set instance `private_network_ids` read value: %v", err)
	}
	if err := d.Set("vpc_ids", vpcs); err != nil {
		return diag.Errorf("unable to set instance `vpc_ids` read value: %v", err)
	}

	return nil
}
