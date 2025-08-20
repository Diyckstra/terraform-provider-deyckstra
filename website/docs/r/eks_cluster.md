---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster"
description: |-
  Manages an EKS cluster.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[eks-clusters]: https://docs.k2.cloud/en/services/kubernetes/eks_cluster.html
[ha-clusters]: https://docs.k2.cloud/en/services/kubernetes/overview.html#ha-control-plane
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_eks_cluster

Manages an EKS cluster. For details about EKS clusters, see the [user documentation][eks-clusters].

## Example Usage

### EKS High-Availability Cluster

->  **Note**
By default, Terraform creates [high availability clusters][ha-clusters].

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 1)
  availability_zone = "ru-msk-vol52"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_eks_cluster" "example" {
  name    = "tf-cluster-ha"
  version = "1.30.2"

  vpc_config {
    subnet_ids = [aws_subnet.example.id]
  }
}
```

### EKS Cluster with High-Availability Disabled

~> **Note**
This example uses the same VPC and subnet as in the [EKS high-availability cluster example](#eks-high-availability-cluster).

```terraform
resource "aws_eks_cluster" "example" {
  name    = "tf-cluster-disabled-ha"
  version = "1.30.2"

  legacy_cluster_params {
    master_config {
      high_availability = false
      instance_type     = "c5.large"
      volume_type       = "gp2"
      volume_size       = 64
    }
  }

  vpc_config {
    subnet_ids = [aws_subnet.example.id]
  }
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) The name of the cluster. Must be between 1-100 characters in length. Must begin with an alphanumeric character, and must only contain alphanumeric characters, dashes and underscores (`^[0-9A-Za-z][A-Za-z0-9\-_]+$`).
* `version` - (Required) The Kubernetes server version for the cluster.
* `vpc_config` - (Required) Configuration block for the VPC associated with your cluster. Detailed below. Also contains attributes detailed in the Attributes section.

The following arguments are optional:

* `kubernetes_network_config` - (Optional) Configuration block with kubernetes network configuration for the cluster. Detailed below. If removed, Terraform will only perform drift detection if a configuration value is provided.
* `legacy_cluster_params` - (Optional) The parameters for fine-tuning the Kubernetes cluster.
  The structure of this block is [described below](#legacy_cluster_params).
* `tags` - (Optional) Map of tags to assign to the cluster. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

### vpc_config Arguments

* `security_group_ids` - (Optional) List of security group IDs.
* `subnet_ids` - (Required) List of subnet IDs.

### kubernetes_network_config

The following arguments are supported in the `kubernetes_network_config` configuration block:

* `ip_family` - (Optional) The IP family used to assign Kubernetes pod and service addresses.
    * _Valid values:_ `ipv4`
* `service_ipv4_cidr` - (Optional) The CIDR block to assign Kubernetes service IP addresses from. If you don't specify a block, Kubernetes assigns addresses from 10.96.0.0/12 CIDR block.
The block must meet the following requirements:
    * Within one of the following private IP address blocks: 10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16.
    * Doesn't overlap with any CIDR block assigned to the VPC that you selected for VPC.
    * Between /24 and /12.

### legacy_cluster_params

The `legacy_cluster_params` block has the following structure:

* `master_config` - (Optional) The configuration of the master node of the cluster.
  The structure of this block is [described below](#master_config).

#### master_config

The `master_config` block has the following structure:

* `high_availability` - (Required) Indicates whether to deploy a high-availability cluster.
* `instance_type` - (Required) The instance type of the master node.
* `public_ip` - (Optional) The public IP address at which the master node can be accessed.
* `volume_iops` - (Optional) The number of read/write operations per second for the master node volume.
  The parameter must be set if `volume_type` is `io2`.
* `volume_size` - (Required) The size of the master node volume in GiB.
* `volume_type` - (Required) The type of the master node volume.
    * _Valid values:_ `st2`, `gp2`, `io2`

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - Cluster ID.
* `certificate_authority` - Nested attribute containing `certificate-authority-data` for your cluster.
    * `data` - The base64 encoded certificate data required to communicate with your cluster. Add this to the `certificate-authority-data` section of the `kubeconfig` file for your cluster.
* `created_at` - The Unix epoch time stamp in seconds for when the cluster was created.
* `id` - The name of the cluster.
* `platform_version` - The platform version for the cluster.
* `status` - The status of the EKS cluster. One of `CLAIMED`, `CREATING`, `DELETED`, `DELETING`, `ERROR`, `MODIFYING`, `PENDING`, `PROVISIONING`, `READY`, `REPAIRING`.
* `tags_all` - Map of tags assigned to the cluster, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `vpc_config` -  Nested list containing VPC configuration for the cluster.
    * `cluster_security_group_id` - The cluster security group that was created by the cloud for the cluster.
    * `vpc_id` - The VPC associated with your cluster.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`enabled_cluster_log_types`, `encryption_config`, `endpoint`, `identity`, `role_arn`, `vpc_config.endpoint_private_access`, `vpc_config.endpoint_public_access`, `vpc_config.public_access_cidrs`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `30 minutes`) How long to wait for the EKS cluster to be created.
* `update` - (Default `60 minutes`) How long to wait for the EKS cluster to be updated.
Note that the `update` timeout is used separately for both `version` and `vpc_config` update timeouts.
* `delete` - (Default `15 minutes`) How long to wait for the EKS cluster to be deleted.

## Import

EKS clusters can be imported using the `name`, e.g.,

```
$ terraform import aws_eks_cluster.my_cluster my_cluster
```
