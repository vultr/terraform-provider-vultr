---
layout: "vultr"
page_title: "Vultr: vultr_region"
sidebar_current: "docs-vultr-datasource-region"
description: |-
  Get information about a Vultr region.
---

# vultr_region

Get information about a Vultr region.

## Example Usage

Get the information for a region by `id`:

```hcl
data "vultr_region" "my_region" {
  filter {
    name   = "id"
    values = ["sea"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding regions.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `continent` - The continent the region is in.
* `country` - The country the region is in.
* `city` - The city the region is in.
* `options` - Shows whether options like ddos protection or block storage are available in the region.