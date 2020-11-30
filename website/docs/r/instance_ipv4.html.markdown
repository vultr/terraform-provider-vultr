---
layout: "vultr"
page_title: "Vultr: vultr_instance_ipv4"
sidebar_current: "docs-vultr-resource-instance-ipv4"
description: |-
  Provides a instance IPv4 resource. This can be used to create, read, and modify a IPv4 address.
---

# vultr_reverse_ipv4

Provides a Vultr instance IPv4 resource. This can be used to create, read, and
modify a IPv4 address. instance is rebooted by default.

## Example Usage

Create a new IPv4 address for a instance:

```hcl
resource "vultr_instance" "my_instance" {
	plan = "vc2-1c-1gb"
	region = "ewr"
	os_id = 167
	enable_ipv6 = true
}

resource "vultr_instance_ipv4" "my_instance_ipv4" {
	instance_id = "${vultr_instance.my_instance.id}"
	reboot = "false"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the instance you want to set an IPv4
  reverse DNS record for.
* `reboot` - (Optional) Default true. Determines whether or not the server is rebooted after adding the IPv4 address.

## Attributes Reference

The following attributes are exported:

* `id` - The ID is the IPv4 address in canonical format.
* `instance_id` - The ID of the server the IPv4 was set for.
* `ip` - The IPv4 address in canonical format.
* `gateway` - The gateway IP address.
* `netmask` - The IPv4 netmask in dot-decimal notation.
* `reverse` - The reverse DNS information for this IP address.
