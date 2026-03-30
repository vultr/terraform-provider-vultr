---
layout: "vultr"
page_title: "Vultr: vultr_organization_role_policy_attachment"
sidebar_current: "docs-vultr-resource-organization-role-polcy-attachment"
description: |-
  Provides a Vultr organization role policy attachment resource. This can be used to attach an organization role resource to a policy resource on your Vultr account.
---

# vultr_organization_role_policy_attachment

Provides a Vultr organization role policy attachment resource. This can be used to attach an organization role resource to a policy resource on your Vultr account.

## Example Usage

Create a new organization role policy attachment.

```hcl
resource "vultr_organization_role_policy_attachment" "role-policy-attachment" {
  role_id = "02883b84-9237-4abe-9572-db0cea335508"  
  policy_id = "1de3d445-c439-40e1-826e-12a445571d20"
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) The ID of the role.
* `policy_id` - (Required) The ID of the policy.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role policy attachment (a composite of the role_id and policy_id).
* `role_id` - The ID of the role.
* `policy_id` - The ID of the policy.
