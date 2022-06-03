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

Get the information for a plan by `id`:

```hcl
data "vultr_plan" "my_plan" {
  filter {
    name   = "id"
    values = ["vc2-1c-1gb"]
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

* `vcpu_count` - The number of virtual CPUs available on the plan.
* `ram` - The amount of memory available on the plan in MB.
* `disk` - The amount of disk space in GB available on the plan.
* `bandwidth` - The bandwidth available on the plan in GB.
* `monthly_cost` - The price per month of the plan in USD.
* `type` - The type of plan it is.
* `gpu_vram` - For GPU plans, the VRAM available in the plan.
* `gpu_type` - For GPU plans, the GPU card type.
* `locations` - A list of DCIDs (used as `region` in Terraform) where the plan can be deployed.
* `disk_count` - The number of disks that this plan offers.
