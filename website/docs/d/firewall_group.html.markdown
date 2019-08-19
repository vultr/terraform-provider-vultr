---
layout: "vultr"
page_title: "Vultr: vultr_firewall_group"
sidebar_current: "docs-vultr-datasource-firewall-group"
description: |-
  Get information about a firewall group on your Vultr account.
---

# vultr_firewall_group

Get information about a firewall group on your Vultr account.

## Example Usage

Get the information for a firewall group by `description`:
```hcl
data "vultr_firewall_group" "my_fwg" {
  filter {
    name   = "description"
    values = ["fwg-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding firewall groups.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `description` - The description of the firewall group.
* `date_created` - The date the firewall group was added to your Vultr account.
* `date_modified` - The date the firewall group was last modified.
* `instance_count` - The number of instances this firewall group is applied to.
* `rule_count` - The number of rules added to this firewall group.
* `max_rule_count` - The maximum number of rules this firewall group can have.