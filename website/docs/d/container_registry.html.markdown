---
layout: "vultr"
page_title: "Vultr: vultr_container_registry"
sidebar_current: "docs-vultr-datasource-container-registry"
description: |-
Get information about a Vultr container registry resource. 
---

# vultr_container_registry

Get information about a Vultr container registry.

## Example Usage

```hcl
data "vultr_container_registry" "vcr-ds" {
  filter {
    name = "name"
    values = ["examplecontainerregistry"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding the container registry.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.


## Attributes Reference

The following attributes are exported:
* `id` - The container registry ID.
* `name` - The name of the container registry.
* `public` - Boolean indicating whether or not the requires login credentials.
* `urn` - The URN of the container registry.
* `storage` - A listing of current storage usage relevant to the container registry.
* `root_user` - The user associated with the container registry.
* `date_created` - A date-time denoting when the container registry was created.
* `repositories` - Listing of the repositories created within the registry and their metadata.

`storage`

* `used` - Amount of storage space used in gigabytes.
* `allowed` - Amount of storage space available in gigabytes.

`root_user`

* `id` - The ID of the root user.
* `username` - The username used to login as the root user.
* `password` - The password used to login as the root user.
* `date_created` - A date-time of when the root user was created.
* `date_modified` - A date-time of when the root user was last modified.

`repositories`

* `name` - The name of the repository.
* `image` - The image name in the repository.
* `description` - A description of the repo, if set.
* `date_create` - The date-time when the repository was created.
* `date_modified` - The date-time that the repository was last updated.
* `pull_count` - A count of the number of pulls against the repository.
* `artifact_count` - A count of the artifacts in the repository.
