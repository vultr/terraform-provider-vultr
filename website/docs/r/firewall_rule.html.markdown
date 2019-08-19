---
layout: "vultr"
page_title: "Vultr: vultr_firewall_rule"
sidebar_current: "docs-vultr-resource-firewall-rule"
description: |-
  Provides a Vultr Firewall Rule resource. This can be used to create, read, modify, and delete Firewall rules.
---

# vultr_firewall_rule

Provides a Vultr Firewall Rule resource. This can be used to create, read, modify, and delete Firewall rules.

## Example Usage

Create a Firewall Rule

```hcl
resource "vultr_firewall_group" "my_firewallgroup" {
    description = "base firewall"
}

resource "vultr_firewall_rule" "my_firewallrule" {
    firewall_group_id = "${vultr_firewall_group.my_firewallgroup.id}"
    protocol = "tcp"
    network = "0.0.0.0/0"
    from_port = "8085"
    to_port = "8090"
}
```

## Argument Reference

The following arguments are supported:

* `firewall_group_id` - (Required) The firewall group that the firewall rule will belong to.
* `protocol` - (Required) The type of protocol for this firewall rule. Possible values (icmp, tcp, udp, gre) **Note** they must be lowercase
* `network` - (Required) IP address that you want to define for this firewall rule.
* `from_port` - (Optional) Port that you want to define for this rule.
* `to_port` - (Optional) This can be used with the from port if you want to define multiple ports. Example from port 8085 to port 8090
* `notes` - (Optional) A simple note for a given firewall rule

## Attributes Reference

The following attributes are exported:

* `id` - The given ID for a firewall rule.
* `firewall_group_id` - The firewall group that the firewall rule belongs to.
* `protocol` - The type of protocol for this firewall rule. Possible values (icmp, tcp, udp, gre)
* `network` - IP address that is defined for this rule.
* `from_port` - Port that is defined for this rule.
* `to_port` - This can be used with the from port if you want to define multiple ports. Example from port 8085 to port 8090
* `notes` - A simple note for a given firewall rule
* `ip_type` - The type of ip this rule is - may be either v4 or v6.
