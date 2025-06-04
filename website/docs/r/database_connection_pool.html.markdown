---
layout: "vultr"
page_title: "Vultr: vultr_database_connection_pool"
sidebar_current: "docs-vultr-resource-database-connection-pool"
description: |-
  Provides a Vultr database connection pool resource. This can be used to create, read, modify, and delete connection pools for a PostgreSQL managed database on your Vultr account.
---

# vultr_database_connection_pool

Provides a Vultr database connection pool resource. This can be used to create, read, modify, and delete connection pools for a PostgreSQL managed database on your Vultr account.

## Example Usage

Create a new database connection pool:

```hcl
resource "vultr_database_connection_pool" "my_database_connection_pool" {
	database_id = vultr_database.my_database.id
	name = "my_database_connection_pool_name"
	database = "defaultdb"
	username = "vultradmin"
	mode = "transaction"
	size = 3
}
```

## Argument Reference


~> Updating the database ID or name will cause a `force new`. This behavior is in place because a database connection pool cannot be moved from one managed database to another and pool names cannot be updated.

The following arguments are supported:

* `database_id` - (Required) The managed database ID you want to attach this connection pool to.
* `name` - (Required) The name of the new managed database connection pool.
* `database` - (Required) The logical database to use for the new managed database connection pool.
* `username` - (Required) The database user to use for the new managed database connection pool.
* `mode` - (Required) The mode to configure for the new managed database connection pool (`session`, `transaction`, `statement`).
* `size` - (Required) The size of the new managed database connection pool.

## Attributes Reference

The following attributes are exported:

* `database_id` - The managed database ID.
* `name` - The name of the managed database connection pool.
* `database` - The logical database associated with the managed database connection pool.
* `username` - The database user associated with the managed database connection pool.
* `mode` - The configured mode of the managed database connection pool.
* `size` - The size of the managed database connection pool.
