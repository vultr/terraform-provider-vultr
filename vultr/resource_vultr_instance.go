package vultr

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v2"
)

func resourceVultrInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrInstanceCreate,
		ReadContext:   resourceVultrInstanceRead,
		UpdateContext: resourceVultrInstanceUpdate,
		DeleteContext: resourceVultrInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			//Required
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
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
			"image_id": {
				Type:     schema.TypeString,
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
			},
			"enable_private_network": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    false,
				Deprecated: "In the next release of this provider we will be removing `enable_private_network` due to issues that may cause drift and having to maintain private network ip state. Please switch to using private_network_ids to manage your private network fields.",
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
				ForceNew: true,
			},
			"activation_email": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ddos_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Changing the hostname after initial deployment will trigger a reinstall",
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"reserved_ip_id": {
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
			"backups": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "disabled",
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"backups_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"daily", "weekly", "monthly", "daily_alt_even", "daily_alt_odd"}, false),
						},
						"hour": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"dow": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"dom": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceVultrInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// five unique options to image your server
	osID := d.Get("os_id")
	appID, appOK := d.GetOk("app_id")
	imageID, imageOK := d.GetOk("image_id")
	isoID, isoOK := d.GetOk("iso_id")
	snapID, snapOK := d.GetOk("snapshot_id")
	backups := d.Get("backups").(string)
	backupSchedule, backupsScheduleOk := d.GetOk("backups_schedule")

	if backups == "enabled" && !backupsScheduleOk {
		return diag.Errorf("Backups are set to enabled please provide a backups_schedule")
	} else if backups == "disabled" && backupsScheduleOk {
		return diag.Errorf("Backups are set to disabled please remove backups_schedule")
	}

	osOptions := map[string]bool{"app_id": appOK, "iso_id": isoOK, "snapshot_id": snapOK, "image_id": imageOK}
	osOption, err := optionCheck(osOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*Client).govultrClient()

	req := &govultr.InstanceCreateReq{
		EnableIPv6:           govultr.BoolToBoolPtr(d.Get("enable_ipv6").(bool)),
		EnablePrivateNetwork: govultr.BoolToBoolPtr(d.Get("enable_private_network").(bool)),
		Label:                d.Get("label").(string),
		Backups:              backups,
		UserData:             base64.StdEncoding.EncodeToString([]byte(d.Get("user_data").(string))),
		ActivationEmail:      govultr.BoolToBoolPtr(d.Get("activation_email").(bool)),
		DDOSProtection:       govultr.BoolToBoolPtr(d.Get("ddos_protection").(bool)),
		Hostname:             d.Get("hostname").(string),
		Tag:                  d.Get("tag").(string),
		FirewallGroupID:      d.Get("firewall_group_id").(string),
		ScriptID:             d.Get("script_id").(string),
		ReservedIPv4:         d.Get("reserved_ip_id").(string),
		Region:               d.Get("region").(string),
		Plan:                 d.Get("plan").(string),
	}

	// If no osOptions where selected and osID has a real value then set the osOptions to osID
	if osOption == "" && osID.(int) != 0 {
		osOption = "os_id"
	} else if osOption != "" && osID.(int) != 0 {
		return diag.Errorf(fmt.Sprintf("Please do not set %s with os_id", osOption))
	}

	switch osOption {
	case "os_id":
		req.OsID = osID.(int)
	case "app_id":
		req.AppID = appID.(int)
	case "image_id":
		req.ImageID = imageID.(string)
	case "iso_id":
		req.ISOID = isoID.(string)
	case "snapshot_id":
		req.SnapshotID = snapID.(string)
	default:
		return diag.Errorf("error occurred while getting your intended os type")
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
	instance, err := client.Instance.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating server: %v", err)
	}

	d.SetId(instance.ID)
	d.Set("default_password", instance.DefaultPassword)

	if _, err = waitForServerAvailable(ctx, d, "active", []string{"pending", "installing"}, "status", meta); err != nil {
		return diag.Errorf("error while waiting for Server %s to be completed: %s", d.Id(), err)
	}

	if _, err = waitForServerAvailable(ctx, d, "running", []string{"stopped"}, "power_status", meta); err != nil {
		return diag.Errorf("error while waiting for Server %s to be in a active state : %s", d.Id(), err)
	}

	if backups == "enabled" {
		backupReq := generateBackupSchedule(backupSchedule)
		if err := client.Instance.SetBackupSchedule(context.Background(), instance.ID, backupReq); err != nil {
			return diag.Errorf("error setting backup schedule: %v", err)
		}
	}

	return resourceVultrInstanceRead(ctx, d, meta)
}

func resourceVultrInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instance, err := client.Instance.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "invalid instance ID") {
			log.Printf("[WARN] Removing instance (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting instance (%s): %v", d.Id(), err)
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
	d.Set("hostname", instance.Hostname)

	backup, err := client.Instance.GetBackupSchedule(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting backup schedule: %v", err)
	}

	d.Set("backups", backupStatus(backup.Enabled))

	if backupStatus(backup.Enabled) != "disabled" {
		var bs []map[string]interface{}
		backupScheduleInfo := map[string]interface{}{
			"type": backup.Type,
			"hour": backup.Hour,
			"dom":  backup.Dom,
			"dow":  backup.Dow,
		}
		bs = append(bs, backupScheduleInfo)

		if err := d.Set("backups_schedule", bs); err != nil {
			return diag.Errorf("error setting `backups_schedule`: %v", err)
		}
	} else {
		d.Set("backups_schedule", nil)
	}

	pn, err := getPrivateNetworks(client, d.Id())
	if err != nil {
		return diag.Errorf(err.Error())
	}

	d.Set("private_network_ids", pn)

	return nil
}
func resourceVultrInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.InstanceUpdateReq{
		Label:           d.Get("label").(string),
		Tag:             d.Get("tag").(string),
		FirewallGroupID: d.Get("firewall_group_id").(string),
		EnableIPv6:      govultr.BoolToBoolPtr(d.Get("enable_ipv6").(bool)),
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

	if d.HasChange("enable_private_network") {
		log.Printf("[INFO] Updating private networking")
		_, newVal := d.GetChange("enable_private_network")
		req.EnablePrivateNetwork = govultr.BoolToBoolPtr(newVal.(bool))
	}

	bs, bsOK := d.GetOk("backups_schedule")
	_, newBackupValue := d.GetChange("backups")
	if d.HasChange("backups") {
		log.Printf("[INFO] Updating Backups")
		backups := newBackupValue.(string)
		req.Backups = backups

		if backups == "disabled" && bsOK {
			return diag.Errorf("Backups are being set to disabled please remove backups_schedule")
		} else if backups == "enabled" && !bsOK {
			return diag.Errorf("Backups are being set to enabled please add backups_schedule")
		}
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

	if _, err := client.Instance.Update(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating instance %s : %s", d.Id(), err.Error())
	}

	if d.HasChange("iso_id") {
		log.Printf("[INFO] Updating ISO")

		_, newISOId := d.GetChange("iso_id")
		if newISOId == "" {
			if err := client.Instance.DetachISO(ctx, d.Id()); err != nil {
				return diag.Errorf("error detaching iso from instance %s : %v", d.Id(), err)
			}
		} else {
			if err := client.Instance.AttachISO(ctx, d.Id(), newISOId.(string)); err != nil {
				return diag.Errorf("error attaching iso to instance %s : %v", d.Id(), err)
			}
		}
	}

	if newBackupValue.(string) == "enabled" && !bsOK {
		return diag.Errorf("Backups are being set to enabled please add backups_schedule")
	}

	// If we are disabling backups we don't do anything.
	// On the read that gets called we will nil out backups_schedule.
	if newBackupValue.(string) != "disabled" && d.HasChange("backups_schedule") {
		schedule := generateBackupSchedule(bs)
		if err := client.Instance.SetBackupSchedule(ctx, d.Id(), schedule); err != nil {
			return diag.Errorf("error setting backup for %s : %v", d.Id(), err)
		}
	}

	// Changing the hostname can only be done via a reinstall
	// Since this is a full reinstall we also put in the status waits so TF doesn't continue processing until the instance is fully up
	if d.HasChange("hostname") {
		req := &govultr.ReinstallReq{Hostname: d.Get("hostname").(string)}
		if _, err := client.Instance.Reinstall(ctx, d.Id(), req); err != nil {
			return diag.Errorf("error changing hostname for %s : %v", d.Id(), err)
		}

		if _, err := waitForServerAvailable(ctx, d, "active", []string{"pending", "installing"}, "status", meta); err != nil {
			return diag.Errorf("error while waiting for Server %s to be completed: %s", d.Id(), err)
		}

		if _, err := waitForServerAvailable(ctx, d, "running", []string{"stopped"}, "power_status", meta); err != nil {
			return diag.Errorf("error while waiting for Server %s to be in a active state : %s", d.Id(), err)
		}
	}

	// There is a delay between the API data returning the newly updated plan change
	// This will wait until the plan has been updated before going to the read call
	if d.HasChange("plan") {
		oldP, newP := d.GetChange("plan")
		if _, err := waitForUpgrade(ctx, d, newP.(string), []string{oldP.(string)}, "plan", meta); err != nil {
			return diag.Errorf("error while waiting for instance %s to have updated plan : %s", d.Id(), err)
		}
	}

	return resourceVultrInstanceRead(ctx, d, meta)
}

func resourceVultrInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting instance (%s)", d.Id())

	if networkIDs, networkOK := d.GetOk("private_network_ids"); networkOK {
		detach := &govultr.InstanceUpdateReq{}
		for _, v := range networkIDs.([]interface{}) {
			detach.DetachPrivateNetwork = append(detach.DetachPrivateNetwork, v.(string))
		}

		if _, err := client.Instance.Update(ctx, d.Id(), detach); err != nil {
			return diag.Errorf("error detaching private networks prior to deleting instance %s : %v", d.Id(), err)
		}
	}

	if err := client.Instance.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying instance %s : %v", d.Id(), err)
	}

	return nil
}

func optionCheck(options map[string]bool) (string, error) {

	var result []string
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

func waitForServerAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for Server (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newServerStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newServerStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating Server")
		server, err := client.Instance.Get(ctx, d.Id())
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

func waitForUpgrade(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for instance (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newInstancePlanRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newInstancePlanRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Upgrading instance")
		instance, err := client.Instance.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving instance %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] The instances plan is %s", instance.Plan)
		return instance, instance.Plan, nil
	}
}

func generateBackupSchedule(backup interface{}) *govultr.BackupScheduleReq {
	k := backup.([]interface{})

	config := k[0].(map[string]interface{})
	return &govultr.BackupScheduleReq{
		Type: config["type"].(string),
		Hour: govultr.IntToIntPtr(config["hour"].(int)),
		Dom:  config["dom"].(int),
		Dow:  govultr.IntToIntPtr(config["dow"].(int)),
	}
}

func backupStatus(status *bool) string {
	if *status {
		return "enabled"
	} else {
		return "disabled"
	}
}
