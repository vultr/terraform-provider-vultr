---
layout: "vultr"
page_title: "Vultr: vultr_snapshot"
sidebar_current: "docs-vultr-resource-snapshot"
description: |-
  Provides a Vultr Snapshot resource. This can be used to create, read, modify, and delete Snapshot.
---

# vultr_snapshot

Provides a Vultr Snapshot resource. This can be used to create, read, modify, and delete Snapshot.

## Example Usage

Create a new Snapshot
```hcl
resource "vultr_server" "my_server" {
    label = "my_server"
    region_id = "1"
    plan_id = 201
    os_id = 147
}
resource "vultr_snapshot" "my_snapshot" {
    vps_id       = "${vultr_server.snap.id}"
    description  = "my servers snapshot"
}
```

## Argument Reference

The following arguments are supported:

* `vps_id` - (Required) ID of a given server that you want to create a snapshot from.
* `description` - (Optional) The description for the given snapshot.

## Attributes Reference

The following attributes are exported:

* `id` - The ID for the given snapshot.
* `vps_id` - The ID of the server that the snapshot was created from.
* `description` - The description for the given snapshot.
* `date_created` - The date the snapshot was created.
* `size` - The size of the snapshot in Bytes.
* `status` - The status for the given snapshot.
* `os_id` - The os id which the snapshot is associated with.
* `app_id` - The app id which the snapshot is associated with.

## Import

Snapshots can be imported using the Snapshot `SNAPSHOTID`, e.g.

```
terraform import vultr_snapshot_url.my_snapshot 9735ced831ed2
```