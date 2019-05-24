---
layout: "vultr"
page_title: "Provider: Vultr"
sidebar_current: "docs-vultr-index"
description: |-
  The Vultr provider is used to interact with the resources supported by Vultr. The provider needs to be configured with the proper credentials before it can be used.
---

# Vultr Provider

The Vultr provider is used to interact with the
resources supported by [Vultr](https://www.vultr.com). The provider needs to be configured
with the proper credentials before it can be used.

## Example Usage

```hcl
# Configure the Vultr Provider
provider "vultr" {
  api_key = "VULTR_API_KEY"
}

# Create a web server
resource "vultr_server" "web" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `api_key` - (Required) This is the [Vultr API API](https://my.vultr.com/settings/#settingsapi). This can also be specified with the VULTR_API_KEY shell environment variable.
