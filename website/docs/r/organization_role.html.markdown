---
layout: "vultr"
page_title: "Vultr: vultr_organization_role"
sidebar_current: "docs-vultr-resource-organization-role"
description: |-
  Provides a Vultr organization role resource. This can be used to create, read, update and delete an organization role resource on your Vultr account.
---

# vultr_organization_role

Provides a Vultr organization role resource. This can be used to create, read, update and delete an organization role resource on your Vultr account.

## Example Usage

Create a new organization role.

```hcl
resource "vultr_organization_role" "my_role" {
  name = "my-role"
  description = "my role from terraform"
  type = "assumable"
  max_session_duration = 3600
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization role.
* `description` - (Required) A description of the organization role.
* `type` - (Required) A type for the organization role.
* `max_session_duration` - (Required) The max session length for the organization role.

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

## Import

Organization roles can be imported using the `ID`, e.g.

```
terraform import vultr_organization_role.my_role 378cbaa3-7e6b-4184-b71f-5d259618d811
```
