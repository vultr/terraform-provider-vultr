---
layout: "vultr"
page_title: "Vultr: vultr_snapshot_from_url"
sidebar_current: "docs-vultr-resource-snapshot_from_url"
description: |-
  Provides a Vultr Snapshots from URL resource. This can be used to create, read, modify, and delete Snapshots from URL.
---

# vultr_snapshot_from_url

Provides a Vultr Snapshots from URL resource. This can be used to create, read, modify, and delete Snapshots from URL.

## Example Usage

Create a new Snapshots from URL

```hcl
resource "vultr_snapshot_from_url" "my_snapshot" {
	url = "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.1-x86_64.iso"
}
```

## Argument Reference

The following arguments are supported:

* `url` - (Required) URL of the given resource you want to create a snapshot from.

## Attributes Reference

The following attributes are exported:

* `id` - The ID for the given snapshot.
* `description` - The description for the given snapshot.
* `url` - The url from where the raw image was used to create the snapshot.
* `date_created` - The date the snapshot was created.
* `size` - The size of the snapshot in Bytes.
* `status` - The status for the given snapshot.
* `os_id` - The os id which the snapshot is associated with.
* `app_id` - The app id which the snapshot is associated with.



## Import

Snapshots from Url can be imported using the Snapshot `ID`, e.g.

```
terraform import vultr_snapshot_from_url.my_snapshot e60dc0a2-9313-4bab-bffc-57ffe33d99f6
```
