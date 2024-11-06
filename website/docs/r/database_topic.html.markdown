---
layout: "vultr"
page_title: "Vultr: vultr_database_topic"
sidebar_current: "docs-vultr-resource-database-topic"
description: |-
  Provides a Vultr database topic resource. This can be used to create, read, modify, and delete topics for a managed database on your Vultr account.
---

# vultr_database_topic

Provides a Vultr database topic resource. This can be used to create, read, modify, and delete topics for a managed database on your Vultr account.

## Example Usage

Create a new database topic:

```hcl
resource "vultr_database_topic" "my_database_topic" {
	database_id = vultr_database.my_database.id
	name = "my_database_topic"
	partitions = "3"
	replication = "2"
	retention_hours = "120"
	retention_bytes = "-1"
}
```

## Argument Reference


~> Updating the database ID will cause a `force new`. This behavior is in place because a database topic cannot be moved from one managed database to another.

The following arguments are supported:

* `database_id` - (Required) The managed database ID you want to attach this topic to.
* `name` - (Required) The name for the new managed database topic.
* `partitions` - (Required) The number of partitions for the new managed database topic.
* `replication` - (Required) The replication factor for the new managed database topic.
* `retention_hours` - (Required) The retention hours for the new managed database topic.
* `retention_bytes` - (Required) The retention bytes for the new managed database topic.

## Attributes Reference

The following attributes are exported:

* `database_id` - The managed database ID.
* `name` - The name of the managed database topic.
* `partitions` - The number of partitions for the managed database topic.
* `replication` - The replication factor for the managed database topic.
* `retention_hours` - The retention hours for the managed database topic.
* `retention_bytes` - The retention bytes for the managed database topic.
