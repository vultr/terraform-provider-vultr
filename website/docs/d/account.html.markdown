---
layout: "vultr"
page_title: "Vultr: vultr_account"
sidebar_current: "docs-vultr-datasource-account"
description: |-
  Get information about your Vultr account.
---

# vultr_account

Get information about your Vultr account. This data source provides the balance, pending charges, last payment date, and last payment amount for your Vultr account.

## Example Usage

Get the information for an account:

```hcl
data "vultr_account" "my_account" {}
```

## Argument Reference

This data source does not take any arguments. It will return the account information associated with the Vultr API key you have set.

## Attributes Reference

The following attributes are exported:

* `balance` - The current balance on your Vultr account.
* `pending_charges` - The pending charges on your Vultr account.
* `last_payment_date` - The date of the last payment made on your Vultr account.
* `last_payment_amount` - The amount of the last payment made on your Vultr account.
