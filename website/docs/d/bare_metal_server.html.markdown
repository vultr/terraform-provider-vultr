---
layout: "vultr"
page_title: "Vultr: vultr_bare_metal_server"
sidebar_current: "docs-vultr-datasource-bare-metal"
description: |-
  Get information about a Vultr bare metal server.
---

# vultr_bare_metal_server

Get information about a Vultr bare metal server.

## Example Usage

Get the information for a server by `label`:

```hcl
data "vultr_bare_metal_server" "my_server" {
  filter {
    name   = "label"
    values = ["my-server-label"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding servers.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `os` - The operating system of the server.
* `ram` - The amount of memory available on the server in MB.
* `disk` - The description of the disk(s) on the server.
* `main_ip` - The server's main IP address.
* `cpu_count` - The number of CPUs available on the server.
* `location` - The location of the server.
* `region` - The region ID of the server.
* `default_password` - The server's default password.
* `date_created` - The date the server was added to your Vultr account.
* `status` - The status of the server's subscription.
* `netmask_v4` - The server's IPv4 netmask.
* `gateway_v4` - The server's IPv4 gateway.
* `plan` - The server's plan ID.
* `v6_networks` - A list of the server's IPv6 networks.
* `label` - The server's label.
* `tag` - The server's tag.
* `tags` - A list of tags applied to the server.
* `user_scheme` - The scheme used for the default user (linux servers only). 
* `os_id` - The server's operating system ID.
* `app_id` - The server's application ID.
* `image_id` - The Marketplace ID for this application.
* `vpc2_ids` - (Deprecated) list of VPC 2.0 IDs attached to the server.
