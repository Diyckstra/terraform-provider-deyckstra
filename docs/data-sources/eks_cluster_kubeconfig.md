---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster_kubeconfig"
description: |-
  Provides the kubeconfig for an EKS cluster.
---

# Data Source: aws_eks_cluster_kubeconfig

Returns the kubeconfig string for an existing EKS cluster. The kubeconfig can
be used to configure `kubectl` or other Kubernetes tooling.

~> **NOTE:** The `kubeconfig` attribute is marked as sensitive because it
contains credentials for the cluster.

## Example Usage

```terraform
data "aws_eks_cluster_kubeconfig" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) The name of the cluster.

## Attribute Reference

* `id` - The name of the cluster.
* `kubeconfig` - The kubeconfig for the cluster. This attribute is sensitive.
