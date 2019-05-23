---
layout: "vultr"
page_title: "Vultr: vultr_account"
sidebar_current: "docs-vultr-datasource-account"
description: |-
  Get information about a Vultr account.
---

# vultr_account

Get information about a Vultr account. This data source provides the balance, pending charges, last payment date, and last payment amount for your Vultr account.

## Example Usage

Get the information for an account:
```hcl
data "vultr_account" "my_account" {
	
}
```

## Argument Reference

The following arguments are supported:

* `foo` - Description

## Attributes Reference

The following attributes are exported:

* `bar` - Description
