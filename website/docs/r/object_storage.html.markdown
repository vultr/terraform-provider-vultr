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
    cluster_id = 9
    tier_id = 4
    label = "vultr-object-storage"

    bucket {
      name = "my-bucket"
      enable_versioning = true
      enable_lock = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the region that you want the object storage to be deployed in.
* `tier_id` - (Required) The ID of the tier to deploy the storage under.
* `label` - (Optional) The description you want to give your object storage.

* `bucket` (Optional) supports the following

**NOTE** Bucket support in Terraform relies solely on the local Terraform state.
If you change these values outside of Terraform it will not detect the state
drift. This will be added in a future version.

**NOTE** Updating any element of a bucket will necessitate destroying and
re-creating of the bucket

* `name` - (Required) A name for the bucket
* `enable_versioning` - (Optional) Whether or not the versioning is enabled in the bucket 
* `enable_lock` - (Optional) Whether or not object lock is enabled in the bucket 

## Attributes Reference

The following attributes are exported:

* `id` - The id of the object storage subscription.
* `label` - The label of the object storage subscription.
* `location` - The location which this subscription resides in.
* `cluster_id` - The identifying cluster ID.
* `region` - The region ID of the object storage subscription.
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
