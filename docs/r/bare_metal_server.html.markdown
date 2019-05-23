---
layout: "vultr"
page_title: "Vultr: vultr_bare_metal_server"
sidebar_current: "docs-vultr-resource-bare-metal-server"
description: |-
  Provides a Vultr Bare Metal Server resource. This can be used to create, modify, and delete Bare Metal Servers.
---

# vultr_bare_metal_server

Provides a Vultr Bare Metal Server resource. This can be used to create, modify, and delete Bare Metal Servers.

## Example Usage

Create a new Bare Metal Server
```hcl
resource "vultr_bare_metal_server" "my_baremetal" {
	
}
```

## Argument Reference

The following arguments are supported:

* `foo` - Description

## Attributes Reference

The following attributes are exported:

* `bar` - Description

## Import

Bare Metal Servers can be imported using the Server `SUBID`, e.g.

```
terraform import vultr_bare_metal_server.my_baremetal 954327
```