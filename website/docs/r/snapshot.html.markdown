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
resource "vultr_instance" "my_instance" {
    label = "my_instance"
    region = "ewr"
    plan = 201
    os_id = 167
}
resource "vultr_snapshot" "my_snapshot" {
    instance_id       = "${vultr_instance.my_instance.id}"
    description  = "my instances snapshot"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of a given instance that you want to create a snapshot from.
* `description` - (Optional) The description for the given snapshot.

Snapshots often exceed the default timeout built in to all create requests in
the provider. In order to customize that, you may specify a custome value in a
`timeouts` block of the resource definition

* `timeouts` - (Optional) The create timeout value can be manually set

``` hcl
resource "vultr_snapshot" "sn" {
  instance_id = resource.vultr_instance.test-inst.id
  description = "terraform timeout test"

  timeouts {
    create = "60m"
  }
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID for the given snapshot.
* `instance_id` - The ID of the instance that the snapshot was created from.
* `description` - The description for the given snapshot.
* `date_created` - The date the snapshot was created.
* `size` - The size of the snapshot in Bytes.
* `status` - The status for the given snapshot.
* `os_id` - The os id which the snapshot is associated with.
* `app_id` - The app id which the snapshot is associated with.

## Import

Snapshots can be imported using the Snapshot `ID`, e.g.

```
terraform import vultr_snapshot_url.my_snapshot 283941e8-0783-410e-9540-71c86b833992
```
