---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_clusters"
description: |-
  Provides a list of EKS clusters names.
---

# Data Source: aws_eks_clusters

Provides the list of EKS clusters names matching the specified criteria.

## Example Usage

```terraform
data "aws_eks_clusters" "example" {}

data "aws_eks_cluster" "example" {
  for_each = toset(data.aws_eks_clusters.example.names)
  name     = each.value
}
```

## Attribute Reference

* `id` - Region.
* `names` - Set of EKS clusters names.
