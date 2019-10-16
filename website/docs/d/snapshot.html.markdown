---
layout: "vultr"
page_title: "Vultr: vultr_snapshot"
sidebar_current: "docs-vultr-datasource-snapshot"
description: |-
  Get information about a Vultr snapshot.
---

# vultr_snapshot

Get information about a Vultr snapshot.

## Example Usage

Get the information for a snapshot by `description`:

```hcl
data "vultr_snapshot" "my_snapshot" {
  filter {
    name = "description"
    values = ["my-snapshot-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding snapshots.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `description` - The description of the snapshot.
* `size` - The size of the snapshot in bytes.
* `status` - The status of the snapshot.
* `date_created` - The date the snapshot was added to your Vultr account.
* `os_id` - The operating system ID of the snapshot.
* `app_id` - The application ID of the snapshot.