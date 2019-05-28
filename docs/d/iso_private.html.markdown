---
layout: "vultr"
page_title: "Vultr: vultr_iso_private"
sidebar_current: "docs-vultr-datasource-iso-private"
description: |-
  Get information about an ISO file uploaded to your Vultr account.
---

# vultr_iso_private

Get information about an ISO file uploaded to your Vultr account.

## Example Usage

Get the information for a ISO file by `filename`:
```hcl
data "vultr_iso_private" "my_iso" {
  filter {
    name   = "filename"
    values = ["my-iso-filename"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding ISO files.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `filename` - The ISO file's filename.
* `status` - The status of the ISO file.
* `size` - The size of the ISO file in bytes.
* `md5sum` - The md5 hash of the ISO file.
* `sha512sum` - The sha512 hash of the ISO file.
* `date_created` - The date the ISO file was added to your Vultr account.