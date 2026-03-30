---
layout: "vultr"
page_title: "Vultr: vultr_organization_role_group_attachment"
sidebar_current: "docs-vultr-resource-organization-role-group-attachment"
description: |-
  Provides a Vultr organization role group attachment resource. This can be used to attach an organization role resource to a group resource on your Vultr account.
---

# vultr_organization_role_group_attachment

Provides a Vultr organization role group attachment resource. This can be used to attach an organization role resource to a group resource on your Vultr account.

## Example Usage

Create a new organization role group attachment.

```hcl
resource "vultr_organization_role_group_attachment" "role-group-attachment" {
  role_id = "bf8587f2-72e6-43e5-9ebc-dd6267eb7f4c"  
  group_id = "2b72adbd-3f25-4ec8-b067-307a8fa4d8fa"
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) The ID of the role.
* `group_id` - (Required) The ID of the group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role group attachment (a composite of the role_id and group_id).
* `role_id` - The ID of the role.
* `group_id` - The ID of the group.
