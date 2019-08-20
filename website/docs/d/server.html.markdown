---
layout: "vultr"
page_title: "Vultr: vultr_server"
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

* `filter` - Query parameters for finding servers.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `os` - The operating system of the server.
* `ram` - The amount of memory available on the server in MB.
* `disk` - The description of the disk(s) on the server.
* `main_ip` - The server's main IP address.
* `vps_cpu_count` - The number of virtual CPUs available on the server.
* `location` - The physical location of the server.
* `region_id` - The region ID (`DCID` in the Vultr API) of the server.
* `default_password` - The server's default password.
* `date_created` - The date the server was added to your Vultr account.
* `pending_charges` - Charges pending for this server's subscription in USD.
* `cost_per_month` - The server's cost per month in USD.
* `current_bandwidth` - The server's current bandwidth usage in GB.
* `allowed_bandwidth` - The server's allowed bandwidth usage in GB.
* `netmask_v4` - The server's IPv4 netmask.
* `gateway_v4` - The server's IPv4 gateway.
* `status` - The status of the server's subscription.
* `power_status` - Whether the server is powered on or not.
* `server_state` - A more detailed server status (none, locked, installingbooting, isomounting, ok).
* `plan_id` - The server's plan ID.
* `v6_networks` - A list of the server's IPv6 networks.
* `label` - The server's label.
* `internal_ip` - The server's internal IP address.
* `kvm_url` - The server's current KVM URL. This URL will change periodically. It is not advised to cache this value.
* `auto_backups` - Whether auto backups are enabled on this server.
* `tag` - The server's tag.
* `os_id` - The server's operating system ID.
* `app_id` - The server's application ID.
* `firewall_group_id` - The ID of the firewall group applied to this server.