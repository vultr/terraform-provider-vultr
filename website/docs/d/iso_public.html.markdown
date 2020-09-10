---
layout: "vultr"
page_title: "Vultr: vultr_iso_public"
sidebar_current: "docs-vultr-datasource-iso-public"
description: |-
  Get information about a public ISO file offered in the Vultr ISO library.
---

# vultr_iso_public

Get information about an ISO file offered in the Vultr ISO library.

## Example Usage

Get the information for a ISO file by `description`:

```hcl
data "vultr_iso_public" "my_iso" {
  filter {
    name   = "description"
    values = ["iso-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding ISO files.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `name` - The ISO file's name.
* `description` - The description of the ISO file.
* `md5sum` - The MD5Sum of the ISO file.