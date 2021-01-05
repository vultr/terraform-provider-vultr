---
layout: "vultr"
page_title: "Vultr: vultr_instance_ipv4"
sidebar_current: "docs-vultr-datasource-instance-ipv4"
description: |-
  Get information about a Vultr instance IPv4.
---

# vultr_reverse_ipv4

Get information about a Vultr instance IPv4.

## Example Usage

Get the information for an IPv4 address by `instance_id`:

```hcl
data "vultr_instance_ipv4" "my_instance_ipv4" {
  filter {
    name = "ip"
    values = ["123.123.123.123"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding IPv4 address.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values to filter with.

## Attributes Reference

The following attributes are exported:

* `instance_id` - The ID of the instance the IPv4 address.
* `ip` - The IPv4 address in canonical format.
* `gateway` - The gateway IP address.
* `netmask` - The IPv4 netmask in dot-decimal notation.
* `reverse` - The reverse DNS information for this IP address.
