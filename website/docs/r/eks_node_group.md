---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_node_group"
description: |-
  Manages an EKS node group.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[eks-node-groups]: https://docs.k2.cloud/en/services/kubernetes/eks_cluster.html#id7
[lifecycle]: https://www.terraform.io/docs/configuration/meta-arguments/lifecycle.html
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_eks_node_group

Manages an EKS node group, which can provision and optionally update an autoscaling group of Kubernetes worker nodes compatible with EKS.
For details about EKS node groups, see the [user documentation][eks-node-groups].

## Example Usage

```terraform
resource "aws_eks_node_group" "example" {
  cluster_name    = aws_eks_cluster.example.name
  instance_types  = ["c5.large"]
  node_group_name = "example"
  subnet_ids      = aws_subnet.example[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  update_config {
    max_unavailable = 2
  }
}
```

### Ignoring Changes to Desired Size

You can utilize the generic Terraform resource [lifecycle configuration block][lifecycle] with `ignore_changes` to create an EKS node group with an initial size of running instances, then ignore any changes to that count caused externally.

```terraform
resource "aws_eks_node_group" "example" {
  # ... other configurations ...

  scaling_config {
    # Example: Create EKS node group with 2 instances to start
    desired_size = 2

    # ... other configurations ...
  }

  # Optional: Allow external changes without Terraform plan difference
  lifecycle {
    ignore_changes = [scaling_config[0].desired_size]
  }
}
```

### Example Subnets for EKS Node Group

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_subnet" "example" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 8, count.index)
  vpc_id            = aws_vpc.example.id

  tags = {
    "kubernetes.io/cluster/${aws_eks_cluster.example.name}" = "shared"
  }
}
```

## Argument Reference

The following arguments are required:

* `cluster_name` – (Required) Name of the EKS cluster. Must be between 1-100 characters in length. Must begin with an alphanumeric character, and must only contain alphanumeric characters, dashes and underscores (`^[0-9A-Za-z][A-Za-z0-9\-_]+$`).
* `instance_types` - (Required) List of instance types associated with the EKS node group.
* `scaling_config` - (Required) Configuration block with scaling settings. Detailed below.
* `subnet_ids` – (Required) IDs of EC2 subnets to associate with the EKS node group.

The following arguments are optional:

* `capacity_type` - (Optional) Type of capacity associated with the EKS node group. Valid values: `ON_DEMAND`. Terraform will only perform drift detection if a configuration value is provided.
* `disk_size` - (Optional) Disk size in GiB for worker nodes. Defaults to `20`. Terraform will only perform drift detection if a configuration value is provided.
* `labels` - (Optional) Key-value map of Kubernetes labels. Only labels that are applied with the EKS API are managed by this argument. Other Kubernetes labels applied to the EKS node group will not be managed.
* `node_group_name` – (Optional) Name of the EKS node group. If omitted, Terraform will assign a random, unique name. Conflicts with `node_group_name_prefix`.
* `node_group_name_prefix` – (Optional) Creates a unique name beginning with the specified prefix. Conflicts with `node_group_name`.
* `remote_access` - (Optional) Configuration block with remote access settings. Detailed below.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.
* `taint` - (Optional) The Kubernetes taints to be applied to the nodes in the node group. Maximum of 50 taints per node group. Detailed below.

### remote_access Configuration Block

* `ec2_ssh_key` - (Optional) EC2 key pair name that provides access for SSH communication with the worker nodes in the EKS node group.

### scaling_config Configuration Block

* `desired_size` - (Required) Desired number of worker nodes.
* `max_size` - (Required) Maximum number of worker nodes.
* `min_size` - (Required) Minimum number of worker nodes.

### taint Configuration Block

* `key` - (Required) The key of the taint. Maximum length of 63.
* `value` - (Optional) The value of the taint. Maximum length of 63.
* `effect` - (Required) The effect of the taint. Valid values: `NO_SCHEDULE`, `NO_EXECUTE`, `PREFER_NO_SCHEDULE`.

### update_config Configuration Block

The following arguments are mutually exclusive.

* `max_unavailable` - (Optional) Desired max number of unavailable worker nodes during node group update.
* `max_unavailable_percentage` - (Optional) Desired max percentage of unavailable worker nodes during node group update.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - EKS node group ID.
* `id` - EKS cluster name and EKS node group name separated by a colon (`:`).
* `launch_template` - Configuration block with launch template settings.
    * `id` - EC2 launch template ID.
    * `name` - Name of the EC2 launch template.
    * `version` - EC2 launch template version number.
* `resources` - List of objects containing information about underlying resources.
    * `autoscaling_groups` - List of objects containing information about autoscaling groups.
        * `name` - Name of the autoscaling group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `status` - Status of the EKS node group. One of `CREATING`, `ACTIVE`, `PENDING`, `UPDATING`, `DELETING`, `CREATE_FAILED`, `DELETE_FAILED`, `DEGRADED`.
* `version` – Kubernetes version.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`ami_type`, `force_update_version`, `node_role_arn`, `release_version`, `remote_access.source_security_group_ids`, `resources.remote_access_security_group_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `60 minutes`) How long to wait for the EKS node group to be created.
* `update` - (Default `60 minutes`) How long to wait for the EKS node group to be updated.
* Note that the `update` timeout is used separately for both configuration and version update operations.
* `delete` - (Default `60 minutes`) How long to wait for the EKS node group to be deleted.

## Import

EKS node groups can be imported using the `cluster_name` and `node_group_name` separated by a colon (`:`), e.g.,

```
$ terraform import aws_eks_node_group.my_node_group my_cluster:my_node_group
```
