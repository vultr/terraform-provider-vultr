---
layout: "vultr"
page_title: "Vultr: vultr_oidc_issuer"
sidebar_current: "docs-vultr-resource-oidc-issuer"
description: |-
  Provides a Vultr OIDC issuer resource.
---

# vultr_oidc_issuer

Provides a Vultr OIDC issuer resource.

## Example Usage

Create a new issuer.

```hcl
resource "vultr_oidc_issuer" "my_issuer" {
  name = "external"
  
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the issuer.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the issuer.
* `name` - The name of the organization.
* `issuer_id` - The issuer ID.
