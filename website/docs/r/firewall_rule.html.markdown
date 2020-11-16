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
    firewall_group_id = vultr_firewall_group.my_firewallgroup.id
    protocol = "tcp"
    ip_type = "v4"
    subnet = "0.0.0.0"
    subnet_size = 0
    port = "8090"
    notes = "my firewall rule"
}
```

## Argument Reference

The following arguments are supported:

* `firewall_group_id` - (Required) The firewall group that the firewall rule will belong to.
* `protocol` - (Required) The type of protocol for this firewall rule. Possible values (icmp, tcp, udp, gre, esp, ah) **Note** they must be lowercase
* `ip_type` - (Required) The type of ip for this firewall rule. Possible values (v4, v6) **Note** they must be lowercase
* `subnet` - (Required) IP address that you want to define for this firewall rule.
* `subnet_size` - (Required) The number of bits for the subnet in CIDR notation. Example: 32.
* `port` - (Optional) TCP/UDP only. This field can be a specific port or a colon separated port range.
* `notes` - (Optional) A simple note for a given firewall rule
* `source` - (Optional) Possible values ("", cloudflare)

## Attributes Reference

The following attributes are exported:

* `id` - The given ID for a firewall rule.
* `firewall_group_id` - The firewall group that the firewall rule belongs to.
* `protocol` - The type of protocol for this firewall rule. Possible values (icmp, tcp, udp, gre, esp, ah)
* `network` - IP address that is defined for this rule.
* `port` - This field can be a specific port or a colon separated port range.
* `notes` - A simple note for a given firewall rule
* `ip_type` - The type of ip this rule is - may be either v4 or v6.

## Import

Firewall Rules can be imported using the Firewall Group `ID` and Firewall Rule `ID`, e.g.

```
terraform import vultr_firewall_rule.my_rule b6a859c5-b299-49dd-8888-b1abbc517d08,1
```