---
layout: "vultr"
page_title: "Vultr: vultr_database_connector"
sidebar_current: "docs-vultr-resource-database-connector"
description: |-
  Provides a Vultr database connector resource. This can be used to create, read, modify, and delete connectors for a managed database on your Vultr account.
---

# vultr_database_connector

Provides a Vultr database connector resource. This can be used to create, read, modify, and delete connectors for a managed database on your Vultr account.

## Example Usage

Create a new database connector:

```hcl
resource "vultr_database_connector" "my_database_connector" {
	database_id = vultr_database.my_database.id
	name = "my_database_connector"
	class = "com.couchbase.connect.kafka.CouchbaseSinkConnector"
	topics = "my_database_topic"
	config = jsonencode({
		couchbase.seed.nodes = 3
		couchbase.username = "some_username"
		couchbase.password = "some_password"
	})
}
```

## Argument Reference


~> Updating the database ID, name, or class will cause a `force new`. This behavior is in place because a database connector cannot be moved from one managed database to another and connector names/classes cannot be updated.

The following arguments are supported:

* `database_id` - (Required) The managed database ID you want to attach this connector to.
* `name` - (Required) The name for the new managed database connector.
* `class` - (Required) The class for the new managed database connector.
* `topics` - (Required) A comma-separated list of topics to use with the new managed database connector.
* `config` - (Optional) A JSON string containing the configuration properties you wish to use with the new managed database connector.

## Attributes Reference

The following attributes are exported:

* `database_id` - The managed database ID.
* `name` - The name of the managed database connector.
* `class` - The class for the managed database connector.
* `topics` - A comma-separated list of topics to use with the managed database connector.
* `config` - A JSON string containing the configuration properties currently set for the managed database connector.
