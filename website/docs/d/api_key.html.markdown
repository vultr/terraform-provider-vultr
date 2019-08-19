---
layout: "vultr"
page_title: "Vultr: vultr_api_key"
sidebar_current: "docs-vultr-datasource-api-key"
description: |-
  Get information about your Vultr API key.
---

# vultr_api_key

Get information about your Vultr API key. This data source provides the name, email, and access control list for your Vultr API key.

## Example Usage

Get the information for your API key:

```hcl
data "vultr_api_key" "my_api_key" {}
```

## Argument Reference

This data source does not take any arguments. It will return the API key information associated with the Vultr API key you have set.

## Attributes Reference

The following attributes are exported:

* `name` - The name associated with your Vultr API key.
* `email` - The email associated with your Vultr API key.
* `acl` - The access control list for your Vultr API key.
