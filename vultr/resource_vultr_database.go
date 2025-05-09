package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

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
			"database_engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_engine_version": {
				Type:     schema.TypeString,
				Required: true,
			},
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
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maintenance_dow": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"maintenance_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"backup_hour": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_minute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_time_zone": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"trusted_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mysql_sql_modes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mysql_slow_query_log": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"mysql_require_primary_key": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mysql_long_query_time": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"eviction_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_disk": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
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
				Optional: true,
			},
			"plan_brokers": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ferretdb_credentials": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				Elem:     schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("database_engine") != "ferretpg"
				},
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_host": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sasl_port": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"access_cert": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"latest_backup": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_replicas": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: readReplicaSchema(false),
				},
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
		VPCID:                  d.Get("vpc_id").(string),
		MaintenanceDOW:         d.Get("maintenance_dow").(string),
		MaintenanceTime:        d.Get("maintenance_time").(string),
		MySQLRequirePrimaryKey: govultr.BoolToBoolPtr(true),
		MySQLSlowQueryLog:      govultr.BoolToBoolPtr(false),
		MySQLLongQueryTime:     d.Get("mysql_long_query_time").(int),
		EvictionPolicy:         d.Get("eviction_policy").(string),
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

	mysqlRequirePrimaryKey, mysqlRequirePrimaryKeyOK := d.GetOk("mysql_require_primary_key")
	if mysqlRequirePrimaryKeyOK {
		req.MySQLRequirePrimaryKey = govultr.BoolToBoolPtr(mysqlRequirePrimaryKey.(bool))
	}

	if mysqlSlowQueryLog, mysqlSlowQueryLogOK := d.GetOk("mysql_slow_query_log"); mysqlSlowQueryLogOK {
		req.MySQLSlowQueryLog = govultr.BoolToBoolPtr(mysqlSlowQueryLog.(bool))
	}

	if backupHour, backupHourOK := d.GetOk("backup_hour"); backupHourOK {
		req.BackupHour = govultr.StringToStringPtr(backupHour.(string))
	}

	if backupMinute, backupMinuteOK := d.GetOk("backup_minute"); backupMinuteOK {
		req.BackupMinute = govultr.StringToStringPtr(backupMinute.(string))
	}

	log.Printf("[INFO] Creating database")
	database, _, err := client.Database.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating database: %v", err)
	}

	d.SetId(database.ID)
	pendStatuses := []string{"Rebalancing", "Rebuilding", "Configuring", "Error"}
	_, errWait := waitForDatabaseAvailable(ctx, d, "Running", pendStatuses, "status", meta)
	if errWait != nil {
		return diag.Errorf("error while waiting for Managed Database %s to be in an active state : %s", d.Id(), err)
	}

	// Some values can only be properly set after creation
	req2 := &govultr.DatabaseUpdateReq{}
	if clusterTimeZone, clusterTimeZoneOK := d.GetOk("cluster_time_zone"); clusterTimeZoneOK {
		log.Printf("[INFO] Updating database default time zone")
		req2.ClusterTimeZone = clusterTimeZone.(string)
	}

	// Empty/default backup schedule settings are only honored on update and not create
	if req.DatabaseEngine != "kafka" && req.BackupHour == nil && req.BackupMinute == nil {
		log.Printf("[INFO] Updating database backup schedule")
		req2.BackupHour = govultr.StringToStringPtr("")
		req2.BackupMinute = govultr.StringToStringPtr("")
	}

	// Perform an update if needed
	if req2.ClusterTimeZone != "" || req2.BackupHour != nil || req2.BackupMinute != nil {
		if _, _, err := client.Database.Update(ctx, d.Id(), req2); err != nil {
			return diag.Errorf("error updating post-creation values for database: %v", err)
		}
	}

	// Default user (vultradmin) password can only be changed after creation
	if password, passwordOK := d.GetOk("password"); passwordOK && d.Get("database_engine").(string) != "valkey" { //nolint:lll
		req3 := &govultr.DatabaseUserUpdateReq{
			Password: password.(string),
		}

		log.Printf("[INFO] Updating default user password")
		if _, _, err := client.Database.UpdateUser(ctx, d.Id(), "vultradmin", req3); err != nil {
			return diag.Errorf("error updating default user: %v", err)
		}
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

	if database.DatabaseEngine != "valkey" {
		if err := d.Set("plan_disk", database.PlanDisk); err != nil {
			return diag.Errorf("unable to set resource database `plan_disk` read value: %v", err)
		}
	}

	if err := d.Set("plan_ram", database.PlanRAM); err != nil {
		return diag.Errorf("unable to set resource database `plan_ram` read value: %v", err)
	}

	if err := d.Set("plan_vcpus", database.PlanVCPUs); err != nil {
		return diag.Errorf("unable to set resource database `plan_vcpus` read value: %v", err)
	}

	if database.DatabaseEngine != "kafka" {
		if err := d.Set("plan_replicas", database.PlanReplicas); err != nil {
			return diag.Errorf("unable to set resource database `plan_replicas` read value: %v", err)
		}
	} else {
		if err := d.Set("plan_brokers", database.PlanBrokers); err != nil {
			return diag.Errorf("unable to set resource database `plan_brokers` read value: %v", err)
		}
	}

	if err := d.Set("region", database.Region); err != nil {
		return diag.Errorf("unable to set resource database `region` read value: %v", err)
	}

	if err := d.Set("database_engine", database.DatabaseEngine); err != nil {
		return diag.Errorf("unable to set resource database `database_engine` read value: %v", err)
	}

	if err := d.Set("database_engine_version", database.DatabaseEngineVersion); err != nil {
		return diag.Errorf("unable to set resource database `database_engine_version` read value: %v", err)
	}

	if err := d.Set("vpc_id", database.VPCID); err != nil {
		return diag.Errorf("unable to set resource database `vpc_id` read value: %v", err)
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

	if err := d.Set("dbname", database.DBName); err != nil {
		return diag.Errorf("unable to set resource database `dbname` read value: %v", err)
	}

	if database.DatabaseEngine == "ferretpg" {
		if err := d.Set("ferretdb_credentials", flattenFerretDBCredentials(database)); err != nil {
			return diag.Errorf("unable to set resource database `ferretdb_credentials` read value: %v", err)
		}
	}

	if err := d.Set("host", database.Host); err != nil {
		return diag.Errorf("unable to set resource database `host` read value: %v", err)
	}

	if database.PublicHost != "" {
		if err := d.Set("public_host", database.PublicHost); err != nil {
			return diag.Errorf("unable to set resource database `public_host` read value: %v", err)
		}
	}

	if err := d.Set("port", database.Port); err != nil {
		return diag.Errorf("unable to set resource database `port` read value: %v", err)
	}

	if database.DatabaseEngine == "kafka" {
		if err := d.Set("sasl_port", database.SASLPort); err != nil {
			return diag.Errorf("unable to set resource database `sasl_port` read value: %v", err)
		}
	}

	if err := d.Set("user", database.User); err != nil {
		return diag.Errorf("unable to set resource database `user` read value: %v", err)
	}

	if err := d.Set("password", database.Password); err != nil {
		return diag.Errorf("unable to set resource database `password` read value: %v", err)
	}

	if database.DatabaseEngine == "kafka" {
		if err := d.Set("access_key", database.AccessKey); err != nil {
			return diag.Errorf("unable to set resource database `access_key` read value: %v", err)
		}

		if err := d.Set("access_cert", database.AccessCert); err != nil {
			return diag.Errorf("unable to set resource database `access_cert` read value: %v", err)
		}
	}

	if err := d.Set("maintenance_dow", database.MaintenanceDOW); err != nil {
		return diag.Errorf("unable to set resource database `maintenance_dow` read value: %v", err)
	}

	if err := d.Set("maintenance_time", database.MaintenanceTime); err != nil {
		return diag.Errorf("unable to set resource database `maintenance_time` read value: %v", err)
	}

	if database.DatabaseEngine != "kafka" {
		if err := d.Set("backup_hour", *database.BackupHour); err != nil {
			return diag.Errorf("unable to set resource database `backup_hour` read value: %v", err)
		}

		if err := d.Set("backup_minute", *database.BackupMinute); err != nil {
			return diag.Errorf("unable to set resource database `backup_minute` read value: %v", err)
		}
	}

	if err := d.Set("latest_backup", database.LatestBackup); err != nil {
		return diag.Errorf("unable to set resource database `latest_backup` read value: %v", err)
	}

	if err := d.Set("trusted_ips", database.TrustedIPs); err != nil {
		return diag.Errorf("unable to set resource database `trusted_ips` read value: %v", err)
	}

	if database.DatabaseEngine == "mysql" {
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
	}

	if database.DatabaseEngine == "valkey" {
		if err := d.Set("eviction_policy", database.EvictionPolicy); err != nil {
			return diag.Errorf("unable to set resource database `eviction_policy` read value: %v", err)
		}
	}

	if database.DatabaseEngine != "valkey" {
		if err := d.Set("cluster_time_zone", database.ClusterTimeZone); err != nil {
			return diag.Errorf("unable to set resource database `cluster_time_zone` read value: %v", err)
		}
	}

	if err := d.Set("read_replicas", flattenReplicas(database)); err != nil {
		return diag.Errorf("unable to set resource database `read_replicas` read value: %v", err)
	}

	return nil
}
func resourceVultrDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.DatabaseUpdateReq{
		Label: d.Get("label").(string),
	}

	if d.HasChange("region") {
		log.Printf("[INFO] Updating Region")
		_, newVal := d.GetChange("region")
		region := newVal.(string)
		req.Region = region
	}

	if d.HasChange("plan") {
		log.Printf("[INFO] Updating Plan")
		_, newVal := d.GetChange("plan")
		plan := newVal.(string)
		req.Plan = plan
	}

	if d.HasChange("tag") {
		log.Printf("[INFO] Updating Tag")
		_, newVal := d.GetChange("tag")
		tag := newVal.(string)
		req.Tag = tag
	}

	if d.HasChange("vpc_id") {
		log.Printf("[INFO] Updating VPC ID")
		_, newVal := d.GetChange("vpc_id")
		vpc := newVal.(string)
		req.VPCID = govultr.StringToStringPtr(vpc)
	}

	if d.HasChange("maintenance_dow") {
		log.Printf("[INFO] Updating Maintenance DOW")
		_, newVal := d.GetChange("maintenance_dow")
		maintenanceDOW := newVal.(string)
		req.MaintenanceDOW = maintenanceDOW
	}

	if d.HasChange("maintenance_time") {
		log.Printf("[INFO] Updating Maintenance Time")
		_, newVal := d.GetChange("maintenance_time")
		maintenanceTime := newVal.(string)
		req.MaintenanceTime = maintenanceTime
	}

	if d.HasChange("backup_hour") {
		log.Printf("[INFO] Updating Backup Hour")
		_, newVal := d.GetChange("backup_hour")
		backupHour := newVal.(string)
		req.BackupHour = govultr.StringToStringPtr(backupHour)
	}

	if d.HasChange("backup_minute") {
		log.Printf("[INFO] Updating Backup Minute")
		_, newVal := d.GetChange("backup_minute")
		backupMinute := newVal.(string)
		req.BackupMinute = govultr.StringToStringPtr(backupMinute)
	}

	if d.HasChange("cluster_time_zone") {
		log.Printf("[INFO] Updating Cluster Time Zone")
		_, newVal := d.GetChange("cluster_time_zone")
		clusterTimeZone := newVal.(string)
		req.ClusterTimeZone = clusterTimeZone
	}

	if d.HasChange("trusted_ips") {
		log.Printf("[INFO] Updating Trusted IPs")
		_, newVal := d.GetChange("trusted_ips")

		var newIPs []string
		for _, v := range newVal.(*schema.Set).List() {
			newIPs = append(newIPs, v.(string))
		}

		req.TrustedIPs = newIPs
	}

	if d.HasChange("mysql_sql_modes") {
		log.Printf("[INFO] Updating MySQL SQL Modes")
		_, newVal := d.GetChange("mysql_sql_modes")

		var newModes []string
		for _, v := range newVal.(*schema.Set).List() {
			newModes = append(newModes, v.(string))
		}

		req.MySQLSQLModes = newModes
	}

	if d.HasChange("mysql_require_primary_key") {
		log.Printf("[INFO] Updating MySQL Require Primary Key")
		_, newVal := d.GetChange("mysql_require_primary_key")
		mysqlRequirePrimaryKey := newVal.(bool)
		req.MySQLRequirePrimaryKey = &mysqlRequirePrimaryKey
	}

	if d.HasChange("mysql_slow_query_log") {
		log.Printf("[INFO] Updating MySQL Slow Query Log")
		_, newVal := d.GetChange("mysql_slow_query_log")
		mysqlSlowQueryLog := newVal.(bool)
		req.MySQLSlowQueryLog = &mysqlSlowQueryLog
	}

	if d.HasChange("mysql_long_query_time") {
		log.Printf("[INFO] Updating MySQL Long Query Time")
		_, newVal := d.GetChange("mysql_long_query_time")
		mysqlLongQueryTime := newVal.(int)
		req.MySQLLongQueryTime = mysqlLongQueryTime
	}

	if d.HasChange("eviction_policy") {
		log.Printf("[INFO] Updating Eviction Policy")
		_, newVal := d.GetChange("eviction_policy")
		evictionPolicy := newVal.(string)
		req.EvictionPolicy = evictionPolicy
	}

	if _, _, err := client.Database.Update(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating database %s : %s", d.Id(), err.Error())
	}

	if d.HasChange("region") || d.HasChange("plan") || d.HasChange("vpc_id") {
		pendStatuses := []string{"Rebalancing", "Rebuilding", "Configuring", "Error"}
		_, errAvail := waitForDatabaseAvailable(ctx, d, "Running", pendStatuses, "status", meta)
		if errAvail != nil {
			return diag.Errorf(
				"error while waiting for Managed Database %s to be in an active state : %s",
				d.Id(),
				errAvail,
			)
		}
	}

	// Updating the default user password requires a separate API call
	if d.HasChange("password") && d.Get("database_engine").(string) != "valkey" {
		_, newVal := d.GetChange("password")
		password := newVal.(string)
		reqP := &govultr.DatabaseUserUpdateReq{
			Password: password,
		}

		log.Printf("[INFO] Updating default user password")
		if _, _, err := client.Database.UpdateUser(ctx, d.Id(), "vultradmin", reqP); err != nil {
			return diag.Errorf("error updating default user: %v", err)
		}
	}

	// Version changes have their own API protocol/checks
	if d.HasChange("database_engine_version") {
		// Check available versions against input
		log.Printf("[INFO] Checking available version upgrades")
		availableVersions, _, err := client.Database.ListAvailableVersions(ctx, d.Id())
		if err != nil {
			return diag.Errorf("error checking available version upgrades %s : %s", d.Id(), err.Error())
		}
		_, newVal := d.GetChange("database_engine_version")
		databaseEngineVersion := newVal.(string)
		if !versionCompare(availableVersions, databaseEngineVersion) {
			return diag.Errorf("invalid version %s provided for database %s", databaseEngineVersion, d.Id())
		}

		// Start version upgrade
		log.Printf("[INFO] Initiating version upgrade")
		req2 := &govultr.DatabaseVersionUpgradeReq{
			Version: databaseEngineVersion,
		}
		if _, _, err := client.Database.StartVersionUpgrade(ctx, d.Id(), req2); err != nil {
			return diag.Errorf("error upgrading database version %s : %s", d.Id(), err.Error())
		}

		// Wait for running state
		pendStatuses := []string{"Rebalancing", "Rebuilding", "Configuring", "Error"}
		_, errAvail := waitForDatabaseAvailable(ctx, d, "Running", pendStatuses, "status", meta)
		if errAvail != nil {
			return diag.Errorf(
				"error while waiting for Managed Database %s to be in an active state : %s",
				d.Id(),
				errAvail,
			)
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

func waitForDatabaseAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll
	log.Printf(
		"[INFO] Waiting for Managed Database (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
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

func newDatabaseStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) retry.StateRefreshFunc { //nolint:lll
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

func versionCompare(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
