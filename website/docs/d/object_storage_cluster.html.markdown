---
layout: "vultr"
page_title: "Vultr: vultr_object_storage_cluster"
sidebar_current: "docs-vultr-datasource-object_storage_cluster"
description: |-
Get information about Object Storage Clusters on Vultr.
---

# vultr_object_storage_cluster

Get information about Object Storage Clusters on Vultr.

## Example Usage

Get the information for an object storage cluster by `region`:

```hcl
data "vultr_object_storage_cluster" "s3" {
  filter {
    name   = "region"
    values = ["ewr"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding operating systems.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `id` - The identifying cluster ID.
* `region` - The region ID of the object storage cluster.
* `hostname` - The cluster hostname.
* `deploy` - The Cluster is eligible for Object Storage deployment. (yes or no)


