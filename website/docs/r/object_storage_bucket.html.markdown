---
layout: "vultr"
page_title: "Vultr: vultr_object_storage_bucket"
sidebar_current: "docs-vultr-resource-object_storage_bucket"
description: |-
  Provides a Vultr private object storage bucket resource. This can be used to create, read and delete object storage bucket resources on your Vultr account.
---

# vultr_object_storage_bucket

Provides a Vultr private object storage bucket resource. This can be used to create, read and delete object storage bucket resources on your Vultr account.

## Example Usage

Create a new object storage bucket subscription.

```hcl
resource "vultr_object_storage_bucket" "obj-bucket" {
    object_storage_id = "0d17ed61-f852-4c26-9e39-aeba6381847d"
    name = "my-bucket"
    enable_versioning = true
    enable_lock = true
}
```

## Argument Reference

The following arguments are supported:

**NOTE** Updating any element of a bucket will necessitate destroying and
re-creating of the bucket

* `object_storage_id` - (Required) The ID of the object storage that you want create the bucket in.
* `name` - (Required) The name of the bucket.
* `enable_versioning` - (Optional) Whether or not the versioning is enabled in the bucket 
* `enable_lock` - (Optional) Whether or not object lock is enabled in the bucket 

## Attributes Reference

The following attributes are exported:

* `object_storage_id` - The ID of the object storage that you want create the bucket in.
* `name` - The name of the bucket.
* `enable_versioning` - Whether or not the versioning is enabled in the bucket.
* `enable_lock` - Whether or not object lock is enabled in the bucket.
* `date_created` - Date of creation for the object storage subscription.
