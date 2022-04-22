---
layout: "vultr"
page_title: "Vultr: vultr_block_storage"
sidebar_current: "docs-vultr-resource-block-storage"
description: |-
  Provides a Vultr Block Storage resource. This can be used to create, read, modify, and delete Block Storage.
---

# vultr_block_storage

Provides a Vultr Block Storage resource. This can be used to create, read, modify, and delete Block Storage.

## Example Usage

Create a new Block Storage

```hcl
resource "vultr_block_storage" "my_blockstorage" {
	size_gb = 10
	region = "ewr"
}
```

## Argument Reference

~> Updating `block_type` will cause a `force new`.

The following arguments are supported:

* `size_gb` - (Required) The size of the given block storage.
* `region` - (Required) Region in which this block storage will reside in. (Currently only NJ/NY supported region "ewr")
* `attached_to_instance` - (Optional) VPS ID that you want to have this block storage attached to.
* `label` - (Optional) Label that is given to your block storage.
* `block_type` - (Optional)  Determines on the type of block storage volume that will be created. Soon to become a required parameter. Options are `high_per` or `storage_opt`.
* `live` - (Optional) Boolean value that will allow attachment of the volume to an instance without a restart. Default is false.



## Attributes Reference

The following attributes are exported:

* `id` - The ID for this block storage.
* `size_gb` - The size of the given block storage.
* `region` - Region in which this block storage will reside in. (Currently only NJ/NY supported region "ewr")
* `attached_to_instance` - VPS ID that is attached to this block storage.
* `label` - Label that is given to your block storage.
* `cost` - The monthly cost of this block storage.
* `date_created` - The date this block storage was created.
* `status` - Current status of your block storage.
* `live` - Flag which will determine if a volume should be attached with a restart or not.
* `mount_id` - An ID associated with the instance, when mounted the ID can be found in /dev/disk/by-id prefixed with virtio.
* `block_type` - The type of block storage volume. Values are `high_per` or `storage_opt`.

## Import

Block Storage can be imported using the Block Storage `ID`, e.g.

```
terraform import vultr_block_storage.my_blockstorage e315835e-d466-4e89-9b4c-dfd8788d7685
```
