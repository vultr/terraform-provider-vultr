---
layout: "vultr"
page_title: "Vultr: vultr_oidc_token"
sidebar_current: "docs-vultr-resource-oidc-token"
description: |-
  Provides a Vultr OIDC token resource.
---

# vultr_oidc_token

Provides a Vultr OIDC token resource.

## Example Usage

Create a new token.

```hcl
resource "vultr_oidc_token" "my_token" {
  grant_type = "external"
  client_id = "id"
  client_secret = "secret"
}
```

## Argument Reference

The following arguments are supported:

* `grant_type` - (Required) Must be `authorization_code` or `refresh_token`.
* `client_id` - (Required)
* `client_secret` - (Required)
* `code`
* `redirect_uri`
* `refresh_token`

## Attributes Reference

The following attributes are exported:

* `grant_type` - (Required) Must be `authorization_code` or `refresh_token`.
* `client_id` 
* `client_secret`
* `code`
* `redirect_uri`
* `refresh_token`
* `access_token`
* `token_type`
* `expires_seconds`
* `id_token`
* `scope`
