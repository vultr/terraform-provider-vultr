package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrBareMetalServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrBareMetalServerRead,
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
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
			"mac_address": {
				Type:     schema.TypeInt,
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
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc2_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"user_scheme": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrBareMetalServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics { //nolint:lll
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	serverList := []govultr.BareMetalServer{}
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		servers, meta, _, err := client.BareMetalServer.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting bare metal servers: %v", err)
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
		return diag.Errorf("unable to set bare_metal_server `os` read value: %v", err)
	}
	if err := d.Set("ram", serverList[0].RAM); err != nil {
		return diag.Errorf("unable to set bare_metal_server `ram` read value: %v", err)
	}
	if err := d.Set("disk", serverList[0].Disk); err != nil {
		return diag.Errorf("unable to set bare_metal_server `disk` read value: %v", err)
	}
	if err := d.Set("main_ip", serverList[0].MainIP); err != nil {
		return diag.Errorf("unable to set bare_metal_server `main_ip` read value: %v", err)
	}
	if err := d.Set("cpu_count", serverList[0].CPUCount); err != nil {
		return diag.Errorf("unable to set bare_metal_server `cpu_count` read value: %v", err)
	}
	if err := d.Set("region", serverList[0].Region); err != nil {
		return diag.Errorf("unable to set bare_metal_server `region` read value: %v", err)
	}
	if err := d.Set("date_created", serverList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set bare_metal_server `date_created` read value: %v", err)
	}
	if err := d.Set("status", serverList[0].Status); err != nil {
		return diag.Errorf("unable to set bare_metal_server `status` read value: %v", err)
	}
	if err := d.Set("netmask_v4", serverList[0].NetmaskV4); err != nil {
		return diag.Errorf("unable to set bare_metal_server `netmask_v4` read value: %v", err)
	}
	if err := d.Set("gateway_v4", serverList[0].GatewayV4); err != nil {
		return diag.Errorf("unable to set bare_metal_server `gateway_v4` read value: %v", err)
	}
	if err := d.Set("plan", serverList[0].Plan); err != nil {
		return diag.Errorf("unable to set bare_metal_server `plan` read value: %v", err)
	}
	if err := d.Set("label", serverList[0].Label); err != nil {
		return diag.Errorf("unable to set bare_metal_server `label` read value: %v", err)
	}
	if err := d.Set("tags", serverList[0].Tags); err != nil {
		return diag.Errorf("unable to set bare_metal_server `tags` read value: %v", err)
	}
	if err := d.Set("mac_address", serverList[0].MacAddress); err != nil {
		return diag.Errorf("unable to set bare_metal_server `mac_address` read value: %v", err)
	}
	if err := d.Set("os_id", serverList[0].OsID); err != nil {
		return diag.Errorf("unable to set bare_metal_server `os_id` read value: %v", err)
	}
	if err := d.Set("app_id", serverList[0].AppID); err != nil {
		return diag.Errorf("unable to set bare_metal_server `app_id` read value: %v", err)
	}
	if err := d.Set("image_id", serverList[0].ImageID); err != nil {
		return diag.Errorf("unable to set bare_metal_server `image_id` read value: %v", err)
	}
	if err := d.Set("snapshot_id", serverList[0].SnapshotID); err != nil {
		return diag.Errorf("unable to set bare_metal_server `snapshot_id` read value: %v", err)
	}
	if err := d.Set("v6_network", serverList[0].V6Network); err != nil {
		return diag.Errorf("unable to set bare_metal_server `v6_network` read value: %v", err)
	}
	if err := d.Set("v6_main_ip", serverList[0].V6MainIP); err != nil {
		return diag.Errorf("unable to set bare_metal_server `v6_main_ip` read value: %v", err)
	}
	if err := d.Set("v6_network_size", serverList[0].V6NetworkSize); err != nil {
		return diag.Errorf("unable to set bare_metal_server `v6_network_size` read value: %v", err)
	}
	if err := d.Set("features", serverList[0].Features); err != nil {
		return diag.Errorf("unable to set bare_metal_server `features` read value: %v", err)
	}
	if err := d.Set("user_scheme", serverList[0].UserScheme); err != nil {
		return diag.Errorf("unable to set bare_metal_server `user_scheme` read value: %v", err)
	}

	vpcInfo, _, err := client.BareMetalServer.ListVPCInfo(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting list of attached vpcs during bare metal server data source read : %v", err)
	}

	// only one VPC ever allowed on bare metal server
	var vpcID = ""
	if len(vpcInfo) != 0 {
		vpcID = vpcInfo[0].ID
	}

	if err := d.Set("vpc_id", vpcID); err != nil {
		return diag.Errorf("unable to set data source bare metal server `vpc_id` read value : %v", err)
	}

	vpc2s, err := getBareMetalServerVPC2s(client, d.Id())
	if err != nil {
		return diag.Errorf("%s", err.Error())
	}

	if err := d.Set("vpc2_ids", vpc2s); err != nil {
		return diag.Errorf("unable to set instance `vpc2_ids` read value: %v", err)
	}

	return nil
}
