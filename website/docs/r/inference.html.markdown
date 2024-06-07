---
layout: "vultr"
page_title: "Vultr: vultr_inference"
sidebar_current: "docs-vultr-resource-inference"
description: |-
  Provides a Vultr Serverless Inference resource. This can be used to create, read, modify, and delete inference subscriptions on your Vultr account.
---

# vultr_inference

Provides a Vultr Serverless Inference resource. This can be used to create, read, modify, and delete managed inference subscriptions on your Vultr account.

## Example Usage

Create a new inference subscription:

```hcl
resource "vultr_inference" "my_inference_subscription" {
    label = "my_inference_label"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) A label for the inference subscription.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the inference subscription.
* `date_created` - The date the inference subscription was added to your Vultr account.
* `label` - The inference subscription's label.
* `api_key` - The inference subscription's API key for accessing the Vultr Inference API.

## Import

Inference subscriptions can be imported using the subscription's `ID`, e.g.

```
terraform import vultr_inference.my_inference_subscription b6a859c5-b299-49dd-8888-b1abbc517d08
```
