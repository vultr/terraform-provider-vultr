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

Get the information for a plan by `name`:

```hcl
data "vultr_bare_metal_plan" "my_plan" {
  filter {
    name   = "name"
    values = ["32768 MB RAM,4x 240 GB SSD,1.00 TB BW"]
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

* `name` - The name of the plan.
* `cpu_count` - The number of CPUs available on the plan.
* `cpu_model` - The CPU model of the plan.
* `ram` - The amount of memory available on the plan in MB.
* `disk` - The description of the disk(s) on the plan.
* `bandwidth_tb` - The bandwidth available on the plan in TB.
* `price_per_month` - The price per month of the plan in USD.
* `plan_type` - The type of plan it is.
* `available_locations` - A list of DCIDs (used as `region_id` in Terraform) where the plan can be deployed.
* `deprecated` - Indicates that the plan will be going away in the future. New deployments of it will still be accepted, but you should begin to transition away from its usage.