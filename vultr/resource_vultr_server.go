package vultr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/vultr/govultr"
	"log"
	"strconv"
)

const (
	osAppID  = 186
	osIsoID  = 159
	osSnapID = 164
)

func resourceVultrServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrServerCreate,
		Read:   resourceVultrServerRead,
		Update: resourceVultrServerUpdate,
		Delete: resourceVultrServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			//Required
			"region_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"plan_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			// computed attributes
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
			"default_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pending_charges": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
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
			"v6_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kvm_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// server options
			"iso_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"script_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"enable_private_network": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"network_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"auto_backup": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"application_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"os_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"snapshot_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"notify_activate": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"ddos_protection": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"reserved_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"firewall_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceVultrServerCreate(d *schema.ResourceData, meta interface{}) error {

	// Four unique options to image your server
	osID, osOK := d.GetOk("os_id")
	appID, appOK := d.GetOk("application_id")
	isoID, isoOK := d.GetOk("iso_id")
	// This isn't supported on the client needs to implement this
	_, snapOK := d.GetOk("snapshot_id")

	osOptions := map[string]bool{"os_id": osOK, "application_id": appOK, "iso_id": isoOK, "snapshot_id": snapOK}
	osOption, err := optionCheck(osOptions)

	if err != nil {
		return err
	}

	client := meta.(*Client).govultrClient()

	options := &govultr.ServerOptions{
		//IPXEChain: d.Get() check wtf this is
		ScriptID:             d.Get("script_id").(string),
		EnableIPV6:           d.Get("enable_ipv6").(bool),
		EnablePrivateNetwork: d.Get("enable_private_network").(bool),
		Label:                d.Get("label").(string),
		AutoBackups:          d.Get("auto_backup").(bool),
		UserData:             d.Get("user_data").(string),
		NotifyActivate:       d.Get("notify_activate").(bool),
		DDOSProtection:       d.Get("ddos_protection").(bool),
		ReservedIPV4:         d.Get("reserved_ipv4").(string),
		Hostname:             d.Get("hostname").(string),
		Tag:                  d.Get("tag").(string),
		FirewallGroupID:      d.Get("firewall_group_id").(string),
	}

	var os int

	switch osOption {
	case "os_id":
		if osID == osAppID || osID == osIsoID || osID == osSnapID {
			return fmt.Errorf("Please set a corrosponding attribute ")
		}
		os = osID.(int)

	case "application_id":
		options.AppID = strconv.Itoa(appID.(int))
		os = osAppID

	case "iso_id":
		options.IsoID = isoID.(int)
		os = osIsoID

	//case "snapshot_id":
	//	options.SnapshotID = snapID.(int)
	//	os = osSnapID

	default:
		return errors.New("Error occurred while getting your intended os type")
	}

	regionID := d.Get("region_id").(int)
	planID := d.Get("plan_id").(int)

	// todo we need to loop through network IDs and
	// networkIDs need to handle this separately

	// todo we need to loop through the sshKey and gather those bad boys
	//SSHKeyID: d.Get("") handle this differently

	log.Printf("[INFO] Creating server")
	server, err := client.Server.Create(context.Background(), regionID, planID, os, options)

	if err != nil {
		return fmt.Errorf("Error creating server: %v", err)
	}

	d.SetId(server.VpsID)

	// todo wait for this to be in a "running state"

	// todo call read after this
	return nil
}
func resourceVultrServerUpdate(d *schema.ResourceData, meta interface{}) error { return nil }
func resourceVultrServerRead(d *schema.ResourceData, meta interface{}) error   { return nil }
func resourceVultrServerDelete(d *schema.ResourceData, meta interface{}) error { return nil }

func optionCheck(options map[string]bool) (string, error) {

	result := []string{}
	for k, v := range options {
		if v == true {
			result = append(result, k)
		}
	}

	if len(result) > 1 {
		return "", fmt.Errorf("Too many options have been selected : %v : please select one", result)
	}

	return result[0], nil
}
