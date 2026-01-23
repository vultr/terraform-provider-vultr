---
layout: "vultr"
page_title: "Vultr: vultr_organization_policy_group_attachment"
sidebar_current: "docs-vultr-resource-organization-policy-group-attachment"
description: |-
  Provides a Vultr organization policy group attachment resource. This can be used to attach an organization policy resource to a group resource on your Vultr account.
---

# vultr_organization_policy_group_attachment

Provides a Vultr organization policy group attachment resource. This can be used to attach an organization policy resource to a group resource on your Vultr account.

## Example Usage

Create a new organization policy group attachment.

```hcl
resource "vultr_organization_policy_group_attachment" "policy-group-attachment" {
  policy_id = "62948e0d-b09d-4d08-9554-3985b635cb17"  
  group_id = "2b72adbd-3f25-4ec8-b067-307a8fa4d8fa"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The ID of the policy.
* `group_id` - (Required) The ID of the group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the policy group attachment (a composite of the policy_id and group_id).
* `polciy_id` - The ID of the policy.
* `group_id` - The ID of the group.
