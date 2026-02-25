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

Create a new external issuer:
```hcl
resource "vultr_oidc_issuer" "my_issuer" {
  source = "external"
  uri = "https://auth.example.com"
  kid = "key-001"
  kty = "RSA"
  n = "xGOr-H7..."
  e = "AQAB"
  alg = "RS256"
  use = "sig"
}
```

Create a new VKE issuer:
```hcl
resource "vultr_oidc_issuer" "my_vke_issuer" {
  source = "external"
  source_id = "a795b8c3-a199-4cf9-a076-fba37a1e0934"
}
```

## Argument Reference

The following arguments are supported:

* `source` - (Required) The source of the issuer. Must be either `vke` or `external`.
* `source_id` - (Optional) UUID of the VKE to source. Required for `source` of `vke`

* `uri` - (Optional) Required for `source` of `external`.
* `kid` - (Optional) Required for `source` of `external`.
* `kty` - (Optional) Required for `source` of `external`.
* `n` - (Optional) Required for `source` of `external`.
* `e` - (Optional) Required for `source` of `external`.
* `alg` - (Optional) Required for `source` of `external`.
* `use` - (Optional) Required for `source` of `external`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the issuer.
* `source` 
* `source_id` 
* `uri`
* `kid`
* `kty`
* `n`
* `e`
* `alg`
* `use`
