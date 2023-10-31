---
layout: "vultr"
page_title: "Vultr: vultr_database"
sidebar_current: "docs-vultr-datasource-database"
description: |-
  Get information about a Vultr database.
---

# vultr_database

Get information about a Vultr database.

## Example Usage

Get the information for a database by `label`:

```hcl
data "vultr_database" "my_database" {
  filter {
    name   = "label"
    values = ["my-database-label"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding databases.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `date_created` - The date the managed database was added to your Vultr account.
* `plan` - The managed database's plan ID.
* `plan_disk` - The description of the disk(s) on the managed database.
* `plan_ram` - The amount of memory available on the managed database in MB.
* `plan_vcpus` - The number of virtual CPUs available on the managed database.
* `plan_replicas` - The number of standby nodes available on the managed database.
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
* `user` - The primary admin user for the managed database.
* `password` - The password for the managed database's primary admin user.
* `port` - The connection port for the managed database.
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
