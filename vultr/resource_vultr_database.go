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
	"github.com/vultr/govultr/v3"
)

const defaultTimeout = 60 * time.Minute

func resourceVultrDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDatabaseCreate,
		ReadContext:   resourceVultrDatabaseRead,
		UpdateContext: resourceVultrDatabaseUpdate,
		DeleteContext: resourceVultrDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
			"region": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: IgnoreCase,
			},
			"plan": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
				Required: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"database_engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_engine_version": {
				Type:     schema.TypeString,
				Computed: true,
				Required: true,
			},
			"maintenance_dow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maintenance_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trusted_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mysql_sql_modes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mysql_slow_query_log": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mysql_long_query_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"redis_eviction_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan_raw": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan_vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan_replicas": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_backup": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceVultrDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.InstanceCreateReq{
		EnableIPv6:      govultr.BoolToBoolPtr(d.Get("enable_ipv6").(bool)),
		Label:           d.Get("label").(string),
		UserData:        base64.StdEncoding.EncodeToString([]byte(d.Get("user_data").(string))),
		ActivationEmail: govultr.BoolToBoolPtr(d.Get("activation_email").(bool)),
		DDOSProtection:  govultr.BoolToBoolPtr(d.Get("ddos_protection").(bool)),
		Hostname:        d.Get("hostname").(string),
		FirewallGroupID: d.Get("firewall_group_id").(string),
		ScriptID:        d.Get("script_id").(string),
		ReservedIPv4:    d.Get("reserved_ip_id").(string),
		Region:          d.Get("region").(string),
		Plan:            d.Get("plan").(string),
	}

	if tagsIDs, tagsOK := d.GetOk("tags"); tagsOK {
		for _, v := range tagsIDs.(*schema.Set).List() {
			req.Tags = append(req.Tags, v.(string))
		}
	}

	if len(d.Get("private_network_ids").(*schema.Set).List()) != 0 && len(d.Get("vpc_ids").(*schema.Set).List()) != 0 {
		return diag.Errorf("private_network_ids cannot be used along with vpc_ids. Use only vpc_ids instead.")
	}

	if networkIDs, networkOK := d.GetOk("private_network_ids"); networkOK {
		for _, v := range networkIDs.(*schema.Set).List() {
			req.AttachVPC = append(req.AttachVPC, v.(string))
		}
	}

	if vpcIDs, vpcOK := d.GetOk("vpc_ids"); vpcOK {
		for _, v := range vpcIDs.(*schema.Set).List() {
			req.AttachVPC = append(req.AttachVPC, v.(string))
		}
	}

	if sshKeyIDs, sshKeyOK := d.GetOk("ssh_key_ids"); sshKeyOK {
		for _, v := range sshKeyIDs.([]interface{}) {
			req.SSHKeys = append(req.SSHKeys, v.(string))
		}
	}

	log.Printf("[INFO] Creating server")
	instance, _, err := client.Instance.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating server: %v", err)
	}

	d.SetId(instance.ID)
	if err := d.Set("default_password", instance.DefaultPassword); err != nil {
		return diag.Errorf("unable to set resource instance `default_password` create value: %v", err)
	}

	if _, err = waitForDatabaseAvailable(ctx, d, "active", []string{"pending", "installing"}, "status", meta); err != nil {
		return diag.Errorf("error while waiting for Server %s to be completed: %s", d.Id(), err)
	}

	if _, err = waitForDatabaseAvailable(ctx, d, "running", []string{"stopped"}, "power_status", meta); err != nil {
		return diag.Errorf("error while waiting for Server %s to be in a active state : %s", d.Id(), err)
	}

	return resourceVultrInstanceRead(ctx, d, meta)
}

func resourceVultrDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	instance, _, err := client.Instance.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "invalid instance ID") {
			log.Printf("[WARN] Removing instance (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting instance (%s): %v", d.Id(), err)
	}

	if err := d.Set("os", instance.Os); err != nil {
		return diag.Errorf("unable to set resource instance `os` read value: %v", err)
	}
	if err := d.Set("ram", instance.RAM); err != nil {
		return diag.Errorf("unable to set resource instance `ram` read value: %v", err)
	}
	if err := d.Set("disk", instance.Disk); err != nil {
		return diag.Errorf("unable to set resource instance `disk` read value: %v", err)
	}
	if err := d.Set("main_ip", instance.MainIP); err != nil {
		return diag.Errorf("unable to set resource instance `main_ip` read value: %v", err)
	}
	if err := d.Set("vcpu_count", instance.VCPUCount); err != nil {
		return diag.Errorf("unable to set resource instance `vcpu_count` read value: %v", err)
	}
	if err := d.Set("date_created", instance.DateCreated); err != nil {
		return diag.Errorf("unable to set resource instance `date_created` read value: %v", err)
	}
	if err := d.Set("status", instance.Status); err != nil {
		return diag.Errorf("unable to set resource instance `status` read value: %v", err)
	}
	if err := d.Set("allowed_bandwidth", instance.AllowedBandwidth); err != nil {
		return diag.Errorf("unable to set resource instance `allowed_bandwidth` read value: %v", err)
	}
	if err := d.Set("netmask_v4", instance.NetmaskV4); err != nil {
		return diag.Errorf("unable to set resource instance `netmask_v4` read value: %v", err)
	}
	if err := d.Set("gateway_v4", instance.GatewayV4); err != nil {
		return diag.Errorf("unable to set resource instance `gateway_v4` read value: %v", err)
	}
	if err := d.Set("power_status", instance.PowerStatus); err != nil {
		return diag.Errorf("unable to set resource instance `power_status` read value: %v", err)
	}
	if err := d.Set("server_status", instance.ServerStatus); err != nil {
		return diag.Errorf("unable to set resource instance `server_status` read value: %v", err)
	}
	if err := d.Set("internal_ip", instance.InternalIP); err != nil {
		return diag.Errorf("unable to set resource instance `internal_ip` read value: %v", err)
	}
	if err := d.Set("kvm", instance.KVM); err != nil {
		return diag.Errorf("unable to set resource instance `kvm` read value: %v", err)
	}
	if err := d.Set("v6_network", instance.V6Network); err != nil {
		return diag.Errorf("unable to set resource instance `v6_network` read value: %v", err)
	}
	if err := d.Set("v6_main_ip", instance.V6MainIP); err != nil {
		return diag.Errorf("unable to set resource instance `v6_main_ip` read value: %v", err)
	}
	if err := d.Set("v6_network_size", instance.V6NetworkSize); err != nil {
		return diag.Errorf("unable to set resource instance `v6_network_size` read value: %v", err)
	}
	if err := d.Set("tags", instance.Tags); err != nil {
		return diag.Errorf("unable to set resource instance `tags` read value: %v", err)
	}
	if err := d.Set("firewall_group_id", instance.FirewallGroupID); err != nil {
		return diag.Errorf("unable to set resource instance `firewall_group_id` read value: %v", err)
	}
	if err := d.Set("region", instance.Region); err != nil {
		return diag.Errorf("unable to set resource instance `region` read value: %v", err)
	}
	if err := d.Set("plan", instance.Plan); err != nil {
		return diag.Errorf("unable to set resource instance `plan` read value: %v", err)
	}
	if err := d.Set("os_id", instance.OsID); err != nil {
		return diag.Errorf("unable to set resource instance `os_id` read value: %v", err)
	}
	if err := d.Set("app_id", instance.AppID); err != nil {
		return diag.Errorf("unable to set resource instance `app_id` read value: %v", err)
	}
	if err := d.Set("features", instance.Features); err != nil {
		return diag.Errorf("unable to set resource instance `features` read value: %v", err)
	}
	if err := d.Set("hostname", instance.Hostname); err != nil {
		return diag.Errorf("unable to set resource instance `hostname` read value: %v", err)
	}

	vpcs, err := getVPCs(client, d.Id())
	if err != nil {
		return diag.Errorf(err.Error())
	}

	// Manipulate the read state so that, depending on which value was passed,
	// only one of these values is populated when a VPC or PN is defined for
	// the instance
	if _, pnUpdate := d.GetOk("private_network_ids"); pnUpdate {
		if err := d.Set("private_network_ids", vpcs); err != nil {
			return diag.Errorf("unable to set resource instance `private_network_ids` read value: %v", err)
		}
		if err := d.Set("vpc_ids", nil); err != nil {
			return diag.Errorf("unable to set resource instance `vpc_ids` read value: %v", err)
		}
	}

	// Since VPC is last, if an instance read involves both vpc_ids &
	// private_network_ids, only the vpc_ids will be preserved
	if _, vpcUpdate := d.GetOk("vpc_ids"); vpcUpdate {
		if err := d.Set("vpc_ids", vpcs); err != nil {
			return diag.Errorf("unable to set resource instance `vpc_ids` read value: %v", err)
		}
		if err := d.Set("private_network_ids", nil); err != nil {
			return diag.Errorf("unable to set resource instance `private_network_ids` read value: %v", err)
		}
	}

	return nil
}
func resourceVultrDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.InstanceUpdateReq{
		Label:           d.Get("label").(string),
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

	if len(d.Get("private_network_ids").(*schema.Set).List()) != 0 && len(d.Get("vpc_ids").(*schema.Set).List()) != 0 {
		return diag.Errorf("private_network_ids cannot be used along with vpc_ids. Use only vpc_ids instead.")
	}

	if d.HasChange("private_network_ids") {
		log.Printf("[INFO] Updating private_network_ids")
		oldNetwork, newNetwork := d.GetChange("private_network_ids")

		var oldIDs []string
		for _, v := range oldNetwork.(*schema.Set).List() {
			oldIDs = append(oldIDs, v.(string))
		}

		var newIDs []string
		for _, v := range newNetwork.(*schema.Set).List() {
			newIDs = append(newIDs, v.(string))
		}

		req.AttachPrivateNetwork = append(req.AttachPrivateNetwork, diffSlice(oldIDs, newIDs)...) // nolint
		req.DetachPrivateNetwork = append(req.DetachPrivateNetwork, diffSlice(newIDs, oldIDs)...) // nolint
	}

	if d.HasChange("vpc_ids") {
		log.Printf("[INFO] Updating vpc_ids")
		oldVPC, newVPC := d.GetChange("vpc_ids")

		var oldIDs []string
		for _, v := range oldVPC.(*schema.Set).List() {
			oldIDs = append(oldIDs, v.(string))
		}

		var newIDs []string
		for _, v := range newVPC.(*schema.Set).List() {
			newIDs = append(newIDs, v.(string))
		}

		req.AttachVPC = append(req.AttachVPC, diffSlice(oldIDs, newIDs)...)
		req.DetachVPC = append(req.DetachVPC, diffSlice(newIDs, oldIDs)...)
	}

	if d.HasChange("tags") {
		_, newTags := tfChangeToSlices("tags", d)
		req.Tags = newTags
	}

	if _, _, err := client.Instance.Update(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating instance %s : %s", d.Id(), err.Error())
	}

	if d.HasChange("iso_id") {
		log.Printf("[INFO] Updating ISO")

		_, newISOId := d.GetChange("iso_id")
		if newISOId == "" {
			if _, err := client.Instance.DetachISO(ctx, d.Id()); err != nil {
				return diag.Errorf("error detaching iso from instance %s : %v", d.Id(), err)
			}
		} else {
			if _, err := client.Instance.AttachISO(ctx, d.Id(), newISOId.(string)); err != nil {
				return diag.Errorf("error attaching iso to instance %s : %v", d.Id(), err)
			}
		}
	}

	// There is a delay between the API data returning the newly updated plan change
	// This will wait until the plan has been updated before going to the read call
	if d.HasChange("plan") {
		oldP, newP := d.GetChange("plan")
		if _, err := waitForDatabaseUpgrade(ctx, d, newP.(string), []string{oldP.(string)}, "plan", meta); err != nil {
			return diag.Errorf("error while waiting for instance %s to have updated plan : %s", d.Id(), err)
		}
	}

	return resourceVultrInstanceRead(ctx, d, meta)
}

func resourceVultrDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting instance (%s)", d.Id())

	if networkIDs, networkOK := d.GetOk("private_network_ids"); networkOK {
		detach := &govultr.InstanceUpdateReq{}
		for _, v := range networkIDs.(*schema.Set).List() {
			detach.DetachPrivateNetwork = append(detach.DetachPrivateNetwork, v.(string)) // nolint
		}

		if _, _, err := client.Instance.Update(ctx, d.Id(), detach); err != nil {
			return diag.Errorf("error detaching private networks prior to deleting instance %s : %v", d.Id(), err)
		}
	}

	if vpcIDs, vpcOK := d.GetOk("vpc_ids"); vpcOK {
		detach := &govultr.InstanceUpdateReq{}
		for _, v := range vpcIDs.(*schema.Set).List() {
			detach.DetachVPC = append(detach.DetachVPC, v.(string))
		}

		if _, _, err := client.Instance.Update(ctx, d.Id(), detach); err != nil {
			return diag.Errorf("error detaching VPCs prior to deleting instance %s : %v", d.Id(), err)
		}
	}

	if _, isoOK := d.GetOk("iso_id"); isoOK {
		if _, err := client.Instance.DetachISO(ctx, d.Id()); err != nil {
			return diag.Errorf("error detaching ISO prior to deleting instance %s : %v", d.Id(), err)
		}
	}

	if err := client.Instance.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying instance %s : %v", d.Id(), err)
	}

	return nil
}

func waitForDatabaseAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for Server (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{ // nolint:all
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newDatabaseStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newDatabaseStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc { // nolint:all
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating Server")
		server, _, err := client.Instance.Get(ctx, d.Id())
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

func waitForDatabaseUpgrade(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for instance (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{ // nolint:all
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newDatabasePlanRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newDatabasePlanRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc { // nolint:all
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Upgrading instance")
		instance, _, err := client.Instance.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving instance %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] The instances plan is %s", instance.Plan)
		return instance, instance.Plan, nil
	}
}
