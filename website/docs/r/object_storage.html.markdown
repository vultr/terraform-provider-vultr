---
layout: "vultr"
page_title: "Vultr: vultr_object_storage"
sidebar_current: "docs-vultr-resource-object_storage"
description: |-
  Provides a Vultr private object storage resource. This can be used to create, read, update and delete object storage resources on your Vultr account.
---

# vultr_object_storage

Provides a Vultr private object storage resource. This can be used to create, read, update and delete object storage resources on your Vultr account.

## Example Usage

Create a new object storage subscription.

```hcl
resource "vultr_object_storage" "tf" {
    object_storage_cluster_id = 2
    label = "tf-label"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The region ID that you want the network to be created in.
* `label` - (Optional) The description you want to give your network.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the object storage subscription.
* `label` - The label of the object storage subscription.
* `location` - The location which this subscription resides in.
* `cluster_id` - The identifying cluster ID.
* `region_id` - The region ID of the object storage subscription.
* `s3_access_key` - Your access key.
* `s3_hostname` - The hostname for this subscription.
* `s3_secret_key` - Your secret key.
* `status` - Current status of this object storage subscription.
* `date_created` - Date of creation for the object storage subscription.

## Import

Object Storage can be imported using the object storage `ID`, e.g.

```
terraform import vultr_object_storage.my_s3 0e04f918-575e-41cb-86f6-d729b354a5a1
```