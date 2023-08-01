---
layout: "vultr"
page_title: "Vultr: vultr_vpc2"
sidebar_current: "docs-vultr-resource-vpc2"
description: |-
  Provides a Vultr VPC 2.0 resource. This can be used to create, read, and delete VPCs 2.0 on your Vultr account.
---

# vultr_vpc2

Provides a Vultr VPC 2.0 resource. This can be used to create, read, and delete VPCs 2.0 on your Vultr account.

## Example Usage

Create a new VPC 2.0 with an automatically generated CIDR block:

```hcl
resource "vultr_vpc2" "my_vpc2" {
	description = "my vpc2"
	region = "ewr"
}
```

Create a new VPC 2.0 with a specified CIDR block:

```hcl
resource "vultr_vpc2" "my_vpc2" {
	description = "my private vpc2"
	region = "ewr"
	ip_block  = "10.0.0.0"
	prefix_length = 24
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region ID that you want the VPC 2.0 to be created in.
* `description` - (Optional) The description you want to give your VPC 2.0.
* `ip_type` - (Optional) Accepted values: `v4`.
* `ip_block` - (Optional) The IPv4 subnet to be used when attaching instances to this VPC 2.0.
* `prefix_length` - The number of bits for the netmask in CIDR notation. Example: 32

## Attributes Reference

The following attributes are exported:

* `id` - ID of the VPC 2.0.
* `region` - The region ID that the VPC 2.0 operates in.
* `description` - The description of the VPC 2.0.
* `ip_block` - The IPv4 subnet used when attaching instances to this VPC 2.0.
* `prefix_length` - The number of bits for the netmask in CIDR notation. Example: 32
* `date_created` - The date that the VPC 2.0 was added to your Vultr account.

## Import

VPCs 2.0 can be imported using the VPC 2.0 `ID`, e.g.

```
terraform import vultr_vpc2.my_vpc2 0e04f918-575e-41cb-86f6-d729b354a5a1
```
