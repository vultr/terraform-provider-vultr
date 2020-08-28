---
layout: "vultr"
page_title: "Vultr: vultr_reverse_ipv4"
sidebar_current: "docs-vultr-datasource-reverse-ipv4"
description: |-
  Get information about a Vultr Reverse IPv4.
---

# vultr_reverse_ipv4

Get information about a Vultr Reverse IPv4.

## Example Usage

Get the information for an IPv4 reverse DNS record by `reverse`:

```hcl
data "vultr_reverse_ipv4" "my_reverse_ipv4" {
  filter {
    name = "reverse"
    values = ["host.example.com"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding IPv4 reverse DNS records.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values to filter with.

## Attributes Reference

The following attributes are exported:

* `instance_id` - The ID of the server the IPv4 reverse DNS record was set for.
* `ip` - The IPv4 address in canonical format used in the reverse DNS record.
* `reverse` - The hostname used in the IPv4 reverse DNS record.
