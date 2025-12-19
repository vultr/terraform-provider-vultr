---
layout: "vultr"
page_title: "Vultr: vultr_organization_role_session"
sidebar_current: "docs-vultr-resource-organization-role-session"
description: |-
  Provides a Vultr organization role session resource. This can be used to create, read, update and delete an organization role session resource on your Vultr account.
---

# vultr_organization_role_session

Provides a Vultr organization role session resource. This can be used to create, read, update and delete an organization role session resource on your Vultr account.

## Example Usage

Create a new organization role session.

```hcl
resource "vultr_organization_role_session" "sss" {
  user = "4f86871e-fa75-4fb8-8960-fe4bac4e498c"
  role = "2c536798-6472-4907-b901-8b0c50c9789f"
  session_name = "my-terraform-session"
  type = "TemporaryAssumption"
  duration = 2300
  hour_start = 9
  hour_end = 17
  ip_address = "10.0.0.1"
  date_expires = "2025-12-31T23:59:59Z"
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Required) A UUID of a user.
* `role` - (Required) A UUID of a role.
* `type` - (Required) The type of role trust (ie. TemporaryAssumption).
* `session_name` - (Required) A name for the session.
* `duration` - (Required) The duration of the session.
* `ip_address` - (Required) An IP for the session.
* `hour_start` - 24h time that the trust is enabled.
* `hour_end` - 24h time that the trust is disabled.
* `date_expires` - (Required) ISO 8601 date time stamp indicationg expiration or role trust.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the session.
* `token` - The token for the session.
* `user` - The UUID of the session user.
* `role` - The UUID of the session role.
* `session_name` - The name of the session.
* `duration` - The duration of the session.
* `ip_address` - An IP for the session.
* `conditions_met` - A list of conditions met for the session.
* `remaining_duration` - The remaining duration of the session.
* `auth_method` - The auth method of the session.
* `source_ip` - The source IP of the session.
* `hour_start` - 24h time that the trust is enabled.
* `hour_end` - 24h time that the trust is disabled.
* `date_assumed` - The date that the session was assumed.
* `date_expires` - The date that the session expires.

## Import

Organization role sessions can be imported using the `ID`, e.g.

```
terraform import vultr_organization_role_session.my_session 9423eca6-0ecc-442b-a63c-d088722d6322
```
