---
layout: "vultr"
page_title: "Vultr: vultr_vpc2"
sidebar_current: "docs-vultr-datasource-vpc2"
description: |-
  Get information about a Vultr VPC 2.0.
---

** Deprecated: Use `vultr_vpc` instead **

# vultr_vpc2

Get information about a Vultr VPC 2.0.

## Example Usage

Get the information for a VPC 2.0 by `description`:

```hcl
data "vultr_vpc2" "my_vpc2" {
  filter {
    name = "description"
    values = ["my-vpc2-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding VPCs 2.0.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `region` - The ID of the region that the VPC 2.0 is in.
* `ip_block` - The IPv4 network address. For example: 10.1.1.0.
* `prefix_length` - The number of bits for the netmask in CIDR notation. Example: 20
* `description` - The VPC 2.0's description.
* `date_created` - The date the VPC 2.0 was added to your Vultr account.
