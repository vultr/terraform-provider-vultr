---
layout: "vultr"
page_title: "Vultr: vultr_network"
sidebar_current: "docs-vultr-datasource-network"
description: |-
  Get information about a Vultr private network.
---

# vultr_network

Get information about a Vultr private network.

## Example Usage

Get the information for a private network by `description`:
```hcl
data "vultr_network" "my_network" {
  filter {
    name = "description"
    values = ["my-network-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding private networks.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `region_id` - The ID of the region that the private network is in.
* `cidr_block` - The CIDR block of the private network.
* `description` - The private network's description.
* `date_created` - The date the private network was added to your Vultr account.