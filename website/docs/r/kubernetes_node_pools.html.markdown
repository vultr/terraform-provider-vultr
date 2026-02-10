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
    plan = "vc2-4c-8gb"
    label = "my-label"
    tag = "my-tag"
    auto_scaler = true
    min_nodes = 1
    max_nodes = 2

    labels {
	    key = "my-label" 
	    value = "a-label-on-all-nodes"
    }

    labels {
	    key = "my-second-label" 
	    value = "another-label-on-all-nodes"
    }

    taints {
	key = "a-taint"
	value = "is-tainted"
	effect = "NoExecute"
    }

    taints {
	key = "another-taint"
	value = "is-tainted"
	effect = "NoSchedule"
    }

    user_data = base64encode("This will be added to node user data")
}

```

## Argument Reference

The follow arguments are supported:

* `cluster_id` - (Required) The VKE cluster ID you want to attach this nodepool to.
* `node_quantity` - (Required) The number of nodes in this node pool.
* `plan` - (Required) The plan to be used in this node pool. [See plans list](https://www.vultr.com/api/#operation/list-plans) Note the minimum plan requirements must have at least 1 core and 2 gbs of memory.
* `label` - (Required) The label to be used as a prefix for nodes in this node pool.
* `tag` - (Optional) A tag that is assigned to this node pool.
* `auto_scaler` - (Optional, Default to False) Enable the auto scaler for the default node pool.
* `min_nodes` - (Optional, Default to 1) The minimum number of nodes to use with the auto scaler.
* `max_nodes` - (Optional, Default to 1) The maximum number of nodes to use with the auto scaler.
* `user_data` - (Optional) A base64 encoded string containing the user data to apply to nodes in the node pool.

`labels` - (Optional) A list of labels to apply to the nodes in the node pool with these fields:

* `key` - (Required) The key definining the label for kubernetes.
* `value` - (Required) The value of the label for kubernetes.

`taints` - (Optional) A list of taints to apply to the nodes in the node pool with these fields: 

* `key` - (Required) The key definining the taint for kubernetes.
* `value` - (Required) The value of the taint for kubernetes.
* `effect` - (Required) The effect of the taint for kubernetes.  Must be one of `NoSchedule`, `PreferNoSchedule` or `NoExecute`.

## Attributes Reference

The following attributes are exported:
* `id` - The Nodepool ID.
* `status` - Status of node pool.
* `cluster_id` - The VKE cluster ID.
* `date_created` - Date of node pool creation.
* `date_updated` - Date of node pool updates.
* `label` - Label of node pool.
* `plan` - Node plan that nodes are using within this node pool.
* `tag` - Tag for node pool.
* `node_quantity` - Number of nodes within node pool.
* `auto_scaler` - Boolean indicating if the  auto scaler for the default node pool is active.
* `min_nodes` - The minimum number of nodes used by the auto scaler.
* `max_nodes` - The maximum number of nodes used by the auto scaler.
* `user_data` - A base64 encoded string containing the user data in use by all nodes in the node pool.

`labels` - A list of labels applied to the nodes in the node pool with these fields:

* `id` - The ID of the label.
* `key` - The key definining the label for kubernetes.
* `value` - The value of the label for kubernetes.

`taints` - A list of taints to apply to the nodes in the node pool with these fields: 

* `id` - The ID of the taint.
* `key` - The key definining the taint for kubernetes.
* `value` - The value of the taint for kubernetes.
* `effect` - The effect of the taint for kubernetes. 

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
