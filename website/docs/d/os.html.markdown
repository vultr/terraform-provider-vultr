---
layout: "vultr"
page_title: "Vultr: vultr_os"
sidebar_current: "docs-vultr-datasource-os"
description: |-
  Get information about operating systems that can be launched when creating a Vultr VPS.
---

# vultr_os

Get information about operating systems that can be launched when creating a Vultr VPS.

## Example Usage

Get the information for an operating system by `name`:

```hcl
data "vultr_os" "centos" {
  filter {
    name   = "name"
    values = ["CentOS 7 x64"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding operating systems.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the operating system.
* `arch` - The architecture of the operating system.
* `family` - The family of the operating system.
* `windows` - If true, a Windows license will be included with the instance, which will increase the cost.