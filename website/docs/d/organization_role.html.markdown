---
layout: "vultr"
page_title: "Vultr: vultr_organization_role"
sidebar_current: "docs-vultr-datasource-organization-role"
description: |-
  Get information about organization roles.
---

# vultr_organization_role

Get information about organization roles.

## Example Usage

Get information about organization roles by `name`:

```hcl
data "vultr_organization_role" "my_role" {
  filter {
    name = "name"
    type = "my-role-name"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding operating systems.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the organization role.
* `name` - The name of the organization role.
* `description` - A description of the organization role.
* `type` - A type for the organization role.
* `max_session_duration` - The max session length for the organization role.
* `policies` - A list of UUIDs attached to the role.
* `groups` - A list of group UUIDs attached to the role.
* `date_created` - Date of creation for the organization role.
