---
layout: "vultr"
page_title: "Vultr: vultr_user"
sidebar_current: "docs-vultr-datasource-user"
description: |-
  Get information about a Vultr user associated with your account.
---

# vultr_user

Get information about a Vultr user associated with your account. This data source provides the name, email, access control list, and API status for a Vultr user associated with your account.

## Example Usage

Get the information for a user by `email`:

```hcl
data "vultr_user" "my_user" {
  filter {
    name = "email"
    values = ["jdoe@example.com"]
  }
}
```

Get the information for a user by `name`:

```hcl
data "vultr_user" "my_user" {
  filter {
    name = "name"
    values = ["John Doe"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding users.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the user.
* `email` - The email of the user.
* `api_enabled` - Whether API is enabled for the user.
* `acl` - The access control list for the user.