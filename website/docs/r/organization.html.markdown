---
layout: "vultr"
page_title: "Vultr: vultr_organization"
sidebar_current: "docs-vultr-resource-organization"
description: |-
  Provides a Vultr organization resource. This can be used to create, read, update and delete an organization resource on your Vultr account.
---

# vultr_organization

Provides a Vultr organization resource. This can be used to create, read, update and delete an organization resource on your Vultr account.

## Example Usage

Create a new organization.

```hcl
resource "vultr_organization" "my_org" {
  name = "my-org"
  type = "Personal/Team"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization.
* `type` - (Required) The type of organization.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the organization.
* `name` - The name of the organization.
* `type` - The type of organization.
* `date_created` - Date of creation for the organization.

## Import

Organizations can be imported using the `ID`, e.g.

```
terraform import vultr_organization.my_org 0e04f918-575e-41cb-86f6-d729b354a5a1
```
