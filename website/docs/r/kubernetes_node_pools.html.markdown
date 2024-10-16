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
    label = "my-label"
    tag = "my-tag"
	auto_scaler = true
	min_nodes = 1
	max_nodes = 2
}

```

## Argument Reference

The follow arguments are supported:

* `cluster_id` - (Required) The VKE cluster ID you want to attach this nodepool to.
* `node_quantity` - (Required) The number of nodes in this node pool.
* `plan` - (Required) The plan to be used in this node pool. [See Plans List](https://www.vultr.com/api/#operation/list-plans) Note the minimum plan requirements must have at least 1 core and 2 gbs of memory.
* `label` - (Required) The label to be used as a prefix for nodes in this node pool.
* `tag` - (Optional) A tag that is assigned to this node pool.
* `auto_scaler` - (Optional) Enable the auto scaler for the default node pool.
* `min_nodes` - (Optional) The minimum number of nodes to use with the auto scaler.
* `max_nodes` - (Optional) The maximum number of nodes to use with the auto scaler.



## Attributes Reference

The following attributes are exported:
* `id` - The Nodepool ID.
* `cluster_id` - The VKE cluster ID.
* `date_created` - Date of node pool creation.
* `date_updated` - Date of node pool updates.
* `label` - Label of node pool.
* `node_quantity` - Number of nodes within node pool.
* `plan` - Node plan that nodes are using within this node pool.
* `status` - Status of node pool.
* `tag` - Tag for node pool.
* `nodes` - Array that contains information about nodes within this node pool.
* `auto_scaler` - Boolean indicating if the  auto scaler for the default node pool is active.
* `min_nodes` - The minimum number of nodes used by the auto scaler.
* `max_nodes` - The maximum number of nodes used by the auto scaler.

`nodes`

* `date_created` - Date node was created.
* `id` - ID of node.
* `label` - Label of node.
* `status` - Status of node.

## Import
Node pool resources are able to be imported into terraform state like other
resources, however, since they rely on a kubernetes cluster, the import state
requires the UUID of the cluster as well. With that in mind, format the second
argument to the `terraform import` command as a space delimited string of
UUIDs, the first is the cluster ID, the second is the node pool ID. It will
look like this:

```sh
# "clusterID nodePoolID"
terraform import vultr_kubernetes_node_pools.my-k8s-np "7365a98b-5a43-450f-bd27-d768827100e5 ec330340-4f50-4526-858f-a39199f568ac"
```
