---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_node_groups"
description: |-
  Provides information about the EKS node groups names associated with a named EKS cluster.
---

# Data Source: aws_eks_node_groups

Provides information about the EKS node groups names associated with a named EKS cluster.
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

## Attributes Reference

* `id` - Cluster name.
* `names` - A set of all node group names in an EKS cluster.
