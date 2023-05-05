---
layout: "vultr"
page_title: "Vultr: vultr_database_db"
sidebar_current: "docs-vultr-resource-database-db"
description: |-
  Provides a Vultr database DB resource. This can be used to create, read, and delete logical DBs for a managed database on your Vultr account.
---

# vultr_database_db

Provides a Vultr database DB resource. This can be used to create, read, and delete logical DBs for a managed database on your Vultr account.

## Example Usage

Create a new database DB:

```hcl
resource "vultr_database_db" "my_database_db" {
	database_id = vultr_database.my_database.id
	name = "my_database_db"
}
```

## Argument Reference


~> Updating the database ID or name will cause a `force new`. This behavior is in place because a database DB canno tbe moved from one managed database to another and logical databases can only be created/destroyed and not updated.

The following arguments are supported:

* `database_id` - (Required) The managed database ID you want to attach this logical DB to.
* `name` - (Required) The name of the new managed database logical DB.

## Attributes Reference

The following attributes are exported:

* `database_id` - The managed database ID.
* `name` - The name of the managed database logical DB.
