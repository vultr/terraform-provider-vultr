---
layout: "vultr"
page_title: "Vultr: vultr_vpc"
sidebar_current: "docs-vultr-resource-vpc"
description: |-
  Provides a Vultr VPC resource. This can be used to create, read, and delete VPCs on your Vultr account.
---

# vultr_vpc

Provides a Vultr VPC resource. This can be used to create, read, and delete VPCs on your Vultr account.

## Example Usage

Create a new VPC with an automatically generated CIDR block:

```hcl
resource "vultr_vpc" "my_vpc" {
	description = "my vpc"
	region = "ewr"
}
```

Create a new VPC with a specified CIDR block:

```hcl
resource "vultr_vpc" "my_vpc" {
	description = "my private vpc"
	region = "ewr"
	v4_subnet  = "10.0.0.0"
	v4_subnet_mask = 24
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region ID that you want the VPC to be created in.
* `description` - (Optional) The description you want to give your VPC.
* `v4_subnet` - (Optional) The IPv4 subnet to be used when attaching instances to this VPC.
* `v4_subnet_mask` - The number of bits for the netmask in CIDR notation. Example: 32

## Attributes Reference

The following attributes are exported:

* `id` - ID of the VPC.
* `region` - The region ID that the VPC operates in.
* `description` - The description of the VPC.
* `v4_subnet` - The IPv4 subnet used when attaching instances to this VPC.
* `v4_subnet_mask` - The number of bits for the netmask in CIDR notation. Example: 32
* `date_created` - The date that the VPC was added to your Vultr account.

## Import

VPCs can be imported using the VPC `ID`, e.g.

```
terraform import vultr_vpc.my_vpc 0e04f918-575e-41cb-86f6-d729b354a5a1
```
