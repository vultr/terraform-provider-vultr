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
data "vultr_private_network" "my_network" {
  filter {
    name = "description"
    values = ["my-network-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding private networks.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `region` - The ID of the region that the private network is in.
* `v4_subnet` - The IPv4 network address. For example: 10.1.1.0.
* `v4_subnet_mask` - The number of bits for the netmask in CIDR notation. Example: 20
* `description` - The private network's description.
* `date_created` - The date the private network was added to your Vultr account.