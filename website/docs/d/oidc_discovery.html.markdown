---
layout: "vultr"
page_title: "Vultr: vultr_oidc_discovery"
sidebar_current: "docs-vultr-datasource-oidc-discovery"
description: |-
  Get OIDC discovery configuration information.
---

# vultr_oidc_discovery

Get OIDC discovery configuration information.

## Example Usage

Get information about oidc discover based on `provider_id`:.

```hcl
data "vultr_oidc_discovery" "my_oidc" {
  provider_id = "8f10eb02-f6b9-457a-bfa9-10f9ba6847b6"
}
```

## Argument Reference

The following arguments are supported:
* `provider_id` (Required) The ID of the provider to retrieve config.


## Attributes Reference

The following attributes are exported:

* `issuer`
* `authorize_endpoint`
* `token_endpoint`
* `jwks_uri` 
* `user_info_endpoint`
* `response_types_supported`
* `subject_types_supported`
* `id_token_values_supported`
* `scopes_supported`
* `claims_supported`
* `grant_types_supported`
* `token_endpoint_auth_methods_supported`
