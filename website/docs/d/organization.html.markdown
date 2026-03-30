---
layout: "vultr"
page_title: "Vultr: vultr_organization"
sidebar_current: "docs-vultr-datasource-organization"
description: |-
  Get information about organizations.
---

# vultr_organization

Get information about organizations.

## Example Usage

Get information about organizations based on `name`:.

```hcl
data "vultr_organization" "my_org" {
  filter {
    name = "name"
    type = "my-org"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding operating systems.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the organization.
* `name` - The name of the organization.
* `type` - The type of organization.
* `date_created` - Date of creation for the organization.
