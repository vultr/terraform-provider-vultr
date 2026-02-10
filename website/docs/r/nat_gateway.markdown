---
layout: "vultr"
page_title: "Vultr: vultr_nat_gateway"
sidebar_current: "docs-vultr-resource-nat-gateway"
description: |-
  Provides a Vultr NAT Gateway resource. This can be used to create, read, modify, and delete NAT Gateways for a VPC network on your Vultr account.
---

# vultr_nat_gateway

Provides a Vultr NAT Gateway resource. This can be used to create, read, modify, and delete NAT Gateways for a VPC network on your Vultr account.

## Example Usage

Create a new NAT Gateway:

```hcl
resource "vultr_nat_gateway" "my_nat_gateway" {
	vpc_id = vultr_vpc.my_vpc.id
	label = "my_nat_gateway"
	tag = "my tag"
}
```

## Argument Reference

~> Updating the VPC ID will cause a `force new`. This behavior is in place because a NAT Gateway cannot be moved from one VPC to another. Addtionally, a VPC network can only have a single NAT Gateway attached to it.

The following arguments are supported:

* `vpc_id` - (Required) The VPC ID you want to attach this NAT Gateway to.
* `label` - (Optional) The label of the new NAT Gateway.
* `tag` - (Optional) The tag of the new NAT Gateway.

## Attributes Reference

The following attributes are exported:

* `vpc_id` - The VPC ID.
* `label` - The label of the NAT Gateway.
* `tag` - The tag of the NAT Gateway.
* `date_created` - The date the NAT Gateway was created.
* `status` - The status of the NAT Gateway.
* `public_ips` - The public IPv4 addresses of the NAT Gateway.
* `public_ips_v6` - The public IPv6 addresses of the NAT Gateway.
* `private_ips` - The private IP addresses of the NAT Gateway.
* `billing_charges` - The current charges for the NAT Gateway.
* `billing_monthly` - The total monthly charges for the NAT Gateway.
