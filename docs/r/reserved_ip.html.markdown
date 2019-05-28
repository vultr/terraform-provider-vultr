---
layout: "vultr"
page_title: "Vultr: vultr_reserved_ip"
sidebar_current: "docs-vultr-resource-reserved-ip"
description: |-
  Provides a Vultr reserved IP resource. This can be used to create, read, modify, and delete reserved IP addresses on your Vultr account.
---

# vultr_reserved_ip

Provides a Vultr reserved IP resource. This can be used to create, read, modify, and delete reserved IP addresses on your Vultr account.

## Example Usage

Create a new reserved IP:
```hcl
resource "vultr_reserved_ip" "my_reserved_ip" {
	label = "my-reserved-ip"
	region_id = 6
	ip_type = "v4"
}
```

Attach a reserved IP to a server:
```hcl
resource "vultr_reserved_ip" "my_reserved_ip" {
	label = "my-reserved-ip"
	region_id = 6
	ip_type = "v4"
	attached_id = "923483"
}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Required) The region ID that you want the reserved IP to be created in.
* `ip_type` - (Required) The type of reserved IP that you want. Either "v4" or "v6".
* `label` - (Optional) The label you want to give your reserved IP.
* `attached_id` - (Optional) The VPS ID you want this reserved IP to be attached to.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the reserved IP.
* `region_id` - The region ID (`DCID` in the Vultr API) that this reserved IP belongs to.
* `ip_type` - The reserved IP's type.
* `label` - The reserved IP's label.
* `attached_id` - The ID of the server the reserved IP is attached to.
* `subnet` - The reserved IP's subnet.
* `subnet_size` - The reserved IP's subnet size.

## Import

Reserved IPs can be imported using the reserved IP `SUBID`, e.g.

```
terraform import vultr_reserved_ip.my_reserved_ip 1313044
```