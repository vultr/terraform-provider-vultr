---
layout: "vultr"
page_title: "Vultr: vultr_vpc"
sidebar_current: "docs-vultr-datasource-vpc"
description: |-
  Get information about a Vultr VPC.
---

# vultr_vpc

Get information about a Vultr VPC.

## Example Usage

Get the information for a VPC by `description`:

```hcl
data "vultr_vpc" "my_vpc" {
  filter {
    name = "description"
    values = ["my-vpc-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding VPCs.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `region` - The ID of the region that the VPC is in.
* `v4_subnet` - The IPv4 network address. For example: 10.1.1.0.
* `v4_subnet_mask` - The number of bits for the netmask in CIDR notation. Example: 20
* `description` - The VPC's description.
* `date_created` - The date the VPC was added to your Vultr account.
