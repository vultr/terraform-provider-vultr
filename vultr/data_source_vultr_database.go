package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrDatabaseRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": {
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
			"public_host": {
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
			"maintenance_dow": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_backup": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trusted_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"mysql_sql_modes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"mysql_require_primary_key": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mysql_slow_query_log": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mysql_long_query_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"redis_eviction_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_replicas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: readReplicaSchema(false),
				},
			},
		},
	}
}

func dataSourceVultrDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var databaseList []govultr.Database
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.DBListOptions{}
	databases, _, _, err := client.Database.List(ctx, options)
	if err != nil {
		return diag.Errorf("error getting databases: %v", err)
	}

	for s := range databases {
		// we need convert the a struct INTO a map so we can easily manipulate the data here
		sm, err := structToMap(databases[s])

		if err != nil {
			return diag.FromErr(err)
		}

		if filterLoop(f, sm) {
			databaseList = append(databaseList, databases[s])
		}
	}

	if len(databaseList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(databaseList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(databaseList[0].ID)

	if err := d.Set("date_created", databaseList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set database `date_created` read value: %v", err)
	}

	if err := d.Set("plan", databaseList[0].Plan); err != nil {
		return diag.Errorf("unable to set resource database `plan` read value: %v", err)
	}

	if databaseList[0].DatabaseEngine != "redis" {
		if err := d.Set("plan_disk", databaseList[0].PlanDisk); err != nil {
			return diag.Errorf("unable to set resource database `plan_disk` read value: %v", err)
		}
	}

	if err := d.Set("plan_ram", databaseList[0].PlanRAM); err != nil {
		return diag.Errorf("unable to set resource database `plan_ram` read value: %v", err)
	}

	if err := d.Set("plan_vcpus", databaseList[0].PlanVCPUs); err != nil {
		return diag.Errorf("unable to set resource database `plan_vcpus` read value: %v", err)
	}

	if err := d.Set("plan_replicas", databaseList[0].PlanReplicas); err != nil {
		return diag.Errorf("unable to set resource database `plan_replicas` read value: %v", err)
	}

	if err := d.Set("region", databaseList[0].Region); err != nil {
		return diag.Errorf("unable to set resource database `region` read value: %v", err)
	}

	if err := d.Set("database_engine", databaseList[0].DatabaseEngine); err != nil {
		return diag.Errorf("unable to set resource database `database_engine` read value: %v", err)
	}

	if err := d.Set("database_engine_version", databaseList[0].DatabaseEngineVersion); err != nil {
		return diag.Errorf("unable to set resource database `database_engine_version` read value: %v", err)
	}

	if err := d.Set("vpc_id", databaseList[0].VPCID); err != nil {
		return diag.Errorf("unable to set resource database `vpc_id` read value: %v", err)
	}

	if err := d.Set("status", databaseList[0].Status); err != nil {
		return diag.Errorf("unable to set resource database `status` read value: %v", err)
	}

	if err := d.Set("label", databaseList[0].Label); err != nil {
		return diag.Errorf("unable to set resource database `label` read value: %v", err)
	}

	if err := d.Set("tag", databaseList[0].Tag); err != nil {
		return diag.Errorf("unable to set resource database `tag` read value: %v", err)
	}

	if err := d.Set("dbname", databaseList[0].DBName); err != nil {
		return diag.Errorf("unable to set resource database `dbname` read value: %v", err)
	}

	if err := d.Set("host", databaseList[0].Host); err != nil {
		return diag.Errorf("unable to set resource database `host` read value: %v", err)
	}

	if databaseList[0].PublicHost != "" {
		if err := d.Set("public_host", databaseList[0].PublicHost); err != nil {
			return diag.Errorf("unable to set resource database `public_host` read value: %v", err)
		}
	}

	if err := d.Set("user", databaseList[0].User); err != nil {
		return diag.Errorf("unable to set resource database `user` read value: %v", err)
	}

	if err := d.Set("password", databaseList[0].Password); err != nil {
		return diag.Errorf("unable to set resource database `password` read value: %v", err)
	}

	if err := d.Set("port", databaseList[0].Port); err != nil {
		return diag.Errorf("unable to set resource database `port` read value: %v", err)
	}

	if err := d.Set("maintenance_dow", databaseList[0].MaintenanceDOW); err != nil {
		return diag.Errorf("unable to set resource database `maintenance_dow` read value: %v", err)
	}

	if err := d.Set("maintenance_time", databaseList[0].MaintenanceTime); err != nil {
		return diag.Errorf("unable to set resource database `maintenance_time` read value: %v", err)
	}

	if err := d.Set("latest_backup", databaseList[0].LatestBackup); err != nil {
		return diag.Errorf("unable to set resource database `latest_backup` read value: %v", err)
	}

	if err := d.Set("trusted_ips", databaseList[0].TrustedIPs); err != nil {
		return diag.Errorf("unable to set resource database `trusted_ips` read value: %v", err)
	}

	if databaseList[0].DatabaseEngine == "mysql" {
		if err := d.Set("mysql_sql_modes", databaseList[0].MySQLSQLModes); err != nil {
			return diag.Errorf("unable to set resource database `mysql_sql_modes` read value: %v", err)
		}

		if err := d.Set("mysql_require_primary_key", databaseList[0].MySQLRequirePrimaryKey); err != nil {
			return diag.Errorf("unable to set resource database `mysql_require_primary_key` read value: %v", err)
		}

		if err := d.Set("mysql_slow_query_log", databaseList[0].MySQLSlowQueryLog); err != nil {
			return diag.Errorf("unable to set resource database `mysql_slow_query_log` read value: %v", err)
		}

		if err := d.Set("mysql_long_query_time", databaseList[0].MySQLLongQueryTime); err != nil {
			return diag.Errorf("unable to set resource database `mysql_long_query_time` read value: %v", err)
		}
	}

	if databaseList[0].DatabaseEngine == "redis" {
		if err := d.Set("redis_eviction_policy", databaseList[0].RedisEvictionPolicy); err != nil {
			return diag.Errorf("unable to set resource database `redis_eviction_policy` read value: %v", err)
		}
	}

	if databaseList[0].DatabaseEngine != "redis" {
		if err := d.Set("cluster_time_zone", databaseList[0].ClusterTimeZone); err != nil {
			return diag.Errorf("unable to set resource database `cluster_time_zone` read value: %v", err)
		}
	}

	if err := d.Set("read_replicas", flattenReplicas(&databaseList[0])); err != nil {
		return diag.Errorf("unable to set resource database `read_replicas` read value: %v", err)
	}

	return nil
}

func flattenReplicas(db *govultr.Database) []map[string]interface{} {
	var replicas []map[string]interface{}
	for v := range db.ReadReplicas {
		r := map[string]interface{}{
			"id":                        db.ReadReplicas[v].ID,
			"date_created":              db.ReadReplicas[v].DateCreated,
			"plan":                      db.ReadReplicas[v].Plan,
			"plan_disk":                 db.ReadReplicas[v].PlanDisk,
			"plan_ram":                  db.ReadReplicas[v].PlanRAM,
			"plan_vcpus":                db.ReadReplicas[v].PlanVCPUs,
			"plan_replicas":             db.ReadReplicas[v].PlanReplicas,
			"region":                    db.ReadReplicas[v].Region,
			"database_engine":           db.ReadReplicas[v].DatabaseEngine,
			"database_engine_version":   db.ReadReplicas[v].DatabaseEngineVersion,
			"vpc_id":                    db.ReadReplicas[v].VPCID,
			"status":                    db.ReadReplicas[v].Status,
			"label":                     db.ReadReplicas[v].Label,
			"tag":                       db.ReadReplicas[v].Tag,
			"dbname":                    db.ReadReplicas[v].DBName,
			"host":                      db.ReadReplicas[v].Host,
			"public_host":               db.ReadReplicas[v].PublicHost,
			"user":                      db.ReadReplicas[v].User,
			"password":                  db.ReadReplicas[v].Password,
			"port":                      db.ReadReplicas[v].Port,
			"maintenance_dow":           db.ReadReplicas[v].MaintenanceDOW,
			"maintenance_time":          db.ReadReplicas[v].MaintenanceTime,
			"latest_backup":             db.ReadReplicas[v].LatestBackup,
			"trusted_ips":               db.ReadReplicas[v].TrustedIPs,
			"mysql_sql_modes":           db.ReadReplicas[v].MySQLSQLModes,
			"mysql_require_primary_key": db.ReadReplicas[v].MySQLRequirePrimaryKey,
			"mysql_slow_query_log":      db.ReadReplicas[v].MySQLSlowQueryLog,
			"mysql_long_query_time":     db.ReadReplicas[v].MySQLLongQueryTime,
			"redis_eviction_policy":     db.ReadReplicas[v].RedisEvictionPolicy,
			"cluster_time_zone":         db.ReadReplicas[v].ClusterTimeZone,
		}

		if db.PublicHost == "" {
			delete(r, "public_host")
		}

		if db.DatabaseEngine != "mysql" {
			delete(r, "mysql_sql_modes")
			delete(r, "mysql_require_primary_key")
			delete(r, "mysql_slow_query_log")
			delete(r, "mysql_sql_long_query_time")
		}

		if db.DatabaseEngine != "redis" {
			delete(r, "redis_eviction_policy")
		}

		if db.DatabaseEngine == "redis" {
			delete(r, "cluster_time_zone")
		}

		replicas = append(replicas, r)
	}

	return replicas
}
