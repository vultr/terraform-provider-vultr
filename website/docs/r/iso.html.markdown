---
layout: "vultr"
page_title: "Vultr: vultr_iso_private"
sidebar_current: "docs-vultr-resource-iso"
description: |-
  Provides a Vultr ISO file resource. This can be used to create, read, and delete ISO files on your Vultr account.
---

# vultr_iso_private

Provides a Vultr ISO file resource. This can be used to create, read, and delete ISO files on your Vultr account.

## Example Usage

Create a new ISO

```hcl
resource "vultr_iso_private" "my_iso" {
	url = "http://dl-cdn.alpinelinux.org/alpine/v3.9/releases/x86_64/alpine-virt-3.9.3-x86_64.iso"
}
```

## Argument Reference

The following arguments are supported:

* `url` - (Required) URL pointing to the ISO file.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the ISO.
* `url` - URL pointing to the ISO file.
* `date_created` - The date the ISO was created.
* `filename` - The ISO filename.
* `size` - The ISO size in bytes.
* `md5sum` - The md5 hash of the ISO file.
* `sha512sum` - The sha512 hash of the ISO file.
* `status` - The status of the ISO file.

## Import

ISOs can be imported using the ISO `ISOID`, e.g.

```
terraform import vultr_iso_private.my_iso 2349859
```