---
layout: "vultr"
page_title: "Vultr: vultr_instance"
sidebar_current: "docs-vultr-datasource-server"
description: |-
  Get information about a Vultr server.
---

# vultr_server

Get information about a Vultr server.

## Example Usage

Get the information for a server by `label`:

```hcl
data "vultr_server" "my_server" {
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
* `vcpu_count` - The number of virtual CPUs available on the server.
* `region` - The region ID of the server.
* `date_created` - The date the server was added to your Vultr account.
* `allowed_bandwidth` - The server's allowed bandwidth usage in GB.
* `netmask_v4` - The server's IPv4 netmask.
* `gateway_v4` - The server's IPv4 gateway.
* `status` - The status of the server's subscription.
* `power_status` - Whether the server is powered on or not.
* `server_status` - A more detailed server status (none, locked, installingbooting, isomounting, ok).
* `plan` - The server's plan ID.
* `v6_network` - The IPv6 subnet.
* `v6_main_ip` - The main IPv6 network address.
* `v6_network_size` - The IPv6 network size in bits.
* `label` - The server's label.
* `internal_ip` - The server's internal IP address.
* `kvm` - The server's current KVM URL. This URL will change periodically. It is not advised to cache this value.
* `tag` - The server's tag.
* `os_id` - The server's operating system ID.
* `app_id` - The server's application ID.
* `firewall_group_id` - The ID of the firewall group applied to this server.
* `features` - Array of which features are enabled.