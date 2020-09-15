---
layout: "vultr"
page_title: "Vultr: vultr_server_ipv4"
sidebar_current: "docs-vultr-resource-server-ipv4"
description: |-
  Provides a server IPv4 resource. This can be used to create, read, and modify a IPv4 address.
---

# vultr_reverse_ipv4

Provides a Vultr server IPv4 resource. This can be used to create, read, and
modify a IPv4 address. Server is rebooted by default.

## Example Usage

Create a new IPv4 address for a server:

```hcl
resource "vultr_server" "my_server" {
	plan_id = "201"
	region_id = "6"
	os_id = "167"
	enable_ipv4 = true
}

resource "vultr_server_ipv4" "my_server_ipv4" {
	instance_id = "${vultr_server.my_server.id}"
	reboot = "false"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the server you want to set an IPv4
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
