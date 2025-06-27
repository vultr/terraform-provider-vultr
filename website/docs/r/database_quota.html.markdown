---
layout: "vultr"
page_title: "Vultr: vultr_database_quota"
sidebar_current: "docs-vultr-resource-database-quota"
description: |-
  Provides a Vultr database quota resource. This can be used to create, read, and delete quotas for a managed database on your Vultr account.
---

# vultr_database_quota

Provides a Vultr database quota resource. This can be used to create, read, modify, and delete quotas for a managed database on your Vultr account.

## Example Usage

Create a new database quota:

```hcl
resource "vultr_database_quota" "my_database_quota" {
	database_id = vultr_database.my_database.id
	client_id = "my_database_quota"
	consumer_byte_rate = "3"
	producer_byte_rate = "2"
	request_percentage = "120"
	user = "my_database_user"
}
```

## Argument Reference


~> Updating the database ID, client ID, or user will cause a `force new`. This behavior is in place because a database quota cannot be moved from one managed database to another and because quotas exist specifically for client ID/user pairs. 

The following arguments are supported:

* `database_id` - (Required) The managed database ID you want to attach this quota to.
* `client_id` - (Required) The client ID for the new database quota.
* `consumer_byte_rate` - (Required) The consumer byte rate for the new managed database quota.
* `producer_byte_rate` - (Required) The producer byte rate for the new managed database quota.
* `request_percentage` - (Required) The CPU request percentage for the new managed database quota.
* `user` - (Required) The user for the new managed database quota.

## Attributes Reference

The following attributes are exported:

* `database_id` - The managed database ID.
* `client_id` - The client ID for the new database quota.
* `consumer_byte_rate` - The consumer byte rate for the new managed database quota.
* `producer_byte_rate` - The producer byte rate for the new managed database quota.
* `request_percentage` - The CPU request percentage for the new managed database quota.
* `user` - The user for the new managed database quota.
