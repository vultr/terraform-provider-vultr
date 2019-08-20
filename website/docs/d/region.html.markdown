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

Get the information for a region by `name`:

```hcl
data "vultr_region" "my_region" {
  filter {
    name   = "name"
    values = ["Miami"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding regions.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the region.
* `continent` - The continent the region is in.
* `country` - The country the region is in.
* `state` - The state the region is in.
* `ddos_protection` - Whether the region has DDoS protection.
* `block_storage` - Whether the region has block storage.
* `regioncode` - The region code of the region.
