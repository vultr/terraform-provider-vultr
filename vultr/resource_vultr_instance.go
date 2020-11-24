package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrInstanceCreate,
		Read:   resourceVultrInstanceRead,
		Update: resourceVultrInstanceUpdate,
		Delete: resourceVultrInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			//Required
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"iso_id": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"os_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"script_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_private_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"private_network_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Default:  nil,
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
			"backups": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "disabled",
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"user_data": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringIsBase64,
			},
			"activation_email": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  true,
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
				Type:         schema.TypeString,
				ForceNew:     true, // force new?
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"firewall_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// Computed
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
			"default_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
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
			"power_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_status": {
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
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kvm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func resourceVultrInstanceCreate(d *schema.ResourceData, meta interface{}) error {
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

	req := &govultr.InstanceCreateReq{
		EnableIPv6:           d.Get("enable_ipv6").(bool),
		EnablePrivateNetwork: d.Get("enable_private_network").(bool),
		Label:                d.Get("label").(string),
		Backups:              d.Get("backups").(string),
		UserData:             d.Get("user_data").(string),
		ActivationEmail:      d.Get("activation_email").(bool),
		DDOSProtection:       d.Get("ddos_protection").(bool),
		Hostname:             d.Get("hostname").(string),
		Tag:                  d.Get("tag").(string),
		FirewallGroupID:      d.Get("firewall_group_id").(string),
		ScriptID:             d.Get("script_id").(string),
		ReservedIPv4:         d.Get("reserved_ip").(string),
		Region:               d.Get("region").(string),
		Plan:                 d.Get("plan").(string),
	}

	// If no osOptions where selected and osID has a real value then set the osOptions to osID
	if osOption == "" && osID.(int) != 0 {
		osOption = "os_id"
	} else if osOption != "" && osID.(int) != 0 {
		return fmt.Errorf(fmt.Sprintf("Please do not set %s with os_id", osOption))
	}

	switch osOption {
	case "os_id":
		req.OsID = osID.(int)
	case "app_id":
		req.AppID = appID.(int)
	case "iso_id":
		req.ISOID = isoID.(string)
	case "snapshot_id":
		req.SnapshotID = snapID.(string)
	default:
		return fmt.Errorf("error occurred while getting your intended os type")
	}

	if networkIDs, networkOK := d.GetOk("private_network_ids"); networkOK {
		for _, v := range networkIDs.([]interface{}) {
			req.AttachPrivateNetwork = append(req.AttachPrivateNetwork, v.(string))
		}
	}

	if sshKeyIDs, sshKeyOK := d.GetOk("ssh_key_ids"); sshKeyOK {
		for _, v := range sshKeyIDs.([]interface{}) {
			req.SSHKeys = append(req.SSHKeys, v.(string))
		}
	}

	log.Printf("[INFO] Creating server")
	instance, err := client.Instance.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("error creating server: %v", err)
	}

	d.SetId(instance.ID)
	d.Set("default_password", instance.DefaultPassword)

	if _, err = waitForServerAvailable(d, "active", []string{"pending", "installing"}, "status", meta); err != nil {
		return fmt.Errorf("error while waiting for Server %s to be completed: %s", d.Id(), err)
	}

	if _, err = waitForServerAvailable(d, "running", []string{"stopped"}, "power_status", meta); err != nil {
		return fmt.Errorf("error while waiting for Server %s to be in a active state : %s", d.Id(), err)
	}

	return resourceVultrInstanceRead(d, meta)
}

func resourceVultrInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	instance, err := client.Instance.Get(context.Background(), d.Id())
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid server") {
			log.Printf("[WARN] Removing instance (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error getting instance (%s): %v", d.Id(), err)
	}

	d.Set("os", instance.Os)
	d.Set("ram", instance.RAM)
	d.Set("disk", instance.Disk)
	d.Set("main_ip", instance.MainIP)
	d.Set("vcpu_count", instance.VCPUCount)
	d.Set("date_created", instance.DateCreated)
	d.Set("status", instance.Status)
	d.Set("allowed_bandwidth", instance.AllowedBandwidth)
	d.Set("netmask_v4", instance.NetmaskV4)
	d.Set("gateway_v4", instance.GatewayV4)
	d.Set("power_status", instance.PowerStatus)
	d.Set("server_status", instance.ServerStatus)
	d.Set("internal_ip", instance.InternalIP)
	d.Set("kvm", instance.KVM)
	d.Set("v6_network", instance.V6Network)
	d.Set("v6_main_ip", instance.V6MainIP)
	d.Set("v6_network_size", instance.V6NetworkSize)
	d.Set("tag", instance.Tag)
	d.Set("firewall_group_id", instance.FirewallGroupID)
	d.Set("region", instance.Region)
	d.Set("plan", instance.Plan)
	d.Set("os_id", instance.OsID)
	d.Set("app_id", instance.AppID)
	d.Set("features", instance.Features)

	return nil
}
func resourceVultrInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	req := &govultr.InstanceUpdateReq{
		Label:                d.Get("label").(string),
		Tag:                  d.Get("tag").(string),
		FirewallGroupID:      d.Get("firewall_group_id").(string),
		EnableIPv6:           d.Get("enable_ipv6").(bool),
		EnablePrivateNetwork: d.Get("enable_private_network").(bool),
		UserData:             d.Get("user_data").(string),
	}

	if d.HasChange("plan") {
		log.Printf("[INFO] Updating Plan")
		_, newVal := d.GetChange("plan")
		plan := newVal.(string)
		req.Plan = plan
	}

	if d.HasChange("ddos_protection") {
		log.Printf("[INFO] Updating DDOS Protection")
		_, newVal := d.GetChange("ddos_protection")
		ddos := newVal.(bool)
		req.DDOSProtection = &ddos
	}

	if d.HasChange("backups") {
		log.Printf("[INFO] Updating Backups")
		_, newVal := d.GetChange("backups")
		backups := newVal.(string)
		req.Backups = backups
	}

	if d.HasChange("private_network_ids") {
		log.Printf("[INFO] Updating private_network_ids")
		oldNetwork, newNetwork := d.GetChange("private_network_ids")

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
			req.AttachPrivateNetwork = append(req.AttachPrivateNetwork, v)
		}

		for _, v := range diff(newIDs, oldIDs) {
			req.DetachPrivateNetwork = append(req.DetachPrivateNetwork, v)
		}

	}

	if err := client.Instance.Update(context.Background(), d.Id(), req); err != nil {
		return fmt.Errorf("error updating instance %s : %s", d.Id(), err.Error())
	}

	if d.HasChange("iso_id") {
		log.Printf("[INFO] Updating ISO")

		_, newISOId := d.GetChange("iso_id")
		if newISOId == "" {
			if err := client.Instance.DetachISO(context.Background(), d.Id()); err != nil {
				return fmt.Errorf("error detaching iso from instance %s : %v", d.Id(), err)
			}
		} else {
			if err := client.Instance.AttachISO(context.Background(), d.Id(), newISOId.(string)); err != nil {
				return fmt.Errorf("error detaching iso from instance %s : %v", d.Id(), err)
			}
		}
	}

	return resourceVultrInstanceRead(d, meta)
}

func resourceVultrInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting instance (%s)", d.Id())

	if err := client.Instance.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error destroying instance %s : %v", d.Id(), err)
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
		return "", fmt.Errorf("too many options have been selected : %v : please select one", result)
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
		server, err := client.Instance.Get(context.Background(), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Server %s : %s", d.Id(), err)
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
