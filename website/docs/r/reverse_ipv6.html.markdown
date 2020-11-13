---
layout: "vultr"
page_title: "Vultr: vultr_reverse_ipv6"
sidebar_current: "docs-vultr-resource-reverse-ipv6"
description: |-
  Provides a Vultr Reverse IPv6 resource. This can be used to create, read, modify, and delete reverse DNS records for IPv6 addresses.
---

# vultr_reverse_ipv6

Provides a Vultr Reverse IPv6 resource. This can be used to create, read,
modify, and delete reverse DNS records for IPv6 addresses. Upon success, DNS
changes may take 6-12 hours to become active.

## Example Usage

Create a new reverse DNS record for an IPv6 address:

```hcl
resource "vultr_server" "my_server" {
	plan = "vc2-1c-1gb"
	region = "ewr"
	os_id = 167
	enable_ipv6 = true
}

resource "vultr_reverse_ipv6" "my_reverse_ipv6" {
	instance_id = "${vultr_server.my_server.id}"
	ip = "${vultr_server.my_server.v6_networks[0].v6_main_ip}"
	reverse = "host.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the server you want to set an IPv6
  reverse DNS record for.
* `ip` - (Required) The IPv6 address used in the reverse DNS record.
* `reverse` - (Required) The hostname used in the IPv6 reverse DNS record.

## Attributes Reference

The following attributes are exported:

* `id` - The ID is the IPv6 address in canonical format.
* `instance_id` - The ID of the server the IPv6 reverse DNS record was set for.
* `ip` - The IPv6 address in canonical format used in the reverse DNS record.
* `reverse` - The hostname used in the IPv6 reverse DNS record.
