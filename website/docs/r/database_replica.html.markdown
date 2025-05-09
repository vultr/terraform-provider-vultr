---
layout: "vultr"
page_title: "Vultr: vultr_database_replica"
sidebar_current: "docs-vultr-resource-database-replica"
description: |-
  Provides a Vultr database replica resource. This can be used to create, read, modify, and delete managed database read replicas on your Vultr account.
---

# vultr_database_replica

Provides a Vultr database replica resource. This can be used to create, read, modify, and delete managed database read replicas on your Vultr account.

## Example Usage

Create a new database replica:

```hcl
resource "vultr_database_replica" "my_database_replica" {
	database_id = vultr_database.my_database.id
	region = "sea"
    label = "my_database_replica_label"
	tag = "test tag"
}
```

## Argument Reference


~> Updating the database ID will cause a `force new`. This behavior is in place because the parent database cannot be changed without creating a new resource to serve as a read replica for the selected database.

The following arguments are supported:

* `database_id` - (Required) The managed database ID you want to attach this replica to.
* `region` - (Required) The ID of the region that the managed database read replica is to be created in. [See List Regions](https://www.vultr.com/api/#operation/list-regions)
* `label` - (Required) A label for the managed database read replica.
* `tag` - (Optional) The tag to assign to the managed database read replica.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the managed database read replica.
* `database_id` - The managed database ID.
* `date_created` - The date the managed database read replica was added to your Vultr account.
* `plan` - The managed database read replica's plan ID.
* `plan_disk` - The description of the disk(s) on the managed database read replica.
* `plan_ram` - The amount of memory available on the managed database read replica in MB.
* `plan_vcpus` - The number of virtual CPUs available on the managed database read replica.
* `plan_replicas` - The number of standby nodes available on the managed database read replica.
* `region` - The region ID of the managed database read replica.
* `status` - The current status of the managed database read replica (poweroff, rebuilding, rebalancing, configuring, running).
* `label` - The managed database read replica's label.
* `tag` - The managed database read replica's tag.
* `database_engine` - The database engine of the managed database read replica.
* `database_engine_version` - The database engine version of the managed database read replica.
* `vpc_id` - The ID of the VPC Network attached to the managed database read replica.
* `dbname` - The managed database read replica's default logical database.
* `host` - The hostname assigned to the managed database read replica.
* `public_host` - The public hostname assigned to the managed database read replica (VPC-attached only).
* `user` - The primary admin user for the managed database read replica.
* `password` - The password for the managed database read replica's primary admin user.
* `port` - The connection port for the managed database read replica.
* `maintenance_dow` - The preferred maintenance day of week for the managed database read replica.
* `maintenance_time` - The preferred maintenance time for the managed database read replica.
* `backup_hour` - The preferred hour of the day (UTC) for daily backups to take place (unavailable for Kafka engine types).
* `backup_minute` - The preferred minute of the backup hour for daily backups to take place (unavailable for Kafka engine types).
* `latest_backup` - The date of the latest backup available on the managed database read replica.
* `trusted_ips` - A list of allowed IP addresses for the managed database read replica.
* `mysql_sql_modes` - A list of SQL modes currently configured for the managed database read replica (MySQL engine types only).
* `mysql_require_primary_key` - The configuration value for whether primary keys are required on the managed database read replica (MySQL engine types only).
* `mysql_slow_query_log` - The configuration value for slow query logging on the managed database read replica (MySQL engine types only).
* `mysql_long_query_time` - The configuration value for the long query time (in seconds) on the managed database read replica (MySQL engine types only).
* `eviction_policy` - The configuration value for the data eviction policy on the managed database read replica (Valkey engine types only).
* `cluster_time_zone` - The configured time zone for the managed database read replica in TZ database format.
