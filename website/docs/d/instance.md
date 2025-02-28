---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_instance"
description: |-
  Provides information about an EC2 instance.
---

[base64decode-function]: https://www.terraform.io/docs/configuration/functions/base64decode.html
[describe-instances]: https://docs.k2.cloud/en/api/ec2/instances/DescribeInstances.html

# Data Source: aws_instance

Provides information about an EC2 instance. This data source can be used to get the ID of an EC2 instance for use in other resources.

## Example Usage

```terraform
data "aws_instance" "selected" {
  instance_id = "i-12345678"

  filter {
    name   = "image-id"
    values = ["cmi-12345678"]
  }

  instance_tags = {
    type = "test"
  }

  filter {
    name   = "tag:Name"
    values = ["example"]
  }
}
```

## Argument Reference

* `get_user_data` - (Optional) Retrieve Base64 encoded user data contents into the `user_data_base64` attribute.
  A SHA-1 hash of the user data contents will always be present in the `user_data` attribute. Defaults to `false`.
* `instance_id` - (Optional) Specify the exact instance ID with which to populate the data source.
* `instance_tags` - (Optional) A map of tags, each pair of which must exactly match a pair on the desired instance.

* `filter` - (Optional) One or more name/value pairs to use as filters.

For more information about filtering, see the [EC2 API documentation][describe-instances].

~> **Note** At least one of `filter`, `instance_tags`, or `instance_id` must be specified.

~> **Note** If anything other than a single match is returned by the search,
Terraform will fail. Ensure that your search is specific enough to return
a single instance ID only.

## Attributes Reference

### Supported attributes

`id` is set to the ID of the found instance. In addition, the following attributes
are exported:

~> **Note** Some values are not always set and may not be available for
interpolation.

* `affinity` - The affinity setting for an instance on a dedicated host.
* `ami` - The ID of the image used to launch the instance.
* `arn` - The ARN of the instance.
* `associate_public_ip_address` - Whether the instance is associated with a public IP address or not.
* `availability_zone` - The availability zone of the instance.
* `ebs_block_device` - The EBS block device mappings of the instance.
    * `delete_on_termination` - If the EBS volume will be deleted on termination.
    * `device_name` - The physical name of the device.
    * `iops` - `0` If the EBS volume is not a provisioned IOPS image, otherwise the supported IOPS count.
    * `snapshot_id` - The ID of the snapshot.
    * `volume_size` - The size of the volume, in GiB.
    * `volume_type` - The volume type.
* `ephemeral_block_device` - The ephemeral block device mappings of the instance.
    * `device_name` - The physical name of the device.
    * `no_device` - Whether the specified device included in the device mapping was suppressed or not.
    * `virtual_name` - The virtual device name.
* `host_id` - The ID of the dedicated host that the instance will be assigned to.
* `instance_state` - The state of the instance. One of: `pending`, `running`, `shutting-down`, `terminated`, `stopping`, `stopped`.
* `instance_type` - The type of the instance.
* `key_name` - The key name of the instance.
* `monitoring` - Whether detailed monitoring is enabled or disabled for the instance.
* `network_interface_id` - The ID of the network interface that was created with the instance.
* `placement_group` - The placement group of the instance.
* `private_dns` - The private DNS name assigned to the instance. For EC2-VPC, this is only available if you've enabled DNS hostnames for your VPC.
* `private_ip` - The private IP address assigned to the instance.
* `secondary_private_ips` - The secondary private IPv4 addresses assigned to the instance's primary network interface in a VPC.
* `public_dns` - The public DNS name assigned to the instance. For EC2-VPC, this is only available if you've enabled DNS hostnames for your VPC.
* `public_ip` - The public IP address assigned to the instance, if applicable. **NOTE**: If you are using an [`aws_eip`](../r/eip.md) with your instance, you should refer to the EIP's address directly and not use `public_ip`, as this field will change after the EIP is attached.
* `root_block_device` - The root block device mappings of the instance
    * `device_name` - The physical name of the device.
    * `delete_on_termination` - If the root block device will be deleted on termination.
    * `iops` - `0` If the volume is not a provisioned IOPS image, otherwise the supported IOPS count.
    * `volume_size` - The size of the volume, in GiB.
    * `volume_type` - The type of the volume.
* `security_groups` - The associated security groups.
* `source_dest_check` - Whether the network interface performs source/destination checking.
* `subnet_id` - The VPC subnet ID.
* `user_data` - SHA-1 hash of user data supplied to the instance.
* `user_data_base64` - Base64 encoded contents of User Data supplied to the instance. Valid UTF-8 contents can be decoded with the [`base64decode` function][base64decode-function]. This attribute is only exported if `get_user_data` is true.
* `tags` - A map of tags assigned to the instance.
* `tenancy` - The placement type.
* `vpc_security_group_ids` - The associated security groups in a non-default VPC.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`credit_specification`, `ebs_block_device.encrypted`, `ebs_block_device.kms_key_id`, `ebs_block_device.throughput`, `ebs_optimized`, `enclave_options`, `get_password_data`, `iam_instance_profile`, `ipv6_addresses`, `maintenance_options`, `metadata_options`, `outpost_arn`, `password_data`, `placement_partition_number`, `root_block_device.encrypted`, `root_block_device.kms_key_id`, `root_block_device.throughput`.
