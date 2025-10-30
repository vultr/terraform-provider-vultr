---
layout: "vultr"
page_title: "Vultr: vultr_organization_policy"
sidebar_current: "docs-vultr-resource-organization-policy"
description: |-
  Provides a Vultr organization policy resource. This can be used to create, read, update and delete an organization policy resource on your Vultr account.
---

# vultr_organization_policy

Provides a Vultr organization policy resource. This can be used to create, read, update and delete an organization policy resource on your Vultr account.

## Example Usage

Create a new organization policy.

```hcl
resource "vultr_organization_policy" "my_policy" {
  name = "my-policy"
  description = "my policy from terraform"
  is_system_policy = false
  groups = [ "048c1f4d-5201-4a05-8095-f3e30e878983" ]
  users = [ "265189eb-781b-46e2-88bb-9c6b829a8334" ]
  document {
    version = "2012-10-20"
    statement {
      effect = "Allow"
      action = ["compute.instance.Create"]
      resource = "*"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization policy.
* `description` - (Required) A description of the organization policy.
* `is_system_policy` - (Required) Whether or not the policy is a system policy.
* `groups` - (Optional) A list of group UUIDs to attach to the organization policy.
* `users` - (Optional) A list of user UUIDs to attach to the organization policy.

* `document` - (Required) A block outlining the organization policy document details.
* `version` - (Required) A version for organization policy document.

* `statement` - (Required) A list of blocks for the organization policy statements.
* `effect` - (Required) The effect of the the policy document statement.
* `action` - (Required) A list of actions for the policy document statement.
* `resource` - (Required) The applicable resources for the policy document statement.

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
* `action` - A list of actions for the policy document statement.
* `resource` - The applicable resources for the policy document statement.

## Import

Organization policies can be imported using the `ID`, e.g.

```
terraform import vultr_organization_policy.my_policy c54db579-15a1-4fcb-9ee0-46e97371ef9b
```
