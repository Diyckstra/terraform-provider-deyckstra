---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_launch_template"
description: |-
  Manages an EC2 launch template.
---

[asg-create]: https://docs.k2.cloud/en/services/compute/autoscaling.html#createautoscalinggroup
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[describe-images]: https://docs.k2.cloud/en/api/ec2/images/DescribeImages.html

# Resource: aws_launch_template

Manages an EC2 launch template. The resource can be used to create instances or Auto Scaling groups.

## Example Usage

```terraform
resource "aws_launch_template" "example" {
  name = "tf-lt"

  block_device_mappings {
    device_name = "disk1"

    ebs {
      volume_size = 20
    }
  }

  disable_api_termination = true

  instance_initiated_shutdown_behavior = "terminate"

  image_id      = "cmi-12345678" # add image id, change instance type if needed
  instance_type = "m1.micro"

  monitoring {
    enabled = true
  }

  placement {
    availability_zone = "ru-msk-vol52"
  }

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name = "tf-lt"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `block_device_mappings` - (Optional) Specify volumes to attach to the instance besides the volumes specified by the image.
  See [Block Devices](#block-devices) below for details.
* `default_version` - (Optional, Conflicts with `update_default_version`) Default version of the launch template.
* `description` - (Optional) Description of the launch template version.
* `disable_api_termination` - (Optional) If `true`, disables the possibility to terminate an instance via API.
* `image_id` - (Required) The ID of the image from which to launch the instance.
* `instance_initiated_shutdown_behavior` - (Optional) Shutdown behavior for the instance. Valid values are `stop`, `terminate`.
* `instance_type` - (Optional) The type of the instance.
* `key_name` - (Optional) The key name to use for the instance.
* `monitoring` - (Optional) The monitoring option for the instance. See [Monitoring](#monitoring) below for more details.
* `name` - (Optional, Conflicts with `name_prefix`) The name of the launch template. If you leave this blank, Terraform will auto-generate a unique name.
* `name_prefix` - (Optional, Conflicts with `name`) Creates a unique name beginning with the specified prefix.
* `network_interfaces` - (Optional) Customize network interfaces to be attached at instance boot time.
  See [Network Interfaces](#network-interfaces) below for more details.
* `placement` - (Optional) The placement of the instance. See [Placement](#placement) below for more details.
* `tag_specifications` - (Optional) The tags to apply to the resources during launch. See [Tag Specifications](#tag-specifications) below for more details.
* `tags` - (Optional) Map of tags to assign to the launch template. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.
* `update_default_version` - (Optional, Conflicts with `default_version`) Whether to update default version each update.
* `user_data` - (Optional) The base64-encoded user data to provide when launching the instance. The text length must not exceed 16 KB.
* `vpc_security_group_ids` - (Optional) A list of security group IDs to associate with.

### Block devices

Configure additional volumes of the instance besides specified by the image.

To find out more information for an existing image to override the configuration, such as `device_name`, use the [EC2 API][describe-images].

Each `block_device_mappings` supports the following:

* `device_name` - (Optional) The name of the device to mount.
* `ebs` - (Optional) Configure EBS volume properties.
* `no_device` - (Optional) Suppresses the specified device included in the block device mapping.

The `ebs` block supports the following:

* `delete_on_termination` - (Optional) Whether the volume should be destroyed on instance termination.
* `iops` - (Optional) The amount of provisioned IOPS. This must be set with a volume_type of `io2`.
* `snapshot_id` - (Optional) The snapshot ID to mount.
* `volume_size` - (Optional) The size of the volume in gigabytes.
* `volume_type` - (Optional) Type of volume. Valid values are `st2`, `gp2`, `io2`.

### Monitoring

The `monitoring` block supports the following:

* `enabled` - If `true`, the launched EC2 instance will have detailed monitoring enabled.

### Network Interfaces

Attaches one or more network interfaces to the instance.

For the details about configuring network interfaces when creating an auto scaling group, see the [user documentation][asg-create].

Each `network_interfaces` block supports the following:

* `associate_public_ip_address` - (Optional) Whether a public IP address should be associated with the network interface.
  The address will be assigned to the `eth0` interface if there are free allocated external addresses.
  This operation is available only for instances running in the VPC and for new network interfaces.
* `delete_on_termination` - (Optional) Whether the network interface should be destroyed on instance termination.
* `description` - (Optional) Description of the network interface.
* `device_index` - (Optional) The integer index of the network interface attachment.
* `network_interface_id` - (Optional) The ID of the network interface to attach.
* `private_ip_address` - (Optional) The primary private IPv4 address.
* `security_groups` - (Optional) A list of security group IDs to associate.
* `subnet_id` - (Optional) The VPC subnet ID to associate.

### Placement

The placement group of the instance.

The `placement` block supports the following:

* `affinity` - (Optional) The affinity setting for an instance on a dedicated host. The parameter could be set to `host` only if `tenancy` is `host`.
* `availability_zone` - (Optional) The availability zone for the instance.
* `group_name` - (Optional) The name of the placement group for the instance.
* `host_id` - (Optional) The ID of the dedicated host for the instance.
* `tenancy` - (Optional) The tenancy of the instance (if the instance is running in a VPC). Valid values are `default`, `host`. Defaults to `default`.

~> **Note** If you use the `host` value, you may encounter the `NotEnoughResourcesForInstanceType` error when running an instance. To avoid this, it is recommended to provide either the `subnet_id` argument or the `availability_zone` argument.

### Tag Specifications

The tags to apply to the resources during launch. You can tag instances and volumes.

Each `tag_specifications` block supports the following:

* `resource_type` - (Optional) The type of resource to tag. Valid values are `instance`, `volume`.
* `tags` - (Optional) Map of tags to assign to the resource.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the launch template.
* `id` - The ID of the launch template.
* `latest_version` - The latest version of the launch template.
* `tags_all` -Map of tags to assign to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`block_device_mappings.ebs.encrypted`, `block_device_mappings.ebs.kms_key_id`, `block_device_mappings.ebs.throughput`, `block_device_mappings.virtual_name`, `capacity_reservation_specification`, `cpu_options`, `credit_specification`, `ebs_optimized`, `elastic_gpu_specifications`, `elastic_inference_accelerator`, `enclave_options`, `hibernation_options`, `iam_instance_profile`, `instance_market_options`, `instance_requirements`, `kernel_id`, `license_specification`, `maintenance_options`, `metadata_options`, `network_interfaces.associate_carrier_ip_address`, `network_interfaces.interface_type`, `network_interfaces.ipv4_address_count`, `network_interfaces.ipv4_addresses`, `network_interfaces.ipv4_prefix_count`, `network_interfaces.ipv4_prefixes`, `network_interfaces.ipv6_address_count`, `network_interfaces.ipv6_addresses`, `network_interfaces.ipv6_prefix_count`, `network_interfaces.ipv6_prefixes`, `network_interfaces.network_card_index`, `placement.host_resource_group_arn`, `placement.spread_domain`, `placement.partition_number`, `private_dns_name_options`, `ram_disk_id`, `security_group_names`.

## Import

Launch Templates can be imported using the `id`, e.g.,

```
$ terraform import aws_launch_template.web lt-12345678
```
