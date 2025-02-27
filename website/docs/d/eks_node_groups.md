---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_node_groups"
description: |-
  Provides the list of EKS node groups names associated with an EKS cluster.
---

# Data Source: aws_eks_node_groups

Provides the list of EKS node groups names associated with an EKS cluster.

## Example Usage

```terraform
data "aws_eks_node_groups" "example" {
  cluster_name = "example"
}

data "aws_eks_node_group" "example" {
  for_each = data.aws_eks_node_groups.example.names

  cluster_name    = "example"
  node_group_name = each.value
}
```

## Argument Reference

* `cluster_name` - (Required) The name of the cluster.

## Attribute Reference

* `id` - Cluster name.
* `names` - A set of all node group names in an EKS cluster.
