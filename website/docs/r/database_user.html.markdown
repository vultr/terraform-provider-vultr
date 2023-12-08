---
layout: "vultr"
page_title: "Vultr: vultr_database_user"
sidebar_current: "docs-vultr-resource-database-user"
description: |-
  Provides a Vultr database user resource. This can be used to create, read, modify, and delete users for a managed database on your Vultr account.
---

# vultr_database_user

Provides a Vultr database user resource. This can be used to create, read, modify, and delete users for a managed database on your Vultr account.

## Example Usage

Create a new database user:

```hcl
resource "vultr_database_user" "my_database_user" {
	database_id = vultr_database.my_database.id
	username = "my_database_user"
	password = "randomTestPW40298"
}
```

## Argument Reference


~> Updating the database ID will cause a `force new`. This behavior is in place because a database user canno tbe moved from one managed database to another.

The following arguments are supported:

* `database_id` - (Required) The managed database ID you want to attach this user to.
* `username` - (Required) The username of the new managed database user.
* `password` - (Required) The password of the new managed database user.
* `encryption` - (Optional) The encryption type of the new managed database user's password (MySQL engine types only - `caching_sha2_password`, `mysql_native_password`).
* `access_control` - (Optional) The access control configuration for the new managed database user (Redis engine types only).

## Attributes Reference

The following attributes are exported:

* `database_id` - The managed database ID.
* `username` - The username of the managed database user.
* `password` - The password of the managed database user.
* `encryption` - The encryption type for the managed database user's password (MySQL engine types only).
* `access_control` - The access control configuration for the managed database user (Redis engine types only).
