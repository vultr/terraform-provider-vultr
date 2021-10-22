---
layout: "vultr"
page_title: "Vultr: vultr_kubernetes_node_pools"
sidebar_current: "docs-vultr-resource-kubernetes-node-pools"
description: |-
Provides a Vultr Kubernetes Engine (VKE) Node Pool resource. This can be used to create, read, modify, and delete VKE clusters on your Vultr account.
---

# vultr_kubernetes_node_pools

Deploy additional node pools to an existing Vultr Kubernetes Engine (VKE) cluster.

## Example Usage

Create a new VKE cluster:

```hcl
resource "vultr_kubernetes_node_pools" "np-1" {
    cluster_id = vultr_kubernetes.k8.id
    node_quantity = 1
    plan = "vc2-2c-4gb"
    label = "my label"
    tag = "my tag"
}

```

## Argument Reference

The follow arguments are supported:

* `node_quantity` - (Required) The number of nodes in this node pool.
* `plan` - (Required) The plan to be used in this node pool. [See Plans List](https://www.vultr.com/api/#operation/list-plans) Note the minimum plan requirements must have at least 1 core and 2 gbs of memory.
* `label` - (Required) The label to be used as a prefix for nodes in this node pool.
* `tag` - (Optional) A tag that is assigned to this node pool.



## Attributes Reference

The following attributes are exported:
* `id` - The VKE cluster ID.
* `date_created` - Date of node pool creation.
* `date_updated` - Date of node pool updates.
* `label` - Label of node pool.
* `node_quantity` - Number of nodes within node pool.
* `plan` - Node plan that nodes are using within this node pool.
* `status` - Status of node pool.
* `tag` - Tag for node pool.
* `nodes` - Array that contains information about nodes within this node pool.

`nodes`

* `date_created` - Date node was created.
* `id` - ID of node.
* `label` - Label of node.
* `status` - Status of node.