---
layout: "vultr"
page_title: "Vultr: vultr_kubernetes_kubeconfig"
sidebar_current: "docs-vultr-datasource-kubernetes-kubeconfig"
description: |-
Retrieve a kubeconfig from a Vultr Kubernetes Engine (VKE) resource.
---

# vultr_kubernetes_kubeconfig

Retrieve a kubeconfig from a Vultr Kubernetes Engine (VKE) resource.

## Example Usage

With an existing kubernetes resource, query for the kubeconfig:

```hcl
data "vultr_kubernetes_kubeconfig" "my_vke_config" {
  cluster_id = resource.vultr_kubernetes.my_vke.id
}
```

## Argument Reference

* `cluster_id` - (Required) The VKE cluster ID from which to pull the kubeconfig.

## Attributes Reference

The following attributes are exported:
* `cluster_id` - The VKE cluster ID.
* `kube_config` - Base64 encoded Kubeconfig for this VKE cluster.
* `cluster_ca_certificate` - The base64 encoded public certificate for the cluster's certificate authority.
* `client_key` - The base64 encoded private key used by clients to access the cluster.
* `client_certificate` - The base64 encoded public certificate used by clients to access the cluster.
