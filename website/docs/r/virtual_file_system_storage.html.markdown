---
layout: "vultr"
page_title: "Vultr: vultr_virtual_file_system_storage"
sidebar_current: "docs-vultr-resource-virtual-file-system-storage"
description: |-
  Provides a Vultr virtual file system storage resource. This can be used to create, read, modify and delete a virtual file system storage.
---

# vultr_virtual_file_system_storage

Provides a Vultr virtual file system storage resource. This can be used to create, read, modify and delete a virtual file system storage.

## Example Usage

Define a virtual file system storage resource:

```hcl
resource "vultr_virtual_file_system_storage" "my_vfs_storage" {
  label = "vultr-vfs-storage"
  size_gb = 10
  region = "ewr"
  tags = ["terraform", "important"]
}
```

## Argument Reference

~> Updating `tags` will cause a `force new`.

The following arguments are supported:

* `size_gb` - (Required) The size of the given virtual file system storage subscription.
* `region` - (Required) The region in which this virtual file system storage will reside.
* `label` - (Required) The label to give to the virtual file system storage subscription.
* `tags` - (Optional) A list of tags to be used on the virtual file system storage subscription.
* `attached_instances` - (Optional) A list of UUIDs to attach to the virtual file system storage subscription.
* `disk_type` - (Optional) The underlying disk type to use for the virtual file system storage.  Default is `nvme`.

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

## Import

Virtual file system storage can be imported using the `ID`, e.g.

```
terraform import vultr_virtual_file_system_storage.my_vfs_storage 79210a84-bc58-494f-8dd1-953685654f7f
```
