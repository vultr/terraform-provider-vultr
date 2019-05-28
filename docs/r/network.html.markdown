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
	region_id = 6
}
```

Create a new private network with a specified CIDR block:
```hcl
resource "vultr_network" "my_network" {
	description = "my private network"
	region_id = 6
	cidr_block  = "10.0.0.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Required) The region ID that you want the network to be created in.
* `description` - (Optional) The description you want to give your network.
* `cidr_block` - (Optional) The IPv4 subnet and subnet mask to be used when attaching servers to this network.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the network.
* `region_id` - The region ID that the network operates in.
* `description` - The description of the network.
* `cidr_block` - The IPv4 subnet and subnet mask to be used when attaching servers to this network.
* `date_created` - The date that the network was added to your Vultr account.

## Import

Networks can be imported using the network `NETWORKID`, e.g.

```
terraform import vultr_network.my_network net539626f0798d7
```