---
layout: "vultr"
page_title: "Vultr: vultr_organization_group"
sidebar_current: "docs-vultr-resource-organization-group"
description: |-
  Provides a Vultr organization group resource. This can be used to create, read, update and delete an organization group resource on your Vultr account.
---

# vultr_organization_group

Provides a Vultr organization group resource. This can be used to create, read, update and delete an organization group resource on your Vultr account.

## Example Usage

Create a new organization group.

```hcl
resource "vultr_organization_group" "my_group" {
  name = "my-group"
  description = "my group from terraform"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization group.
* `description` - (Required) A description of the organization group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the organization group.
* `name` - The name of the organization group.
* `description` - A description of the organization group.
* `users` - A list of users attached to the organization group.
* `roles` - A list of roles attached to the organization group.
* `policies` - **NOTE** This field is temporarily ignored until the API is fixed. A list of policies attached to the organization group.
* `date_created` - Date of creation for the organization group.

## Import

Organization groups can be imported using the `ID`, e.g.

```
terraform import vultr_organization_group.my_group eedbd991-3630-4edc-924b-1d9c38e82fe0
```
