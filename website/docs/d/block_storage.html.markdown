---
layout: "vultr"
page_title: "Vultr: vultr_block_storage"
sidebar_current: "docs-vultr-datasource-block-storage"
description: |-
  Get information about a Vultr block storage subscription.
---

# vultr_block_storage

Get information about a Vultr block storage subscription.

## Example Usage

Get the information for a block storage subscription by `label`:

```hcl
data "vultr_block_storage" "my_block_storage" {
  filter {
    name = "label"
    values = ["my-block-storage-label"]
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

* `label` - The label of the block storage subscription.
* `cost` - The cost per month of the block storage subscription in USD.
* `status` - The status of the block storage subscription.
* `size_gb` - The size of the block storage subscription in GB.
* `region` - The region ID of the block storage subscription.
* `attached_to_instance` - The ID of the VPS the block storage subscription is attached to.
* `date_created` - The date the block storage subscription was added to your Vultr account.