---
layout: "vultr"
page_title: "Vultr: vultr_container_registry"
sidebar_current: "docs-vultr-resource-container-registry"
description: |-
  Provides a Vultr container registry resource. This can be used to create, read, modify, and delete registries on your Vultr account.
---

# vultr_container_registry

Create and update a Vultr container registry.

## Example Usage

Create a new container registry:

```hcl
resource "vultr_container_registry" "vcr1" {
  name = "examplecontainerregistry"
  region = "sjc"
  plan = "start_up"
  public = false
}
```

The `name` for container registries must be all lowercase and only contain alphanumeric characters.

## Argument Reference

The follow arguments are supported:

* `name` - (Required) The name for your container registry.  Must be lowercase and only alphanumeric characters.
* `region` - (Required) The region where your container registry will be deployed. [See available regions](https://www.vultr.com/api/#tag/Container-Registry/operation/list-registry-regions)
* `plan` - (Required) The billing plan for the container registry. [See available plans](https://www.vultr.com/api/#tag/Container-Registry/operation/list-registry-plans)
* `public` - (Required) Boolean indicating if the container registry should be created with public visibility or if it should require credentials.

## Attributes Reference

The following attributes are exported:
* `id` - The container registry ID.
* `name` - The name of the container registry.
* `plan` - The container registry plan.
* `region` - The region in which the container registry is deployed.
* `public` - Boolean indicating whether or not the requires login credentials.
* `urn` - The URN of the container registry.
* `storage` - A listing of current storage usage relevant to the container registry.
* `root_user` - The user associated with the container registry.
* `date_created` - A date-time denoting when the container registry was created.
