---
layout: "vultr"
page_title: "Vultr: vultr_dns_record"
sidebar_current: "docs-vultr-resource-dns-record"
description: |-
  Provides a Vultr DNS Record resource. This can be used to create, read, modify, and delete DNS Records.
---

# vultr_dns_record

Provides a Vultr DNS Record resource. This can be used to create, read, modify, and delete DNS Records.

## Example Usage

Create a new DNS Record
```hcl

resource "vultr_dns_domain" "my_domain" {
	domain = "domain.com"
	server_ip = "66.42.94.227"
}

resource "vultr_dns_record" "my_record" {
	domain = "${vultr_dns_domain.my_domain.id}"
	name = "www"
	data = "66.42.94.227"
	type = "A"
}
```

## Argument Reference

The following arguments are supported:

* `data` - (Required) IP Address of the server the domain is associated with.
* `domain` - (Required) Name of the DNS Domain this record will belong to.
* `name` - (Required) Name (subdomain) for this record.
* `type` - (Required) Type of record.
* `priority` - (Optional) Priority of this record (only required for MX and SRV). 
* `ttl` - (Optional) The time to live of this record.

## Attributes Reference

The following attributes are exported:

* `id` - ID associated with the record.
* `data` - IP Address of the server the domain is associated with.
* `domain` - Name of the DNS Domain this record will belong to.
* `name` - Name for this record (Can be subdomain).
* `type` - Type of record.
* `priority` - Priority of this record (only required for MX and SRV). 
* `ttl` -  The time to live of this record.