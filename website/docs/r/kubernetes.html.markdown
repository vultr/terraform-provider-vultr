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
  region = "ewr"
  label     = "tf-test"
  version = "v1.20.11+1"

  node_pools {
    node_quantity = 1
    plan = "vc2-2c-4gb"
    label = "my label"
  }
} 
```

## Argument Reference

The follow arguments are supported:

* `region` - (Required) The region your VKE cluster will be deployed in. Currently, supported values are `ewr` and `lax`
* `version` - (Required) The version your VKE cluster you want deployed. [See Available Version](https://www.vultr.com/api/#operation/get-kubernetes-versions)
* `label` - (Optional) The VKE clusters label.


`node_pools` (Required) There must be 1 node pool with the kubernetes resource. It supports the following fields

* `node_quantity` - (Required) The number of nodes in this node pool.
* `plan` - (Required) The plan to be used in this node pool. [See Plans List](https://www.vultr.com/api/#operation/list-plans) Note the minimum plan requirements must have at least 1 core and 2 gbs of memory.
* `label` - (Required) The label to be used as a prefix for nodes in this node pool.



## Attributes Reference

The following attributes are exported:
* `id` - The VKE cluster ID.
* `label` - The VKE clusters label.
* `region` - The region your VKE cluster is deployed in.
* `version` - The current kubernetes version your VKE cluster is running on.
* `status` - The overall status of the cluster.
* `service_subnet` - IP range that services will run on this cluster.
* `cluster_subnet` - IP range that your pods will run on in this cluster.
* `endpoint` - Domain for your Kubernetes clusters control plane.
* `ip` - IP address of VKE cluster control plane.
* `date_created` - Date of VKE cluster creation.
* `kube_config` - Base64 encoded Kubeconfig for this VKE cluster.
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

`nodes`

* `date_created` - Date node was created.
* `id` - ID of node.
* `label` - Label of node.
* `status` - Status of node.