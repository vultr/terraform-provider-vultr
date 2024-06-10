---
layout: "vultr"
page_title: "Vultr: vultr_inference"
sidebar_current: "docs-vultr-datasource-inference"
description: |-
  Get information about a Vultr Serverless Inference subscription.
---

# vultr_inference

Get information about a Vultr Serverless Inference subscription.

## Example Usage

Get the information for an inference subscription by `label`:

```hcl
data "vultr_inference" "example_inference" {
  filter {
    name   = "label"
    values = ["my_inference_label"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding inference subscriptions.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `date_created` - The date the inference subscription was added to your Vultr account.
* `label` - The inference subscription's label.
* `api_key` - The inference subscription's API key for accessing the Vultr Inference API.
