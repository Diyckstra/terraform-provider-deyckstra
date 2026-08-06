---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster_auth"
description: |-
  Provides information about an authentication token to communicate with an EKS cluster.
---

# Data Source: aws_eks_cluster_auth

This AWS-compatible data source is not supported by K2 EKS. K2 clusters use a
certificate-based kubeconfig instead of an IAM authentication token. Reading
this data source returns an explicit error.

Use `aws_eks_cluster_kubeconfig`:

```terraform
data "aws_eks_cluster_kubeconfig" "example" {
  name = "example"
}

output "kubeconfig" {
  value     = data.aws_eks_cluster_kubeconfig.example.kubeconfig
  sensitive = true
}
```
