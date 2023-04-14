package vultr

import (
	"context"
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
			"mysql_require_primary_key": {
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
			"plan_ram": {
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

	req := &govultr.DatabaseCreateReq{
		DatabaseEngine:         d.Get("database_engine").(string),
		DatabaseEngineVersion:  d.Get("database_engine_version").(string),
		Region:                 d.Get("region").(string),
		Plan:                   d.Get("plan").(string),
		Label:                  d.Get("label").(string),
		Tag:                    d.Get("tag").(string),
		MaintenanceDOW:         d.Get("maintenance_dow").(string),
		MaintenanceTime:        d.Get("maintenance_time").(string),
		MySQLRequirePrimaryKey: govultr.BoolToBoolPtr(true),
		MySQLSlowQueryLog:      govultr.BoolToBoolPtr(false),
		MySQLLongQueryTime:     d.Get("mysql_long_query_time").(int),
		RedisEvictionPolicy:    d.Get("redis_eviction_policy").(string),
	}

	if trustedIPs, trustedIPsOK := d.GetOk("trusted_ips"); trustedIPsOK {
		for _, v := range trustedIPs.(*schema.Set).List() {
			req.TrustedIPs = append(req.TrustedIPs, v.(string))
		}
	}

	if mysqlSQLModes, mysqlSQLModesOK := d.GetOk("mysql_sql_modes"); mysqlSQLModesOK {
		for _, v := range mysqlSQLModes.(*schema.Set).List() {
			req.MySQLSQLModes = append(req.MySQLSQLModes, v.(string))
		}
	}

	if mysqlRequirePrimaryKey, mysqlRequirePrimaryKeyOK := d.GetOk("mysql_require_primary_key"); mysqlRequirePrimaryKeyOK {
		req.MySQLRequirePrimaryKey = govultr.BoolToBoolPtr(mysqlRequirePrimaryKey.(bool))
	}

	if mysqlSlowQueryLog, mysqlSlowQueryLogOK := d.GetOk("mysql_require_primary_key"); mysqlSlowQueryLogOK {
		req.MySQLSlowQueryLog = govultr.BoolToBoolPtr(mysqlSlowQueryLog.(bool))
	}

	log.Printf("[INFO] Creating database")
	database, _, err := client.Database.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating database: %v", err)
	}

	d.SetId(database.ID)

	if _, err = waitForDatabaseAvailable(ctx, d, "Running", []string{"Rebalancing", "Rebuilding", "Error"}, "status", meta); err != nil {
		return diag.Errorf("error while waiting for Managed Database %s to be in an active state : %s", d.Id(), err)
	}

	return resourceVultrDatabaseRead(ctx, d, meta)
}

func resourceVultrDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	database, _, err := client.Database.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "invalid database ID") {
			log.Printf("[WARN] Removing database (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting database (%s): %v", d.Id(), err)
	}

	if err := d.Set("date_created", database.DateCreated); err != nil {
		return diag.Errorf("unable to set resource database `date_created` read value: %v", err)
	}

	if err := d.Set("plan", database.Plan); err != nil {
		return diag.Errorf("unable to set resource database `plan` read value: %v", err)
	}

	if err := d.Set("plan_disk", database.PlanDisk); err != nil {
		return diag.Errorf("unable to set resource database `plan_disk` read value: %v", err)
	}

	if err := d.Set("plan_ram", database.PlanRAM); err != nil {
		return diag.Errorf("unable to set resource database `plan_ram` read value: %v", err)
	}

	if err := d.Set("plan_vcpus", database.PlanVCPUs); err != nil {
		return diag.Errorf("unable to set resource database `plan_vcpus` read value: %v", err)
	}

	if err := d.Set("plan_replicas", database.PlanReplicas); err != nil {
		return diag.Errorf("unable to set resource database `plan_replicas` read value: %v", err)
	}

	if err := d.Set("region", database.Region); err != nil {
		return diag.Errorf("unable to set resource database `region` read value: %v", err)
	}

	if err := d.Set("status", database.Status); err != nil {
		return diag.Errorf("unable to set resource database `status` read value: %v", err)
	}

	if err := d.Set("label", database.Label); err != nil {
		return diag.Errorf("unable to set resource database `label` read value: %v", err)
	}

	if err := d.Set("tag", database.Tag); err != nil {
		return diag.Errorf("unable to set resource database `tag` read value: %v", err)
	}

	if err := d.Set("database_engine", database.DatabaseEngine); err != nil {
		return diag.Errorf("unable to set resource database `database_engine` read value: %v", err)
	}

	if err := d.Set("database_engine_version", database.DatabaseEngineVersion); err != nil {
		return diag.Errorf("unable to set resource database `database_engine_version` read value: %v", err)
	}

	if err := d.Set("dbname", database.DBName); err != nil {
		return diag.Errorf("unable to set resource database `dbname` read value: %v", err)
	}

	if err := d.Set("host", database.Host); err != nil {
		return diag.Errorf("unable to set resource database `host` read value: %v", err)
	}

	if err := d.Set("user", database.User); err != nil {
		return diag.Errorf("unable to set resource database `user` read value: %v", err)
	}

	if err := d.Set("password", database.Password); err != nil {
		return diag.Errorf("unable to set resource database `password` read value: %v", err)
	}

	if err := d.Set("port", database.Port); err != nil {
		return diag.Errorf("unable to set resource database `port` read value: %v", err)
	}

	if err := d.Set("maintenance_dow", database.MaintenanceDOW); err != nil {
		return diag.Errorf("unable to set resource database `maintenance_dow` read value: %v", err)
	}

	if err := d.Set("maintenance_time", database.MaintenanceTime); err != nil {
		return diag.Errorf("unable to set resource database `maintenance_time` read value: %v", err)
	}

	if err := d.Set("latest_backup", database.LatestBackup); err != nil {
		return diag.Errorf("unable to set resource database `latest_backup` read value: %v", err)
	}

	if err := d.Set("trusted_ips", database.TrustedIPs); err != nil {
		return diag.Errorf("unable to set resource database `trusted_ips` read value: %v", err)
	}

	if err := d.Set("mysql_sql_modes", database.MySQLSQLModes); err != nil {
		return diag.Errorf("unable to set resource database `mysql_sql_modes` read value: %v", err)
	}

	if err := d.Set("mysql_require_primary_key", database.MySQLRequirePrimaryKey); err != nil {
		return diag.Errorf("unable to set resource database `mysql_require_primary_key` read value: %v", err)
	}

	if err := d.Set("mysql_slow_query_log", database.MySQLSlowQueryLog); err != nil {
		return diag.Errorf("unable to set resource database `mysql_slow_query_log` read value: %v", err)
	}

	if err := d.Set("mysql_long_query_time", database.MySQLLongQueryTime); err != nil {
		return diag.Errorf("unable to set resource database `mysql_long_query_time` read value: %v", err)
	}

	if err := d.Set("redis_eviction_policy", database.RedisEvictionPolicy); err != nil {
		return diag.Errorf("unable to set resource database `redis_eviction_policy` read value: %v", err)
	}

	if err := d.Set("cluster_time_zone", database.ClusterTimeZone); err != nil {
		return diag.Errorf("unable to set resource database `cluster_time_zone` read value: %v", err)
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

	return resourceVultrDatabaseRead(ctx, d, meta)
}

func resourceVultrDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting database (%s)", d.Id())

	if err := client.Database.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying database %s : %v", d.Id(), err)
	}

	return nil
}

func waitForDatabaseAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for Managed Database (%s) to have %s of %s",
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

		log.Printf("[INFO] Creating Database")
		server, _, err := client.Database.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Managed Database %s : %s", d.Id(), err)
		}

		if attr == "status" {
			log.Printf("[INFO] The Managed Database Status is %s", server.Status)
			return server, server.Status, nil
		}

		return nil, "", nil
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
