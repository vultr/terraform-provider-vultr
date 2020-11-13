---
layout: "vultr"
page_title: "Vultr: vultr_bare_metal_plan"
sidebar_current: "docs-vultr-datasource-bare-metal-plan"
description: |-
  Get information about a Vultr bare metal server plan.
---

# vultr_bare_metal_plan

Get information about a Vultr bare metal server plan.

## Example Usage

Get the information for a plan by `id`:

```hcl
data "vultr_bare_metal_plan" "my_plan" {
  filter {
    name   = "id"
    values = ["vbm-4c-32gb"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding plans.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `cpu_count` - The number of CPUs available on the plan.
* `cpu_model` - The CPU model of the plan.
* `cpu_threads` - The number of CPU threads.
* `ram` - The amount of memory available on the plan in MB.
* `disk` - The description of the disk(s) on the plan.
* `bandwidth` - The bandwidth available on the plan.
* `monthly_cost` - The price per month of the plan in USD.
* `type` - The type of plan it is.
* `locations` - A list of DCIDs (used as `region` in Terraform) where the plan can be deployed.