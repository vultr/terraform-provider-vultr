---
layout: "vultr"
page_title: "Vultr: vultr_kubernetes_version"
sidebar_current: "docs-vultr-datasource-kubernetes-version"
description: |-
  Get information about the available kubernetes version
---

# vultr_kubernetes_version

Get information about the available kubernetes version.

## Example Usage

No parameters are required if you just want a list of versions. The version
strings are returned in full semantic versioning format:

```hcl
data "vultr_kubernetes_version" "my-k8s-versions" { }

output "my-latest-version" {
  value = data.vultr_kubernetes_version.my-k8s-versions.latest
}
```

If you want to filter on a minor version so that only current or newer versions
will return, you may pass a semantic versioning formatted string in to the
`filter` parameter:

```hcl
data "vultr_kubernetes_version" "my-k8s-new-versions" { 
  filter = "v1.36"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Semantic versioning formatted string which will filter out older k8s versions.

## Attributes Reference

The following attributes are exported:

* `available` - A list of all available kubernetes version, or, if filtered, matching or newer versions.
* `upgrades` - A list of all newer kubernetes version which are, at least, the next newest minor version or newer. Only present when a filter is provided.
* `latest` - The newest available kubernetes version.
