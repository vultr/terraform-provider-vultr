---
layout: "vultr"
page_title: "Vultr: vultr_user"
sidebar_current: "docs-vultr-resource-user"
description: |-
  Provides a Vultr User resource. This can be used to create, read, modify, and delete Users.
---

# vultr_user

Provides a Vultr User resource. This can be used to create, read, modify, and delete Users.

## Example Usage

Create a new User without any ACLs

```hcl
resource "vultr_user" "my_user" {
	name = "my user"
	email = "user@vultr.com"
	password = "myP@ssw0rd"
	api_enabled = true
}
```

Create a new User with all ACLs

```hcl
resource "vultr_user" "my_user" {
	name = "my user"
	email = "user@vultr.com"
	password = "myP@ssw0rd"
	api_enabled = true
	acl = [
	  "manage_users",
	  "subscriptions",
	  "provisioning",
	  "billing",
	  "support",
	  "abuse",
	  "dns",
	  "upgrade",
	]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name for this user.
* `email` - (Required) Email for this user.
* `password` - (Required) Password for this user.
* `api_enabled` - (Optional) Whether API is enabled for the user. Default behavior is set to enabled.
* `acl` - (Optional) The access control list for the user. 


## Attributes Reference

The following attributes are exported:
* `id` - ID associated with the user.
* `name` - Name for this user.
* `email` - Email for this user.
* `api_enabled` - Whether API is enabled for the user.

## Import

Users can be imported using the User `ID`, e.g.

```
terraform import vultr_user.myuser cbe5ced2ae716
```