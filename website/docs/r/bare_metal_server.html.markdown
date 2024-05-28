---
layout: "vultr"
page_title: "Vultr: vultr_bare_metal_server"
sidebar_current: "docs-vultr-resource-bare-metal-server"
description: |-
  Provides a Vultr bare metal server resource. This can be used to create, read, modify, and delete bare metal servers on your Vultr account.
---

# vultr_bare_metal_server

Provides a Vultr bare metal server resource. This can be used to create, read, modify, and delete bare metal servers on your Vultr account.

## Example Usage

Create a new bare metal server:

```hcl
resource "vultr_bare_metal_server" "my_server" {
	plan = "vbm-4c-32gb"
	region = "ewr"
	os_id = 1743
}
```

Create a new bare metal server with options:

```hcl
resource "vultr_bare_metal_server" "my_server" {
	plan = "vbm-4c-32gb"
	region = "ewr"
	os_id = 1743
	label = "my-server-label"
	tags = ["my-server-tag"]
	hostname = "my-server-hostname"
	user_data = "this is my user data"
	enable_ipv6 = true
	activation_email = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The ID of the region that the server is to be created in. [See List Regions](https://www.vultr.com/api/#operation/list-regions)
* `plan` - (Required) The ID of the plan that you want the server to subscribe to. [See List Plans](https://www.vultr.com/api/#tag/plans)
* `os_id` - (Optional) The ID of the operating system to be installed on the server. [See List OS](https://www.vultr.com/api/#operation/list-os)
* `app_id` - (Optional) The ID of the Vultr application to be installed on the server. [See List Applications](https://www.vultr.com/api/#operation/list-applications)
* `image_id` - (Optional) The ID of the Vultr marketplace application to be installed on the server. [See List Applications](https://www.vultr.com/api/#operation/list-applications) Note marketplace applications are denoted by type: `marketplace` and you must use the `image_id` not the id.
* `snapshot_id` - (Optional) The ID of the Vultr snapshot that the server will restore for the initial installation. [See List Snapshots](https://www.vultr.com/api/#operation/list-snapshots)
* `script_id` - (Optional) The ID of the startup script you want added to the server.
* `vpc2_ids` - (Optional) A list of VPC 2.0 IDs to be attached to the server.
* `ssh_key_ids` - (Optional) A list of SSH key IDs to apply to the server on install (only valid for Linux/FreeBSD).
* `user_data` - (Optional) Generic data store, which some provisioning tools and cloud operating systems use as a configuration file. It is generally consumed only once after an instance has been launched, but individual needs may vary.
* `enable_ipv6` - (Optional) Whether the server has IPv6 networking activated.
* `activation_email` - (Optional) Whether an activation email will be sent when the server is ready.
* `hostname` - (Optional) The hostname to assign to the server.
* `tag` - (Deprecated: use `tags` instead) (Optional) The tag to assign to the server.
* `tags` - (Optional) A list of tags to apply to the servier.
* `label` - (Optional) A label for the server.
* `reserved_ipv4` - (Optional) The ID of the floating IP to use as the main IP of this server. [See Reserved IPs](https://www.vultr.com/api/#operation/list-reserved-ips)
* `mdisk_mode` - (Optional) The Raid configuration to apply to the server during provisioning. [See Bare Metal Create Request Body Schema](https://www.vultr.com/api/#tag/baremetal/operation/create-baremetal)
* `app_variables` - (Optional) A map of user-supplied variable keys and values for Vultr Marketplace apps. [See List Marketplace App Variables](https://www.vultr.com/api/#tag/marketplace/operation/list-marketplace-app-variables)

## Attributes Reference

The following attributes are exported:

* `id` - ID of the server.
* `region` - The ID of the region that the server is in.
* `os` - The string description of the operating system installed on the server.
* `ram` - The amount of memory available on the server in MB.
* `disk` - The description of the disk(s) on the server.
* `main_ip` - The server's main IP address.
* `cpu_count` - The number of CPUs available on the server.
* `default_password` - The server's default password.
* `date_created` - The date the server was added to your Vultr account.
* `netmask_v4` - The server's IPv4 netmask.
* `gateway_v4` - The server's IPv4 gateway.
* `status` - The status of the server's subscription.
* `v6_network` - The IPv6 subnet.
* `v6_main_ip` - The main IPv6 network address.
* `v6_network_size` - The IPv6 network size in bits.
* `plan` - The ID of the plan that server is subscribed to.
* `os_id` - The ID of the operating system installed on the server.
* `app_id` - The ID of the Vultr application installed on the server.
* `app_id` - The ID of the Vultr marketplace application installed on the server.
* `snapshot_id` - The ID of the Vultr snapshot that the server was restored from.
* `script_id` - The ID of the startup script that was added to the server.
* `vpc2_ids` - A list of VPC 2.0 IDs to be attached to the server.
* `ssh_key_ids` - A list of SSH key IDs applied to the server on install.
* `user_data` - Generic data store, which some provisioning tools and cloud operating systems use as a configuration file. It is generally consumed only once after an instance has been launched, but individual needs may vary.
* `enable_ipv6` - Whether the server has IPv6 networking activated.
* `activation_email` - Whether an activation email was sent when the server was ready.
* `hostname` - The hostname assigned to the server.
* `tag` - (Deprecated: use `tags` instead) The tag assigned to the server.
* `tags` - A list of tags applied to the server.
* `label` - A label for the server.
* `mac_address` - The MAC address associated with the server.


## Import

Bare Metal Servers can be imported using the server `ID`, e.g.

```
terraform import vultr_bare_metal_server.my_server b6a859c5-b299-49dd-8888-b1abbc517d08
```
