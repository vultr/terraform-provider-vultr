---
layout: "vultr"
page_title: "Vultr: vultr_reverse_ipv4"
sidebar_current: "docs-vultr-resource-reverse-ipv4"
description: |-
  Provides a Vultr Reverse IPv4 resource. This can be used to create, read, and modify reverse DNS records for IPv4 addresses.
---

# vultr_reverse_ipv4

Provides a Vultr Reverse IPv4 resource. This can be used to create, read, and
modify reverse DNS records for IPv4 addresses. Upon success, DNS
changes may take 6-12 hours to become active.

## Example Usage

Create a new reverse DNS record for an IPv4 address:

```hcl
resource "vultr_instance" "my_instance" {
	plan = "vc2-1c-1gb"
	region = "ewr"
	os_id = 167
	enable_ipv6 = true
}

resource "vultr_reverse_ipv4" "my_reverse_ipv4" {
	instance_id = "${vultr_instance.my_instance.id}"
	ip = "${vultr_instance.my_instance.main_ip}"
	reverse = "host.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the instance you want to set an IPv4
  reverse DNS record for.
* `ip` - (Required) The IPv4 address used in the reverse DNS record.
* `reverse` - (Required) The hostname used in the IPv4 reverse DNS record.

## Attributes Reference

The following attributes are exported:

* `id` - The ID is the IPv4 address in canonical format.
* `instance_id` - The ID of the instance the IPv4 reverse DNS record was set for.
* `ip` - The IPv4 address in canonical format used in the reverse DNS record.
* `gateway` - The gateway IP address.
* `netmask` - The IPv4 netmask in dot-decimal notation.
* `reverse` - The reverse DNS information for this IP address.