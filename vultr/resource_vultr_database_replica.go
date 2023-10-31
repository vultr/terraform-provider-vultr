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
	"github.com/vultr/govultr/v3"
)

func resourceVultrDatabaseReplica() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrDatabaseReplicaCreate,
		ReadContext:   resourceVultrDatabaseReplicaRead,
		UpdateContext: resourceVultrDatabaseReplicaUpdate,
		DeleteContext: resourceVultrDatabaseReplicaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: readReplicaSchema(true),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceVultrDatabaseReplicaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	// Wait for at least one backup on the parent database to be available
	if _, err := waitForParentBackupAvailable(ctx, d, "yes", []string{"yes", "no"}, "latest_backup", meta); err != nil {
		return diag.Errorf(
			"error while waiting for parent Managed Database %s to have at least one backup : %s",
			d.Get("database_id").(string), err)
	}

	databaseID := d.Get("database_id").(string)

	req := &govultr.DatabaseAddReplicaReq{
		Region: d.Get("region").(string),
		Label:  d.Get("label").(string),
	}

	log.Printf("[INFO] Creating database read replica")
	database, _, err := client.Database.AddReadOnlyReplica(ctx, databaseID, req)
	if err != nil {
		return diag.Errorf("error creating database read replica: %v", err)
	}

	d.SetId(database.ID)

	if _, err = waitForDatabaseReplicaAvailable(ctx, d, "Running", []string{"Rebalancing", "Rebuilding", "Configuring", "Error"}, "status", meta); err != nil { //nolint:lll
		return diag.Errorf("error while waiting for Managed Database read replica %s to be in an active state : %s", d.Id(), err)
	}

	// Tags for read replicas can only be changed after creation
	if tag, tagOK := d.GetOk("tag"); tagOK {
		req2 := &govultr.DatabaseUpdateReq{
			Tag: tag.(string),
		}

		log.Printf("[INFO] Updating database read replica tag")
		if _, _, err := client.Database.Update(ctx, d.Id(), req2); err != nil {
			return diag.Errorf("error updating database read replica: %v", err)
		}
	}

	return resourceVultrDatabaseReplicaRead(ctx, d, meta)
}

func resourceVultrDatabaseReplicaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	database, _, err := client.Database.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "invalid database ID") {
			log.Printf("[WARN] Removing database read replica (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting database read replica (%s): %v", d.Id(), err)
	}

	if err := d.Set("date_created", database.DateCreated); err != nil {
		return diag.Errorf("unable to set resource database read replica `date_created` read value: %v", err)
	}

	if err := d.Set("plan", database.Plan); err != nil {
		return diag.Errorf("unable to set resource database read replica `plan` read value: %v", err)
	}

	if database.DatabaseEngine != "redis" {
		if err := d.Set("plan_disk", database.PlanDisk); err != nil {
			return diag.Errorf("unable to set resource database read replica `plan_disk` read value: %v", err)
		}
	}

	if err := d.Set("plan_ram", database.PlanRAM); err != nil {
		return diag.Errorf("unable to set resource database read replica `plan_ram` read value: %v", err)
	}

	if err := d.Set("plan_vcpus", database.PlanVCPUs); err != nil {
		return diag.Errorf("unable to set resource database read replica `plan_vcpus` read value: %v", err)
	}

	if err := d.Set("plan_replicas", database.PlanReplicas); err != nil {
		return diag.Errorf("unable to set resource database read replica `plan_replicas` read value: %v", err)
	}

	if err := d.Set("region", database.Region); err != nil {
		return diag.Errorf("unable to set resource database read replica `region` read value: %v", err)
	}

	if err := d.Set("vpc_id", database.VPCID); err != nil {
		return diag.Errorf("unable to set resource database read replica `vpc_id` read value: %v", err)
	}

	if err := d.Set("status", database.Status); err != nil {
		return diag.Errorf("unable to set resource database read replica `status` read value: %v", err)
	}

	if err := d.Set("label", database.Label); err != nil {
		return diag.Errorf("unable to set resource database read replica `label` read value: %v", err)
	}

	if err := d.Set("tag", database.Tag); err != nil {
		return diag.Errorf("unable to set resource database read replica `tag` read value: %v", err)
	}

	if err := d.Set("database_engine", database.DatabaseEngine); err != nil {
		return diag.Errorf("unable to set resource database read replica `database_engine` read value: %v", err)
	}

	if err := d.Set("database_engine_version", database.DatabaseEngineVersion); err != nil {
		return diag.Errorf("unable to set resource database read replica `database_engine_version` read value: %v", err)
	}

	if err := d.Set("dbname", database.DBName); err != nil {
		return diag.Errorf("unable to set resource database read replica `dbname` read value: %v", err)
	}

	if database.DatabaseEngine == "ferretpg" {
		if err := d.Set("ferretdb_credentials", flattenFerretDBCredentials(database)); err != nil {
			return diag.Errorf("unable to set resource database read replica `ferretdb_credentials` read value: %v", err)
		}
	}

	if err := d.Set("host", database.Host); err != nil {
		return diag.Errorf("unable to set resource database read replica `host` read value: %v", err)
	}

	if database.PublicHost != "" {
		if err := d.Set("public_host", database.PublicHost); err != nil {
			return diag.Errorf("unable to set resource database read replica `public_host` read value: %v", err)
		}
	}

	if err := d.Set("user", database.User); err != nil {
		return diag.Errorf("unable to set resource database read replica `user` read value: %v", err)
	}

	if err := d.Set("password", database.Password); err != nil {
		return diag.Errorf("unable to set resource database read replica `password` read value: %v", err)
	}

	if err := d.Set("port", database.Port); err != nil {
		return diag.Errorf("unable to set resource database read replica `port` read value: %v", err)
	}

	if err := d.Set("maintenance_dow", database.MaintenanceDOW); err != nil {
		return diag.Errorf("unable to set resource database read replica `maintenance_dow` read value: %v", err)
	}

	if err := d.Set("maintenance_time", database.MaintenanceTime); err != nil {
		return diag.Errorf("unable to set resource database read replica `maintenance_time` read value: %v", err)
	}

	if err := d.Set("latest_backup", database.LatestBackup); err != nil {
		return diag.Errorf("unable to set resource database read replica `latest_backup` read value: %v", err)
	}

	if err := d.Set("trusted_ips", database.TrustedIPs); err != nil {
		return diag.Errorf("unable to set resource database read replica `trusted_ips` read value: %v", err)
	}

	if database.DatabaseEngine == "mysql" {
		if err := d.Set("mysql_sql_modes", database.MySQLSQLModes); err != nil {
			return diag.Errorf("unable to set resource database read replica `mysql_sql_modes` read value: %v", err)
		}

		if err := d.Set("mysql_require_primary_key", database.MySQLRequirePrimaryKey); err != nil {
			return diag.Errorf("unable to set resource database read replica `mysql_require_primary_key` read value: %v", err)
		}

		if err := d.Set("mysql_slow_query_log", database.MySQLSlowQueryLog); err != nil {
			return diag.Errorf("unable to set resource database read replica `mysql_slow_query_log` read value: %v", err)
		}

		if err := d.Set("mysql_long_query_time", database.MySQLLongQueryTime); err != nil {
			return diag.Errorf("unable to set resource database read replica `mysql_long_query_time` read value: %v", err)
		}
	}

	if database.DatabaseEngine == "redis" {
		if err := d.Set("redis_eviction_policy", database.RedisEvictionPolicy); err != nil {
			return diag.Errorf("unable to set resource database read replica `redis_eviction_policy` read value: %v", err)
		}
	}

	if database.DatabaseEngine != "redis" {
		if err := d.Set("cluster_time_zone", database.ClusterTimeZone); err != nil {
			return diag.Errorf("unable to set resource database read replica `cluster_time_zone` read value: %v", err)
		}
	}

	return nil
}
func resourceVultrDatabaseReplicaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if _, _, err := client.Database.Update(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating database read replica %s : %s", d.Id(), err.Error())
	}

	if d.HasChange("region") || d.HasChange("vpc_id") {
		if _, err := waitForDatabaseReplicaAvailable(ctx, d, "Running", []string{"Rebalancing", "Rebuilding", "Configuring", "Error"}, "status", meta); err != nil { //nolint:lll
			return diag.Errorf("error while waiting for Managed Database read replica %s to be in an active state : %s", d.Id(), err)
		}
	}

	return resourceVultrDatabaseReplicaRead(ctx, d, meta)
}

func resourceVultrDatabaseReplicaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Client).govultrClient()
	log.Printf("[INFO] Deleting database read replica (%s)", d.Id())

	if err := client.Database.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying database read replica %s : %v", d.Id(), err)
	}

	return nil
}

func waitForDatabaseReplicaAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll,dupl
	log.Printf(
		"[INFO] Waiting for Managed Database read replica (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{ // nolint:all
		Pending:        pending,
		Target:         []string{target},
		Refresh:        newDatabaseReplicaStateRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newDatabaseReplicaStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc { // nolint:all
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Creating Database read replica")
		server, _, err := client.Database.Get(ctx, d.Id())

		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Managed Database read replica %s : %s", d.Id(), err)
		}

		if attr == "status" {
			log.Printf("[INFO] The Managed Database read replica Status is %s", server.Status)
			return server, server.Status, nil
		}

		return nil, "", nil
	}
}

func waitForParentBackupAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) { //nolint:lll
	log.Printf(
		"[INFO] Waiting for parent Managed Database (%s) to have %s of %s",
		d.Get("database_id").(string), attribute, target)

	stateConf := &resource.StateChangeConf{ // nolint:all
		Pending:        pending,
		Target:         []string{target},
		Refresh:        parentDatabaseRefresh(ctx, d, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func parentDatabaseRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}, attr string) resource.StateRefreshFunc { // nolint:all
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Waiting for parent Managed Database backup status")
		server, _, err := client.Database.Get(ctx, d.Get("database_id").(string))

		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Managed Database %s : %s", d.Get("database_id"), err)
		}

		if attr == "latest_backup" {
			log.Printf("[INFO] The Managed Database read replica LatestBackup is %s", server.LatestBackup)
			if server.LatestBackup == "" {
				return server, "no", nil
			}
			return server, "yes", nil
		}

		return nil, "", nil
	}
}
