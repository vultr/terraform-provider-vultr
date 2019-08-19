---
layout: "vultr"
page_title: "Vultr: vultr_startup_script"
sidebar_current: "docs-vultr-datasource-startup-script"
description: |-
  Get information about a Vultr startup script.
---

# vultr_startup_script

Get information about a Vultr startup script. This data source provides the name, script, type, creation date, and the last modification date for your Vultr startup script.

## Example Usage

Get the information for an startup script by `name`:
```hcl
data "vultr_startup_script" "my_startup_script" {
  filter {
    name = "name"
    values = ["my-startup-script-name"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding startup scripts.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the startup script.
* `script` - The contents of the startup script.
* `type` - The type of the startup script.
* `date_created` - The date the startup script was added to your Vultr account.
* `date_modified` - The date the startup script was last modified.
