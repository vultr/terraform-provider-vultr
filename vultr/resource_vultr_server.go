package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/vultr/govultr"
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
			"network_macs": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"network_ips": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	snapID, snapOK := d.GetOk("snapshot_id")

	osOptions := map[string]bool{"os_id": osOK, "application_id": appOK, "iso_id": isoOK, "snapshot_id": snapOK}
	osOption, err := optionCheck(osOptions)

	if err != nil {
		return err
	}

	client := meta.(*Client).govultrClient()

	options := &govultr.ServerOptions{
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

	case "snapshot_id":
		options.SnapshotID = snapID.(string)
		os = osSnapID

	default:
		return errors.New("Error occurred while getting your intended os type")
	}

	regionID := d.Get("region_id").(int)
	planID := d.Get("plan_id").(int)

	networkIDs, networkOK := d.GetOk("network_ids")
	if networkOK {
		for _, v := range networkIDs.([]interface{}) {
			options.NetworkID = append(options.NetworkID, v.(string))
		}
	}

	sshKeyIDs, sshKeyOK := d.GetOk("ssh_key_ids")
	if sshKeyOK {
		log.Print(sshKeyOK)
		for _, v := range sshKeyIDs.([]interface{}) {
			options.SSHKeyIDs = append(options.SSHKeyIDs, v.(string))
		}
	}

	log.Printf("[INFO] Creating server")
	server, err := client.Server.Create(context.Background(), regionID, planID, os, options)

	if err != nil {
		return fmt.Errorf("Error creating server: %v", err)
	}

	d.SetId(server.VpsID)

	_, err = waitForServerAvailable(d, "active", []string{"pending", "installing"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error while waiting for Server %s to be completed: %s", d.Id(), err)
	}

	return resourceVultrServerRead(d, meta)
}

func resourceVultrServerRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()

	vps, err := client.Server.GetServer(context.Background(), d.Id())

	if err != nil {
		if strings.HasPrefix(err.Error(), "Invalid server") {
			log.Printf("[WARN] Removing instance (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting instance (%s): %v", d.Id(), err)
	}

	networks, err := client.Server.ListPrivateNetworks(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error getting private networks for server  %s : %v", d.Id(), err)
	}

	var networkIDs []string
	networkIPs := make(map[string]string)
	networkMacs := make(map[string]string)
	for _, v := range networks {
		networkMacs[v.NetworkID] = v.MacAddress
		networkIPs[v.NetworkID] = v.IPAddress
		networkIDs = append(networkIDs, v.NetworkID)
	}

	d.Set("os", vps.Os)
	d.Set("ram", vps.RAM)
	d.Set("disk", vps.Disk)
	d.Set("main_ip", vps.MainIP)
	d.Set("vps_cpu_count", vps.VPSCpus)
	d.Set("location", vps.Location)
	d.Set("default_password", vps.DefaultPassword)
	d.Set("date_created", vps.Created)
	d.Set("pending_charges", vps.PendingCharges)
	d.Set("status", vps.Status)
	d.Set("cost_per_month", vps.Cost)
	d.Set("current_bandwidth", vps.CurrentBandwidth)
	d.Set("allowed_bandwidth", vps.AllowedBandwidth)
	d.Set("netmask_v4", vps.NetmaskV4)
	d.Set("gateway_v4", vps.GatewayV4)
	d.Set("power_status", vps.PowerStatus)
	d.Set("server_status", vps.ServerState)
	d.Set("internal_ip", vps.InternalIP)
	d.Set("kvm_url", vps.KVMUrl)
	d.Set("network_macs", networkMacs)
	d.Set("network_ips", networkIPs)

	d.Set("network_ids", networks)
	var ipv6s []string
	for _, net := range vps.V6Networks {
		ipv6s = append(ipv6s, net.MainIP)
	}
	d.Set("v6_networks", ipv6s)

	d.Set("tag", vps.Tag)
	d.Set("region_id", vps.RegionID)
	d.Set("firewall_group_id", vps.FirewallGroupID)

	return nil
}
func resourceVultrServerUpdate(d *schema.ResourceData, meta interface{}) error { return nil }

func resourceVultrServerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Destroying instance (%s)", d.Id())
	err := client.Server.Destroy(context.Background(), d.Id())

	if err != nil {
		return fmt.Errorf("Error destroying instance %s : %v", d.Id(), err)
	}

	return nil
}

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

func waitForServerAvailable(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for Server (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newServerStateRefresh(d, meta),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForState()
}

func newServerStateRefresh(
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()

	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating Server")
		server, err := client.Server.GetServer(context.Background(), d.Id())

		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving Server %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] The Server Status is %s", server.Status)
		return server, server.Status, nil
	}
}
