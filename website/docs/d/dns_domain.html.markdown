---
layout: "vultr"
page_title: "Vultr: vultr_dns_domain"
sidebar_current: "docs-vultr-datasource-dns-domain"
description: |-
  Get information about a DNS domain associated with your Vultr account.
---

# vultr_dns_domain

Get information about a DNS domain associated with your Vultr account.

## Example Usage

Get the information for a DNS domain:

```hcl
data "vultr_dns_domain" "my_domain" {
  domain = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - The name you're searching for.

## Attributes Reference

The following attributes are exported:

* `domain` - Name of domain.
* `date_created` - The date the DNS domain was added to your Vultr account.