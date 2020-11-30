---
layout: "vultr"
page_title: "Vultr: vultr_firewall_group"
sidebar_current: "docs-vultr-resource-firewall-group"
description: |-
  Provides a Vultr Firewall Group resource. This can be used to create, read, modify, and delete Firewall Group.
---

# vultr_firewall_group

Provides a Vultr Firewall Group resource. This can be used to create, read, modify, and delete Firewall Group.

## Example Usage

Create a new Firewall group

```hcl
resource "vultr_firewall_group" "my_firewallgroup" {
	description = "base firewall"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the firewall group.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the firewall group.
* `description` - Description of the firewall group.
* `date_created` - The date the firewall group was created.
* `date_modified` - The date the firewall group was modified.
* `instance_count` - The number of instances that are currently using this firewall group.
* `max_rule_count` - The number of max firewall rules this group can have.
* `rule_count` - The number of firewall rules this group currently has.

## Import

Firewall Groups can be imported using the Firewall Group `FIREWALLGROUPID`, e.g.

```
terraform import vultr_firewall_group.my_firewallgroup c342f929
```