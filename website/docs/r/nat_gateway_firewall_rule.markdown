---
layout: "vultr"
page_title: "Vultr: vultr_nat_gateway_firewall_rule"
sidebar_current: "docs-vultr-resource-nat-gateway-port-forwarding-rule"
description: |-
  Provides a Vultr NAT Gateway firewall rule resource. This can be used to create, read, modify, and delete firewall rules for a NAT Gateway on your Vultr account.
---

# vultr_nat_gateway_firewall_rule

Provides a Vultr NAT Gateway firewall rule resource. This can be used to create, read, modify, and delete firewall rules for a NAT Gateway on your Vultr account.

## Example Usage

Create a new firewall rule:

```hcl
resource "vultr_nat_gateway_firewall_rule" "my_firewall_rule" {
	vpc_id = vultr_nat_gateway.my_nat_gateway.vpc_id
	nat_gateway_id = vultr_nat_gateway.my_nat_gateway.id
	protocol = "tcp"
	subnet = "1.2.3.4"
	subnet_size = "24"
	port = "123"
	notes = "my notes"
}
```

## Argument Reference

~> Updating any field other than notes will cause a `force new`. This behavior is in place because a firewall rule's main properties cannot be updated after creation, but rules can be added and removed as needed.

~> NAT Gateway firewall rules depend on corresponding port forwarding rules for the same external port and protocol. Prior to creating a firewall rule, it is recommended to first create the necessary port forwarding rule to go along with it.

The following arguments are supported:

* `vpc_id` - (Required) The VPC ID associated with the NAT Gateway you want to attach this to.
* `nat_gateway_id` - (Required) The NAT Gateway ID you want to attach this firewall rule to.
* `protocol` - (Required) The protocol of the new firewall rule (`tcp`, `udp`).
* `subnet` - (Required) The subnet of the new firewall rule.
* `subnet_size` - (Required) The subnet size of the new firewall rule.
* `port` - (Required) The port or port range of the new firewall rule.
* `notes` - (Optional) The notes for the new firewall rule.

## Attributes Reference

The following attributes are exported:

* `vpc_id` - The VPC ID.
* `nat_gateway_id` - The NAT Gateway ID.
* `action` - The action of the firewall rule (always `accept`).
* `protocol` - The protocol of the firewall rule.
* `subnet` - The subnet of the firewall rule.
* `subnet_size` - The subnet size of the firewall rule.
* `port` - The  port of the firewall rule.
* `notes` - The notes for the firewall rule.

## Import

NAT gateway firewall rules rely on a VPC ID and a NAT gateway ID as well as their own ID.
To import we look up the rule via a composite key formatted like
`<VPC UUID>_<NAT UUID>_<rule UUID>` (the three UUIDs separated by underscores).
The import will look like:

```
terraform import vultr_nat_gateway_firewall_rule.my_rule 7ceab21c-e35b-4efb-9e96-6026ff5de5aa_f5dc2455-b748-4a8c-a8ac-5eef9c2ebfb7_0b5b8c4e-6c55-4b8b-9f2e-8a3d6d2f4c11
```
