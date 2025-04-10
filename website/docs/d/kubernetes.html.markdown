---
layout: "vultr"
page_title: "Vultr: vultr_kubernetes"
sidebar_current: "docs-vultr-datasource-kubernetes"
description: |-
Get information about a Vultr Kubernetes Engine (VKE) resource. 
---

# vultr_kubernetes

Get information about a Vultr Kubernetes Engine (VKE) Cluster.

## Example Usage

Create a new VKE cluster:

```hcl
data "vultr_kubernetes" "my_vke" {
  filter {
    name   = "label"
    values = ["my-lb-label"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) Query parameters for finding VKE.

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.


## Attributes Reference

The following attributes are exported:
* `id` - The VKE cluster ID.
* `label` - The VKE clusters label.
* `region` - The region your VKE cluster is deployed in.
* `version` - The current kubernetes version your VKE cluster is running on.
* `ha_controlplanes` - Boolean indicating whether or not the cluster has multiple, highly available controlplanes.
* `firewall_group_id` - The ID of the firewall group managed by this cluster.
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
* `labels` - Kubernetes node labels applied to the node pool.
* `taints` - Kubernetes node taints applied to the node pool.

`nodes`

* `date_created` - Date node was created.
* `id` - ID of node.
* `label` - Label of node.
* `status` - Status of node.
