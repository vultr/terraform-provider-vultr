---
layout: "vultr"
page_title: "Vultr: vultr_dns_domain"
sidebar_current: "docs-vultr-resource-dns-domain"
description: |-
  Provides a Vultr DNS Domain resource. This can be used to create, read, modify, and delete DNS Domains.
---

# vultr_dns_domain

Provides a Vultr DNS Domain resource. This can be used to create, read, modify, and delete DNS Domains.

## Example Usage

Create a new DNS Domain

```hcl
resource "vultr_dns_domain" "my_domain" {
	domain = "domain.com"
	ip = "66.42.94.227"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) Name of domain.
* `ip` - (Optional) instance IP you want associated to domain. If omitted this will create a domain with no records.

## Attributes Reference

The following attributes are exported:

* `id` - The ID is the name of the domain.
* `domain` -  Name of domain.
* `date_created` - The date the domain was added to your account.

## Import

DNS Domains can be imported using the Dns Domain `domain`, e.g.

```
terraform import vultr_dns_domain.name domain.com
```