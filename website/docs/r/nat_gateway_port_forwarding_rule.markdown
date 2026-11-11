---
layout: "vultr"
page_title: "Vultr: vultr_nat_gateway_port_forwarding_rule"
sidebar_current: "docs-vultr-resource-nat-gateway-port-forwarding-rule"
description: |-
  Provides a Vultr NAT Gateway port forwarding rule resource. This can be used to create, read, modify, and delete port forwarding rules for a NAT Gateway on your Vultr account.
---

# vultr_nat_gateway_port_forwarding_rule

Provides a Vultr NAT Gateway port forwarding rule resource. This can be used to create, read, modify, and delete port forwarding rules for a NAT Gateway on your Vultr account.

## Example Usage

Create a new port forwarding rule:

```hcl
resource "vultr_nat_gateway_port_forwarding_rule" "my_port_forwarding_rule" {
	vpc_id = vultr_nat_gateway.my_nat_gateway.vpc_id
	nat_gateway_id = vultr_nat_gateway.my_nat_gateway.id
	name = "my_port_forwarding_rule"
	protocol = "tcp"
	internal_port = "321"
	internal_ip = "10.1.2.3"
	external_port = "123"
	enabled = true
	description = "my description"
}
```

## Argument Reference

~> Updating the VPC or NAT Gateway ID will cause a `force new`. This behavior is in place because a port forwarding rule cannot be moved from one NAT Gateway to another.

~> NAT Gateway firewall rules depend on corresponding port forwarding rules for the same external port and protocol. Prior to deleting or changing a port forwarding rule, it is recommended to first remove any related firewall rules that may depend on it.

The following arguments are supported:

* `vpc_id` - (Required) The VPC ID associated with the NAT Gateway you want to attach this to.
* `nat_gateway_id` - (Required) The NAT Gateway ID you want to attach this port forwarding rule to.
* `name` - (Required) The name of the new port forwarding rule.
* `protocol` - (Required) The protocol of the new port forwarding rule (`tcp`, `udp`, `both`).
* `internal_port` - (Required) The internal port of the new port forwarding rule.
* `internal_ip` - (Required) The internal IP address of the new port forwarding rule.
* `external_port` - (Required) The external port of the new port forwarding rule.
* `enabled` - (Required) Whether the new port forwarding rule is enabled or disabled.
* `description` - (Optional) The description of the new port forwarding rule.

## Attributes Reference

The following attributes are exported:

* `vpc_id` - The VPC ID.
* `nat_gateway_id` - The NAT Gateway ID.
* `name` - The name of the port forwarding rule.
* `protocol` - The protocol of the port forwarding rule.
* `internal_port` - The internal port of the port forwarding rule.
* `internal_ip` - The internal IP address of the port forwarding rule.
* `external_port` - The external port of the port forwarding rule.
* `enabled` - Whether the port forwarding rule is enabled or disabled.
* `description` - The description of the port forwarding rule.
* `date_created` - The date the port forwarding rule was created.
* `date_updated` - The date the port forwarding rule was last updated.

## Import

NAT gateway port forwarding rules rely on a VPC ID and a NAT gateway ID as well as their own ID.
To import we look up the rule via a composite key formatted like
`<VPC UUID>_<NAT UUID>_<rule UUID>` (the three UUIDs separated by underscores).
The import will look like:

```
terraform import vultr_nat_gateway_port_forwarding_rule.my_rule 7ceab21c-e35b-4efb-9e96-6026ff5de5aa_f5dc2455-b748-4a8c-a8ac-5eef9c2ebfb7_0b5b8c4e-6c55-4b8b-9f2e-8a3d6d2f4c11
```
