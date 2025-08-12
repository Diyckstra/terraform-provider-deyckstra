---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_node_group"
description: |-
  Provides information about an EKS node group.
---

# Data Source: aws_eks_node_group

Provides information about an EKS node group.

## Example Usage

```terraform
data "aws_eks_node_group" "example" {
  cluster_name    = "example"
  node_group_name = "example"
}
```

## Argument Reference

* `cluster_name` - (Required) The name of the cluster.
* `node_group_name` - (Required) The name of the node group.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - EKS node group ID.
* `disk_size` - Disk size in GiB for worker nodes.
* `id` - EKS cluster name and EKS node group name separated by a colon (`:`).
* `instance_types` - Set of instance types associated with the EKS node group.
* `labels` - Key-value map of Kubernetes labels. Only labels that are applied with the EKS API are managed by this argument. Other Kubernetes labels applied to the EKS node group will not be managed.
* `remote_access` - Configuration block with remote access settings.
    * `ec2_ssh_key` - Name of key pair that provides access for SSH communication with the worker nodes in the EKS node group.
* `resources` - List of objects containing information about underlying resources.
    * `autoscaling_groups` - List of objects containing information about autoscaling groups.
        * `name` - Name of the autoscaling group.
* `scaling_config` - Configuration block with scaling settings.
  The structure of this block is [described below](#scaling_config).
* `status` - Status of the EKS node group.
    * _Valid values_: `CREATING`, `ACTIVE`, `PENDING`, `UPDATING`, `DELETING`, `CREATE_FAILED`, `DELETE_FAILED`, `DEGRADED`
* `subnet_ids` - Identifiers of EC2 subnets to associate with the EKS node group.
* `tags` - Map of tags assigned to the node group.
* `taints` - List of objects containing information about taints applied to the nodes in the EKS node group.
  The structure of this block is [described below](#taints)
* `version` - Kubernetes version.

#### scaling_config

The `scaling_config` block has the following structure:

* `desired_size` - Desired number of worker nodes.
* `max_size` - Maximum number of worker nodes.
* `min_size` - Minimum number of worker nodes.

#### taints

The `taints` block has the following structure:

* `key` - The key of the taint.
* `value` - The value of the taint.
* `effect` - The effect of the taint.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`ami_type`, `node_role_arn`, `release_version`, `remote_access.source_security_group_ids`, `resources.remote_access_security_group_id`.
