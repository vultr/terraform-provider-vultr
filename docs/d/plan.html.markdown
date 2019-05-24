---
layout: "vultr"
page_title: "Vultr: vultr_plan"
sidebar_current: "docs-vultr-datasource-plan"
description: |-
  Get information about a Vultr plan.
---

# vultr_plan

Get information about a Vultr plan.

## Example Usage

Get the information for a plan by name:
```hcl
data "vultr_plan" "my_plan" {
  filter {
    name   = "name"
    values = ["8192 MB RAM,160 GB SSD,4.00 TB BW"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - Query parameters for finding plans.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the plan.
* `vcpu_count` - The number of virtual CPUs available on the plan.
* `ram` - The amount of memory available on the plan in MB.
* `disk` - The amount of disk space in GB available on the plan.
* `bandwidth` - The bandwidth available on the plan in TB.
* `bandwidth_gb` - The bandwidth available on the plan in GB.
* `price_per_month` - The price per month of the plan in USD.
* `plan_type` - The type of plan it is.
* `available_locations` - A list of DCIDs (used as `region_id` in Terraform) where the plan can be deployed.
* `deprecated` - Indicates that the plan will be going away in the future. New deployments of it will still be accepted, but you should begin to transition away from its usage.