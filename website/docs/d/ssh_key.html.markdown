---
layout: "vultr"
page_title: "Vultr: vultr_ssh_key"
sidebar_current: "docs-vultr-datasource-ssh-key"
description: |-
  Get information about a Vultr SSH key.
---

# vultr_ssh_key

Get information about a Vultr SSH key. This data source provides the name, public SSH key, and the creation date for your Vultr SSH key.

## Example Usage

Get the information for an SSH key by `name`:

```hcl
data "vultr_ssh_key" "my_ssh_key" {
  filter {
    name = "name"
    values = ["my-ssh-key-name"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding SSH keys.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the SSH key.
* `ssh_key` - The public SSH key.
* `date_created` - The date the SSH key was added to your Vultr account.
