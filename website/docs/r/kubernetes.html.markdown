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

## Argument Reference

The follow arguments are supported:

* `region` - (Required) The region your VKE cluster will be deployed in.
* `version` - (Required) The version your VKE cluster you want deployed. [See Available Version](https://www.vultr.com/api/#operation/get-kubernetes-versions)
* `label` - (Required) The label for the cluster.
* `ha_controlplanes` - (Optional, Default to False) Boolean indicating if the cluster should be created with multiple, highly available controlplanes.
* `enable_firewall` - (Optional, Default to False) Boolean indicating if the cluster should be created with a managed firewall.
* `vpc_id` - (Optional) The ID of the VPC to use when creating the cluster. If not provided a new VPC will be created instead.

`node_pools` (Required) Defines the default node pool for a cluster using these fields:

* `node_quantity` - (Required) The number of nodes in this node pool.
* `plan` - (Required) The plan to be used in this node pool. [See plans list](https://www.vultr.com/api/#operation/list-plans) Note the minimum plan requirements must have at least 1 core and 2 gbs of memory.
* `label` - (Required) The label to be used as a prefix for nodes in this node pool.
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

* `id` - The ID of the node pool.
* `label` - Label of node pool.
* `plan` - Node plan that nodes are using within this node pool.
* `tag` - The default tag that is assigned to the default node pool.
* `node_quantity` - Number of nodes within node pool.
* `auto_scaler` - Boolean indicating if the auto scaler for the default node pool is active.
* `min_nodes` - The minimum number of nodes used by the auto scaler.
* `max_nodes` - The maximum number of nodes used by the auto scaler.
* `user_data` - The base64 encoded user data for nodes in the node pool.
* `status` - Status of node pool.
* `tag` - Tag for node pool.
* `date_created` - Date of node pool creation.
* `date_updated` - Date of node pool updates.

`labels` - A list of labels applied to the nodes in the node pool with these fields:

* `id` - The ID of the label.
* `key` - The key definining the label for kubernetes.
* `value` - The value of the label for kubernetes.

`taints` - A list of taints to apply to the nodes in the node pool with these fields: 

* `id` - The ID of the taint.
* `key` - The key definining the taint for kubernetes.
* `value` - The value of the taint for kubernetes.
* `effect` - The effect of the taint for kubernetes. 

`nodes` - Array that contains information about nodes within this node pool.

* `id` - ID of node.
* `date_created` - Date node was created.
* `label` - Label of node.
* `status` - Status of node.

## Import

A kubernetes cluster created outside of terraform can be imported into the
terraform state using the UUID.  One thing to note is that all kubernetes
resources have a default node pool with a tag of `tf-vke-default`. In order to
avoid errors, ensure that there is a node pool with that tag set that the node
pool matches the configuration in the `node_pools` block of the kubernetes
resource.

```sh
terraform import vultr_kubernetes.my-k8s 7365a98b-5a43-450f-bd27-d768827100e5
```
