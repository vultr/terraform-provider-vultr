---
layout: "vultr"
page_title: "Vultr: vultr_organization_policy_user_attachment"
sidebar_current: "docs-vultr-resource-organization-policy-user-attachment"
description: |-
  Provides a Vultr organization policy user attachment resource. This can be used to attach an organization policy resource to a user resource on your Vultr account.
---

# vultr_organization_policy_user_attachment

Provides a Vultr organization policy user attachment resource. This can be used to attach an organization policy resource to a user resource on your Vultr account.

## Example Usage

Create a new organization policy user attachment.

```hcl
resource "vultr_organization_policy_user_attachment" "policy-user-attachment" {
  policy_id = "62948e0d-b09d-4d08-9554-3985b635cb17"  
  user_id = "8aff6ef5-26f9-458b-ae87-2ac06cfe1d4b"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The ID of the policy.
* `user_id` - (Required) The ID of the user.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the policy user attachment (a composite of the policy_id and user_id).
* `polciy_id` - The ID of the policy.
* `user_id` - The ID of the user.
