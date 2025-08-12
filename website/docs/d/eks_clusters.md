---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_clusters"
description: |-
  Provides a list of EKS clusters names.
---

# Data Source: aws_eks_clusters

Provides a list of EKS clusters names.

## Example Usage

```terraform
data "aws_eks_clusters" "example" {}

data "aws_eks_cluster" "example" {
  for_each = toset(data.aws_eks_clusters.example.names)
  name     = each.value
}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The region.
* `names` - Set of EKS clusters names.
