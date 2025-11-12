---
layout: "vultr"
page_title: "Vultr: vultr_organization_policy"
sidebar_current: "docs-vultr-datasource-organization-policy"
description: |-
  Get information about organization policies.
---

# vultr_organization_policy

Get information about organization policies.

## Example Usage

Get information about organization policies by `name`.

```hcl
data "vultr_organization_policy" "my_policy" {
  filter {
    name = "name"
    type = "my-policy-name"
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

* `id` - The ID of the organization policy.
* `name` - The name of the organization policy.
* `description` - A description of the organization policy.
* `is_system_policy` - Whether or not the policy is a system policy.
* `groups` - A list of group UUIDs attached to the organization policy.
* `users` - A list of user UUIDs attached to the organization policy.

* `document` - A block outlining the organization policy document details.
* `version` - A version for organization policy document.

* `statement` - A list of blocks for the organization policy statements.
* `effect` - The effect of the the policy document statement.
* `actions` - A list of actions for the policy document statement.
* `resources` - A list of applicable resources for the policy document statement.
