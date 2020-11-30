---
layout: "vultr"
page_title: "Vultr: vultr_network"
sidebar_current: "docs-vultr-resource-network"
description: |-
  Provides a Vultr private network resource. This can be used to create, read, and delete private networks on your Vultr account.
---

# vultr_network

Provides a Vultr private network resource. This can be used to create, read, and delete private networks on your Vultr account.

## Example Usage

Create a new private network with an automatically generated CIDR block:

```hcl
resource "vultr_network" "my_network" {
	description = "my private network"
	region = "ewr"
}
```

Create a new private network with a specified CIDR block:

```hcl
resource "vultr_network" "my_network" {
	description = "my private network"
	region = "ewr"
	subnet  = "10.0.0.0"
	subnet_size = 24
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region ID that you want the network to be created in.
* `description` - (Optional) The description you want to give your network.
* `v4_subnet` - (Optional) The IPv4 subnet to be used when attaching instances to this network.
* `v4_subnet_size` - The number of bits for the netmask in CIDR notation. Example: 32

## Attributes Reference

The following attributes are exported:

* `id` - ID of the network.
* `region` - The region ID that the network operates in.
* `description` - The description of the network.
* `v4_subnet` - The IPv4 subnet used when attaching instances to this network.
* `v4_subnet_size` - The number of bits for the netmask in CIDR notation. Example: 32
* `date_created` - The date that the network was added to your Vultr account.

## Import

Networks can be imported using the network `ID`, e.g.

```
terraform import vultr_network.my_network 0e04f918-575e-41cb-86f6-d729b354a5a1
```