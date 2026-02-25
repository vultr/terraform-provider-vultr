---
layout: "vultr"
page_title: "Vultr: vultr_oidc_provider"
sidebar_current: "docs-vultr-datasource-oidc-provider"
description: |-
  Get OIDC provider information.
---

# vultr_oidc_provider

Get OIDC provider information.

## Example Usage

Get information about providers:

```hcl
data "vultr_oidc_provider" "my_provider" {
  filter {
    name   = "uri"
    values = ["my-uri"]
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

* `source`
* `uri`
* `n`
* `e`
* `alg`
* `use`
* `jwks_fetched_date`
* `jwks_expiry_date`
