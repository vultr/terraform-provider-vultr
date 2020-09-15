---
layout: "vultr"
page_title: "Vultr: vultr_reserved_ip"
sidebar_current: "docs-vultr-datasource-reserved-ip"
description: |-
  Get information about a Vultr reserved IP address.
---

# vultr_reserved_ip

Get information about a Vultr reserved IP address.

## Example Usage

Get the information for a reserved IP by `label`:

```hcl
data "vultr_reserved_ip" "my_reserved_ip" {
  filter {
    name = "label"
    values = ["my-reserved-ip-label"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding reserved IP addresses.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `region` - The ID of the region that the reserved IP is in.
* `ip_type` - The IP type of the reserved IP.
* `subnet` - The subnet of the reserved IP.
* `subnet_size` - The subnet size of the reserved IP.
* `label` - The label of the reserved IP.
* `attached_to_vps` - The ID of the VPS the reserved IP is attached to.