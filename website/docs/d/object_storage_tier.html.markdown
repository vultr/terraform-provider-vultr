---
layout: "vultr"
page_title: "Vultr: vultr_object_storage_tier"
sidebar_current: "docs-vultr-datasource-object_storage_tier"
description: |-
Get information about Object Storage tiers on Vultr.
---

# vultr_object_storage_tier

Get information about Object Storage tiers on Vultr.

## Example Usage

Get the information for an object storage tier by `slug`:

```hcl
data "vultr_object_storage_tier" "obs-tier" {
  filter {
    name   = "slug"
    values = ["tier_010k_5000m"]
  }
}
```

`slug` values and associated details can be retrieved through [this API call](https://www.vultr.com/api/#tag/s3/operation/list-object-storage-tiers).

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding operating systems.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `id` - The identifying tier ID.
* `slug` - The unique name for the tier.
* `price` - The monthly cost for the tier.
* `rate_limit_bytes` - The byte-per-second rate limit in the tier.
* `rate_limit_operations` - The operations-per-second rate limit in the tier.
* `locations` - A list of locations/clusters where the tier is available.
