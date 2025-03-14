---
layout: "vultr"
page_title: "Vultr: vultr_virtual_file_system_storage"
sidebar_current: "docs-vultr-datasource-virtual-file-system-storage"
description: |-
  Get information about a Vultr virtual file system storage subscription.
---

# vultr_virtual_file_system_storage

Get information about a Vultr virtual file system storage subscription.

## Example Usage

Get the information for a virtual file system storage subscription by `label`:

```hcl
data "vultr_virtual_file_system_storage" "my_vfs_storage" {
  filter {
    name = "label"
    values = ["my-vfs-storage-label"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding block storage subscriptions.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `attached_instances` - A list of instance IDs currently attached to the virtual file system storage.
* `attachments` - A list of attchment states for instances currently attached to the virtual file system storage.
* `charges` - The current pending charges for the virtual file system storage subscription in USD.
* `cost` - The cost per month of the virtual file system storage subscription in USD.
* `date_created` - The date the virtual file system storage subscription was added to your Vultr account.
* `disk_type` - The underlying disk type used by the virtual file system storage subscription.
* `label` - The label of the virtual file system storage subscription.
* `region` - The region ID of the virtual file system storage subscription.
* `status` - The status of the virtual file system storage subscription.
* `size_gb` - The size of the virtual file system storage subscription in GB.
* `tags` - A list of tags used on the virtual file system storage subscription.
