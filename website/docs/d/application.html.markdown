---
layout: "vultr"
page_title: "Vultr: vultr_application"
sidebar_current: "docs-vultr-datasource-application"
description: |-
  Get information about applications that can be launched when creating a Vultr VPS.
---

# vultr_application

Get information about applications that can be launched when creating a Vultr VPS.

## Example Usage

Get the information for an application by `deploy_name`:
```hcl
data "vultr_application" "docker" {
  filter {
    name   = "deploy_name"
    values = ["Docker on CentOS 7 x64"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding applications.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the application.
* `deploy_name` - The deploy name of the application.
* `short_name` - The short name of the application.