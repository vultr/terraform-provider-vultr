---
layout: "vultr"
page_title: "Vultr: vultr_organization_group"
sidebar_current: "docs-vultr-datasource-organization-group"
description: |-
  Get information about organization groups.
---

# vultr_organization_group

Get information about organization groups.

## Example Usage

Get information about organization groups based on `name`:.

```hcl
data "vultr_organization_group" "my_group" {
  filter {
    name = "name"
    type = "my-group-name"
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

* `id` - The ID of the organization group.
* `name` - The name of the organization group.
* `description` - A description of the organization group.
* `users` - A list of users attached to the organization group.
* `roles` - A list of roles attached to the organization group.
* `policies` - A list of policies attached to the organization group.
* `date_created` - Date of creation for the organization group.
