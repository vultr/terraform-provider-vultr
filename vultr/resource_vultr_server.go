package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"iso_id": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
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
				ForceNew: true,
				Optional: true,
			},
			"enable_private_network": {
				Type:     schema.TypeBool,
				Computed: true,
				ForceNew: true,
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
				Default:  false,
				Optional: true,
			},
			"app_id": {
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
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
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
				ForceNew: true,
				Optional: true,
			},
			"ddos_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"reserved_ip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
				Optional: true,
			},
			"firewall_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
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
			"server_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v6_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
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
		},
	}
}

func resourceVultrServerCreate(d *schema.ResourceData, meta interface{}) error {

	// Four unique options to image your server
	osID := d.Get("os_id")
	appID, appOK := d.GetOk("app_id")
	isoID, isoOK := d.GetOk("iso_id")
	snapID, snapOK := d.GetOk("snapshot_id")

	osOptions := map[string]bool{"app_id": appOK, "iso_id": isoOK, "snapshot_id": snapOK}
	osOption, err := optionCheck(osOptions)

	if err != nil {
		return err
	}

	client := meta.(*Client).govultrClient()

	options := &govultr.ServerOptions{
		EnableIPV6:           d.Get("enable_ipv6").(bool),
		EnablePrivateNetwork: d.Get("enable_private_network").(bool),
		Label:                d.Get("label").(string),
		AutoBackups:          d.Get("auto_backup").(bool),
		UserData:             d.Get("user_data").(string),
		NotifyActivate:       d.Get("notify_activate").(bool),
		DDOSProtection:       d.Get("ddos_protection").(bool),
		Hostname:             d.Get("hostname").(string),
		Tag:                  d.Get("tag").(string),
		FirewallGroupID:      d.Get("firewall_group_id").(string),
		ScriptID:             d.Get("script_id").(string),
		ReservedIPV4:         d.Get("reserved_ip").(string),
	}

	var os int

	// If no osOptions where selected and osID has a real value then set the osOptions to osID
	if osOption == "" && osID.(int) != 0 {
		osOption = "os_id"
	} else if osOption != "" && osID.(int) != 0 {
		return errors.New(fmt.Sprintf("Please do not set %s with os_id", osOption))
	}

	switch osOption {
	case "os_id":
		os = osID.(int)

	case "app_id":
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

	d.SetId(server.InstanceID)

	_, err = waitForServerAvailable(d, "active", []string{"pending", "installing"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error while waiting for Server %s to be completed: %s", d.Id(), err)
	}

	_, err = waitForServerAvailable(d, "running", []string{"stopped"}, "power_status", meta)
	if err != nil {
		return fmt.Errorf("Error while waiting for Server %s to be in a active state : %s", d.Id(), err)
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
	d.Set("server_state", vps.ServerState)
	d.Set("internal_ip", vps.InternalIP)
	d.Set("kvm_url", vps.KVMUrl)

	if err := d.Set("network_macs", networkMacs); err != nil {
		return fmt.Errorf("Error setting `network_macs`: %#v", err)
	}

	if err := d.Set("network_ips", networkIPs); err != nil {
		return fmt.Errorf("Error setting `network_ips`: %#v", err)
	}

	if err := d.Set("network_ids", networkIDs); err != nil {
		return fmt.Errorf("Error setting `network_ids`: %#v", err)
	}

	var ipv6s []map[string]string
	for _, net := range vps.V6Networks {
		v6network := map[string]string{
			"v6_network":      net.Network,
			"v6_main_ip":      net.MainIP,
			"v6_network_size": net.NetworkSize,
		}
		ipv6s = append(ipv6s, v6network)
	}
	if err := d.Set("v6_networks", ipv6s); err != nil {
		return fmt.Errorf("Error setting `v6_networks`: %#v", err)
	}

	d.Set("tag", vps.Tag)
	d.Set("firewall_group_id", vps.FirewallGroupID)

	regionID, err := strconv.Atoi(vps.RegionID)
	if err != nil {
		return fmt.Errorf("Error while getting regionID for server : %v", err)
	}
	d.Set("region_id", regionID)

	planID, err := strconv.Atoi(vps.PlanID)
	if err != nil {
		return fmt.Errorf("Error while getting planID for server : %v", err)
	}
	d.Set("plan_id", planID)

	osID, err := strconv.Atoi(vps.OsID)
	if err != nil {
		return fmt.Errorf("Error while getting osID for server : %v", err)
	}
	d.Set("os_id", osID)

	appID, err := strconv.Atoi(vps.AppID)
	if err != nil {
		return fmt.Errorf("Error while getting appID for server : %v", err)
	}
	d.Set("app_id", appID)

	if vps.AutoBackups == "yes" {
		d.Set("auto_backup", true)
	} else {
		d.Set("auto_backup", false)
	}

	return nil
}
func resourceVultrServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	if d.HasChange("ddos_protection") {
		log.Printf("[INFO] Updating DDOS Protection")

		_, newVal := d.GetChange("ddos_protection")
		if newVal.(bool) {
			err := client.Server.EnableDDOS(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("error occured while enabling ddos_protection for server %s : %v", d.Id(), err)
			}
		} else {
			err := client.Server.DisableDDOS(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("error occured while disabling ddos_protection for server %s : %v", d.Id(), err)
			}
		}
	}

	if d.HasChange("auto_backup") {
		log.Printf("[INFO] Updating auto backups")

		_, newVal := d.GetChange("auto_backup")
		if newVal.(bool) {
			err := client.Server.EnableBackup(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("Error occured while enabling auto_backup for server %s : %v", d.Id(), err)
			}
		} else {
			err := client.Server.DisableBackup(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("Error occured while disabling auto_backup for server %s : %v", d.Id(), err)
			}
		}
	}

	if d.HasChange("app_id") {
		log.Printf("[INFO] Updating app_id")
		_, newer := d.GetChange("app_id")
		err := client.Server.ChangeApp(context.Background(), d.Id(), strconv.Itoa(newer.(int)))
		if err != nil {
			return fmt.Errorf("Error occured while updating app_id for server %s : %v", d.Id(), err)
		}

		_, err = waitForServerAvailable(d, "active", []string{"pending", "installing"}, "status", meta)
		if err != nil {
			return fmt.Errorf("Error while waiting for Server %s to be in a active state : %s", d.Id(), err)
		}

		_, err = waitForServerAvailable(d, "running", []string{"stopped"}, "power_status", meta)
		if err != nil {
			return fmt.Errorf("Error while waiting for Server %s to be in a active state : %s", d.Id(), err)
		}

	}

	if d.HasChange("os_id") {
		log.Printf("[INFO] Updating os_id")
		_, newer := d.GetChange("os_id")
		err := client.Server.ChangeOS(context.Background(), d.Id(), strconv.Itoa(newer.(int)))
		if err != nil {
			return fmt.Errorf("Error occured while updating os_id for server %s : %v", d.Id(), err)
		}

		_, err = waitForServerAvailable(d, "active", []string{"pending", "installing"}, "status", meta)
		if err != nil {
			return fmt.Errorf("Error while waiting for Server %s to be in a active state : %s", d.Id(), err)
		}

		_, err = waitForServerAvailable(d, "running", []string{"stopped"}, "power_status", meta)
		if err != nil {
			return fmt.Errorf("Error while waiting for Server %s to be in a active state : %s", d.Id(), err)
		}

	}

	if d.HasChange("user_data") {
		log.Printf("[INFO] Updating user_data")
		err := client.Server.SetUserData(context.Background(), d.Id(), d.Get("user_data").(string))
		if err != nil {
			return fmt.Errorf("Error occured while updating user_data for server %s : %v", d.Id(), err)
		}
	}

	if d.HasChange("firewall_group_id") {
		log.Printf("[INFO] Updating firewall_group_id")
		err := client.Server.SetFirewallGroup(context.Background(), d.Id(), d.Get("firewall_group_id").(string))
		if err != nil {
			return fmt.Errorf("Error occured while updating firewall_group_id for server %s : %v", d.Id(), err)
		}
	}

	if d.HasChange("tag") {
		log.Printf("[INFO] Updating tag")
		err := client.Server.SetTag(context.Background(), d.Id(), d.Get("tag").(string))
		if err != nil {
			return fmt.Errorf("Error occured while updating tag for server %s : %v", d.Id(), err)
		}
	}

	if d.HasChange("label") {
		log.Printf("[INFO] Updating label")
		err := client.Server.SetLabel(context.Background(), d.Id(), d.Get("label").(string))
		if err != nil {
			return fmt.Errorf("Error occured while updating label for server %s : %v", d.Id(), err)
		}
	}

	if d.HasChange("network_ids") {
		log.Printf("[INFO] Updating network_ids")
		oldNetwork, newNetwork := d.GetChange("network_ids")

		var oldIDs []string
		for _, v := range oldNetwork.([]interface{}) {
			oldIDs = append(oldIDs, v.(string))
		}

		var newIDs []string
		for _, v := range newNetwork.([]interface{}) {
			newIDs = append(newIDs, v.(string))
		}

		diff := func(in, out []string) []string {
			var diff []string

			b := map[string]string{}
			for i := range in {
				b[in[i]] = ""
			}

			for i := range out {
				if _, ok := b[out[i]]; !ok {
					diff = append(diff, out[i])
				}
			}

			return diff
		}

		for _, v := range diff(oldIDs, newIDs) {
			err := client.Server.EnablePrivateNetwork(context.Background(), d.Id(), v)

			if err != nil {
				return fmt.Errorf("Error attaching network id %s to server %s : %v", v, d.Id(), err)
			}
		}

		for _, v := range diff(newIDs, oldIDs) {
			err := client.Server.DisablePrivateNetwork(context.Background(), d.Id(), v)

			if err != nil {
				return fmt.Errorf("Error detaching network id %s from server %s : %v", v, d.Id(), err)
			}
		}

	}

	return resourceVultrServerRead(d, meta)
}

func resourceVultrServerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting instance (%s)", d.Id())

	ids, err := client.Server.ListPrivateNetworks(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error grabbing private networks associated to server %s : %v", d.Id(), err)
	}

	for i := range ids {
		err := client.Server.DisablePrivateNetwork(context.Background(), d.Id(), ids[i].NetworkID)

		if err != nil {
			return fmt.Errorf("Error detaching network id %s from server %s : %v", ids[i].NetworkID, d.Id(), err)
		}
	}

	isoID, ok := d.GetOk("os_id")
	if ok {
		err = client.Server.IsoDetach(context.Background(), d.Id())
		if err != nil {
			return fmt.Errorf("error detaching iso (%s) from instance (%s) : %v", d.Id(), isoID, err)
		}
	}

	err = client.Server.Delete(context.Background(), d.Id())

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

	// Return back an empty slice so we can possibly add in osID
	if len(result) == 0 {
		return "", nil
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
		Refresh:        newServerStateRefresh(d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForState()
}

func newServerStateRefresh(d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating Server")
		server, err := client.Server.GetServer(context.Background(), d.Id())

		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving Server %s : %s", d.Id(), err)
		}

		if attr == "status" {
			log.Printf("[INFO] The Server Status is %s", server.Status)
			return server, server.Status, nil
		} else if attr == "power_status" {
			log.Printf("[INFO] The Server Power Status is %s", server.PowerStatus)
			return server, server.PowerStatus, nil
		} else {
			return nil, "", nil
		}
	}
}
