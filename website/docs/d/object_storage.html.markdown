---
layout: "vultr"
page_title: "Vultr: vultr_object_storage"
sidebar_current: "docs-vultr-datasource-object_storage"
description: |-
  Get information about a Object Storage subscription on Vultr.
---

# vultr_object_storage

Get information about a Object Storage subscription on Vultr.

## Example Usage

Get the information for an object storage subscription by `label`:

```hcl
data "vultr_object_storage" "s3" {
  filter {
    name   = "label"
    values = ["my-s3"]
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

* `label` - The label of the object storage subscription.
* `location` - The location which this subscription resides in.
* `object_storage_cluster_id` - The identifying cluster ID.
* `region_id` - The region ID (DCID in the Vultr API) of the object storage subscription.
* `s3_access_key` - Your access key.
* `s3_hostname` - The hostname for this subscription.
* `s3_secret_key` - Your secret key.
* `status` - Current status of this object storage subscription.
* `date_created` - Date of creation for the object storage subscription.


