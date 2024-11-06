---
layout: "vultr"
page_title: "Vultr: vultr_database"
sidebar_current: "docs-vultr-resource-database"
description: |-
  Provides a Vultr database resource. This can be used to create, read, modify, and delete managed databases on your Vultr account.
---

# vultr_database

Provides a Vultr database resource. This can be used to create, read, modify, and delete managed databases on your Vultr account.

## Example Usage

Create a new database:

```hcl
resource "vultr_database" "my_database" {
	database_engine = "pg"
	database_engine_version = "15"
    region = "ewr"
    plan = "vultr-dbaas-startup-cc-1-55-2"
    label = "my_database_label"
}
```

Create a new database with options:

```hcl
resource "vultr_database" "my_database" {
	database_engine = "pg"
	database_engine_version = "15"
    region = "ewr"
    plan = "vultr-dbaas-startup-cc-1-55-2"
    label = "my_database_label"
	tag = "some tag"
	cluster_time_zone = "America/New_York"
	maintenance_dow = "sunday"
	maintenance_time = "01:00"
}
```

## Argument Reference


~> Updating the database engine will cause a `force new`. This behavior is in place because databases cannot be changed from one type to another, although updating the database engine version is supported when an upgraded version is available.

The following arguments are supported:

* `region` - (Required) The ID of the region that the managed database is to be created in. [See List Regions](https://www.vultr.com/api/#operation/list-regions)
* `plan` - (Required) The ID of the plan that you want the managed database to subscribe to. [See List Managed Database Plans](https://www.vultr.com/api/#tag/managed-databases/operation/list-database-plans)
* `database_engine` - (Required) The database engine of the new managed database.
* `database_engine_version` - (Required) The database engine version of the new managed database.
* `label` - (Required) A label for the managed database.
* `vpc_id` - (Optional) The ID of the VPC Network to attach to the Managed Database.
* `tag` - (Optional) The tag to assign to the managed database.
* `maintenance_dow` - (Optional) The preferred maintenance day of week for the managed database.
* `maintenance_time` - (Optional) The preferred maintenance time for the managed database.
* `cluster_time_zone` - (Optional) The configured time zone for the Managed Database in TZ database format (e.g. `UTC`, `America/New_York`, `Europe/London`).
* `trusted_ips` - (Optional) A list of allowed IP addresses for the managed database.
* `maintenance_dow` - (Optional) The preferred maintenance day of week for the managed database.
* `maintenance_time` - (Optional) The preferred maintenance time for the managed database in 24-hour HH:00 format (e.g. `01:00`, `13:00`, `23:00`).
* `mysql_sql_modes` - (Optional) A list of SQL modes to configure for the managed database (MySQL engine types only - `ALLOW_INVALID_DATES`, `ANSI`, `ANSI_QUOTES`, `ERROR_FOR_DIVISION_BY_ZERO`, `HIGH_NOT_PRECEDENCE`, `IGNORE_SPACE`, `NO_AUTO_VALUE_ON_ZERO`, `NO_DIR_IN_CREATE`, `NO_ENGINE_SUBSTITUTION`, `NO_UNSIGNED_SUBTRACTION`, `NO_ZERO_DATE`, `NO_ZERO_IN_DATE`, `ONLY_FULL_GROUP_BY`, `PIPES_AS_CONCAT`, `REAL_AS_FLOAT`, `STRICT_ALL_TABLES`, `STRICT_TRANS_TABLES`, `TIME_TRUNCATE_FRACTIONAL`, `TRADITIONAL`).
* `mysql_require_primary_key` - (Optional) The configuration value for whether primary keys are required on the managed database (MySQL engine types only).
* `mysql_slow_query_log` - (Optional) The configuration value for slow query logging on the managed database (MySQL engine types only).
* `mysql_long_query_time` - (Optional) The configuration value for the long query time (in seconds) on the managed database (MySQL engine types only).
* `redis_eviction_policy` - (Optional) The configuration value for the data eviction policy on the managed database (Redis engine types only - `noeviction`, `allkeys-lru`, `volatile-lru`, `allkeys-random`, `volatile-random`, `volatile-ttl`, `volatile-lfu`, `allkeys-lfu`).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the managed database.
* `date_created` - The date the managed database was added to your Vultr account.
* `plan` - The managed database's plan ID.
* `plan_disk` - The description of the disk(s) on the managed database.
* `plan_ram` - The amount of memory available on the managed database in MB.
* `plan_vcpus` - The number of virtual CPUs available on the managed database.
* `plan_replicas` - The number of standby nodes available on the managed database (excluded for Kafka engine types).
* `plan_brokers` - The number of brokers available on the managed database (Kafka engine types only).
* `region` - The region ID of the managed database.
* `status` - The current status of the managed database (poweroff, rebuilding, rebalancing, configuring, running).
* `label` - The managed database's label.
* `tag` - The managed database's tag.
* `database_engine` - The database engine of the managed database.
* `database_engine_version` - The database engine version of the managed database.
* `vpc_id` - The ID of the VPC Network attached to the Managed Database.
* `dbname` - The managed database's default logical database.
* `ferretdb_credentials` - An associated list of FerretDB connection credentials (FerretDB + PostgreSQL engine types only).
* `host` - The hostname assigned to the managed database.
* `public_host` - The public hostname assigned to the managed database (VPC-attached only).
* `port` - The connection port for the managed database.
* `sasl_port` - The SASL connection port for the managed database (Kafka engine types only).
* `user` - The primary admin user for the managed database.
* `password` - The password for the managed database's primary admin user.
* `access_key` - The private key to authenticate the default user (Kafka engine types only).
* `access_cert` - The certificate to authenticate the default user (Kafka engine types only).
* `maintenance_dow` - The preferred maintenance day of week for the managed database.
* `maintenance_time` - The preferred maintenance time for the managed database.
* `latest_backup` - The date of the latest backup available on the managed database.
* `trusted_ips` - A list of allowed IP addresses for the managed database.
* `mysql_sql_modes` - A list of SQL modes currently configured for the managed database (MySQL engine types only).
* `mysql_require_primary_key` - The configuration value for whether primary keys are required on the managed database (MySQL engine types only).
* `mysql_slow_query_log` - The configuration value for slow query logging on the managed database (MySQL engine types only).
* `mysql_long_query_time` - The configuration value for the long query time (in seconds) on the managed database (MySQL engine types only).
* `redis_eviction_policy` - The configuration value for the data eviction policy on the managed database (Redis engine types only).
* `cluster_time_zone` - The configured time zone for the Managed Database in TZ database format.
* `read_replicas` - A list of read replicas attached to the managed database.


## Import

Database can be imported using the database `ID`, e.g.

```
terraform import vultr_database.my_database b6a859c5-b299-49dd-8888-b1abbc517d08
```
