---
layout: "vultr"
page_title: "Vultr: vultr_organization_role_trust"
sidebar_current: "docs-vultr-resource-organization-role-trust"
description: |-
  Provides a Vultr organization role trust resource. This can be used to create, read, update and delete an organization role trust resource on your Vultr account.
---

# vultr_organization_role_trust

Provides a Vultr organization role trust resource. This can be used to create, read, update and delete an organization role trust resource on your Vultr account.

## Example Usage

Create a new organization role trust.

```hcl
resource "vultr_organization_role_trust" "rot" {
  role_id = "de2e3d3b-c86f-402e-932f-368fe21c29df"
  user_id = "3b37fef9-3fb7-4803-94f6-ceace63b1b90"
  type = "TemporaryAssumption"
  hour_start = 9
  hour_end = 17
  ip_range = ["10.0.0.0/8", "192.168.0.1/16"]
  date_expires = "2025-12-31T23:59:59Z"
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) A UUID of a role.
* `user_id` - (Optional) A UUID of a user.
* `type` - (Required) The type of role trust.
* `hour_start` - (Required) The hour that the role trust begins.
* `hour_end` - (Required) The hour that the role trust ends.
* `ip_range` - (Required) A list of IP ranges allowed for the role trust.
* `date_expires` - (Required) ISO 8601 date time stamp indicating when to expire the role trust.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role trust.
* `role_id` - A UUID of a role.
* `user_id` - A UUID of a user.
* `type` - The type of role trust.
* `hour_start` - The hour that the role trust begins.
* `hour_end` - The hour that the role trust ends.
* `ip_range` - A list of IP ranges allowed for the role trust.
* `date_expires` - The date the role trust expires.
* `date_created` - The date the role trust was created.
