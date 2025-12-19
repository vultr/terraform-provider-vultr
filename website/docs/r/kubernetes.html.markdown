---
layout: "vultr"
page_title: "Vultr: vultr_kubernetes"
sidebar_current: "docs-vultr-resource-kubernetes"
description: |-
  Provides a Vultr Kubernetes Engine (VKE) resource. This can be used to create, read, modify, and delete VKE clusters on your Vultr account.
---

# vultr_kubernetes

Get information about a Vultr Kubernetes Engine (VKE) Cluster.

~> The node pool deployed with this resource adds its own `tag` which is then used as an identifier for Terraform to see which node pool is part of this resource. This resource only supports a single node pool. To deploy additional worker nodes you must use `vultr_kubernetes_node_pools`.

## Example Usage

Create a new VKE cluster:

```hcl
resource "vultr_kubernetes" "k8" {
	region  = "ewr"
	label   = "vke-test"
	version = "v1.28.2+1"

	node_pools {
		node_quantity = 1
		plan          = "vc2-1c-2gb"
		label         = "vke-nodepool"
		auto_scaler   = true
		min_nodes     = 1
		max_nodes     = 2

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
	}
}
```

A default node pool is required when first creating the resource but it can be removed at a later point so long as there is a separate `vultr_kubernetes_node_pools` resource attached. For example:

```hcl
resource "vultr_kubernetes" "k8" {
	region  = "ewr"
	label   = "vke-test"
	version = "v1.28.2+1"
} 

# This resource must be created and attached to the cluster
# before removing the default node from the vultr_kubernetes resource
resource "vultr_kubernetes_node_pools" "np" {
	cluster_id    = vultr_kubernetes.k8.id
	node_quantity = 1
	plan          = "vc2-1c-2gb"
	label         = "vke-nodepool"
	auto_scaler   = true
	min_nodes     = 1
	max_nodes     = 2
}
```

There is still a requirement that there be one node pool attached to the cluster but this should allow more flexibility about which node pool that is.

## Argument Reference

The follow arguments are supported:

* `region` - (Required) The region your VKE cluster will be deployed in.
* `version` - (Required) The version your VKE cluster you want deployed. [See Available Version](https://www.vultr.com/api/#operation/get-kubernetes-versions)
* `label` - (Optional) The VKE clusters label.
* `ha_controlplanes` - (Optional, Default to False) Boolean indicating if the cluster should be created with multiple, highly available controlplanes.
* `enable_firewall` - (Optional, Default to False) Boolean indicating if the cluster should be created with a managed firewall.
* `vpc_id` - (Optional) The ID of the VPC to use when creating the cluster. If not provided a new VPC will be created instead.

`node_pools` (Optional) **NOTE** There must be 1 node pool when the kubernetes resource is first created (see explanation above). It supports the following fields

* `node_quantity` - (Required) The number of nodes in this node pool.
* `plan` - (Required) The plan to be used in this node pool. [See Plans List](https://www.vultr.com/api/#operation/list-plans) Note the minimum plan requirements must have at least 1 core and 2 gbs of memory.
* `label` - (Required) The label to be used as a prefix for nodes in this node pool.
* `auto_scaler` - (Optional) Enable the auto scaler for the default node pool.
* `min_nodes` - (Optional) The minimum number of nodes to use with the auto scaler.
* `max_nodes` - (Optional) The maximum number of nodes to use with the auto scaler.
* `labels` - (Optional) A list of labels to apply to the nodes in the node pool. Should contain `key` and `value`.
* `taints` - (Optional) A list of taints to apply to the nodes in the node pool. Should contain `key`, `value` and `effect`.  The `effect` should be one of `NoSchedule`, `PreferNoSchedule` or `NoExecute`.

## Attributes Reference

The following attributes are exported:
* `id` - The VKE cluster ID.
* `label` - The VKE clusters label.
* `region` - The region your VKE cluster is deployed in.
* `ha_controlplanes` - Boolean indicating whether or not the cluster has multiple, highly available controlplanes.
* `firewall_group_id` - The ID of the firewall group managed by this cluster.
* `version` - The current kubernetes version your VKE cluster is running on.
* `status` - The overall status of the cluster.
* `service_subnet` - IP range that services will run on this cluster.
* `cluster_subnet` - IP range that your pods will run on in this cluster.
* `endpoint` - Domain for your Kubernetes clusters control plane.
* `ip` - IP address of VKE cluster control plane.
* `date_created` - Date of VKE cluster creation.
* `kube_config` - Base64 encoded Kubeconfig for this VKE cluster.
* `cluster_ca_certificate` - The base64 encoded public certificate for the cluster's certificate authority.
* `client_key` - The base64 encoded private key used by clients to access the cluster.
* `client_certificate` - The base64 encoded public certificate used by clients to access the cluster.
* `node_pools` - Contains the default node pool that was deployed.

`node_pools`

* `date_created` - Date of node pool creation.
* `date_updated` - Date of node pool updates.
* `label` - Label of node pool.
* `node_quantity` - Number of nodes within node pool.
* `plan` - Node plan that nodes are using within this node pool.
* `status` - Status of node pool.
* `tag` - Tag for node pool.
* `nodes` - Array that contains information about nodes within this node pool.
* `auto_scaler` - Boolean indicating if the auto scaler for the default node pool is active.
* `min_nodes` - The minimum number of nodes used by the auto scaler.
* `max_nodes` - The maximum number of nodes used by the auto scaler.
* `labels` - A list of labels applied to the node pool. Contains `key`, `value` and `id`.
* `taints` - A list of taints applied to the node pool. Contains `key`, `value`, `effect` and `id`.


`nodes`

* `date_created` - Date node was created.
* `id` - ID of node.
* `label` - Label of node.
* `status` - Status of node.

## Import

A kubernetes cluster created outside of terraform can be imported into the
terraform state using the UUID.  One thing to note is that all kubernetes
resources have a default node pool with a tag of `tf-vke-default`. In order to
avoid errors, ensure that there is a node pool with that tag set.

```sh
terraform import vultr_kubernetes.my-k8s 7365a98b-5a43-450f-bd27-d768827100e5
```
