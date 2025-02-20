---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami"
description: |-
  Provides information about an Amazon Machine Image (AMI).
---

[describe-images]: https://docs.k2.cloud/en/api/ec2/images/DescribeImages.html

# Data Source: aws_ami

Provides information about the ID of a registered image for use in other resources.

## Example Usage

```terraform
data "aws_ami" "example" {
  executable_users = ["self"]
  most_recent      = true
  name_regex       = "^example\\d{1}"
  owners           = ["self"]

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}
```

## Argument Reference

* `owners` – (Required) List of image owners to limit search. At least one value must be specified.
    * _Valid values_: Project ID (`project@customer`) or `self`
* `executable_users` – (Optional) Limit search to project with _explicit_ launch permission on the image.
	  * _Valid values_: Project ID (`project@customer`), `all` or `self`
* `filter` – (Optional) One or more name/value pairs to use as filters.
    * _Valid values_: See supported names and values in [EC2 API documentation][describe-images].
* `most_recent` – (Optional) If more than one result is returned, use the most recent image.
* `name_regex` – (Optional) A regex string to apply to the image list returned by the EC2 API.
  It is recommended to combine this with other options to narrow down the list the EC2 API returns.

~> **Note** The search must return a single match, otherwise Terraform will fail.
Ensure that your search is specific enough to return
a single image ID only, or use `most_recent` to choose the most recent one. If
you want to match multiple images, use the [`aws_ami_ids`](ami_ids.md) data source instead.

## Attributes Reference

### Supported attributes


In addition, the following attributes are exported:

* `arn` – The ARN of the image.
* `architecture` – The OS architecture of the image (e.g., `i386` or `x86_64`).
* `block_device_mappings` – Set of objects with block device mappings of the image.
  The structure of this block is [described below](#block-device-mappings).
* `description` – The description of the image that was provided during image
  creation.
* `id` – The ID of the image.
* `image_id` – The ID of the image. Should be the same as the resource `id`.
* `image_owner_alias` –  The alias of the image owner name.
* `image_type` – The type of image.
* `name` – The name of the image that was provided during image creation.
* `owner_id` – The project ID.
* `platform` – The value is `Windows` for Windows images; otherwise blank.
* `public` – `true` if the image has public launch permissions.
* `root_device_name` – The device name of the root device.
* `root_device_type` – The type of root device (e.g., `ebs` or `instance-store`).
* `root_snapshot_id` – The ID of the snapshot associated with the root device, if any
  (only applies to `ebs` root devices).
* `state` – The current state of the image. If the state is `available`, the image
  is successfully registered and can be used to launch an instance.
* `tags` – Map of tags assigned to the image.
* `virtualization_type` – The type of virtualization of the image (e.g., `hvm`).

#### Block Device Mappings

The `block_device_mappings` block has the following structure:

* `device_name` – The physical name of the device.
* `ebs` – Map containing EBS information, if the device is EBS based.
  Unlike most object attributes, these are accessed directly (e.g., `ebs.volume_size` or `ebs["volume_size"]`) rather than accessed through the first element of a list (e.g., `ebs[0].volume_size`).
  The structure of this subblock is [described below](#ebs).
* `virtual_name` – The virtual device name (for instance stores).

##### EBS

The `ebs` subblock of the [`block_device_mappings`](#block_device_mappings) block has the following structure:

* `delete_on_termination` – `true` if the EBS volume will be deleted on termination.
* `iops` – `0` if the EBS volume is not a provisioned IOPS image, otherwise the supported IOPS count.
* `snapshot_id` – The ID of the snapshot.
* `volume_size` – The size of the volume, in GiB.
* `volume_type` – The volume type.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`boot_mode`, `creation_date`, `deprecation_time`, `ena_support`, `hypervisor`, `image_location`, `kernel_id`, `platform_details`, `product_codes`, `ramdisk_id`, `state_reason`, `sriov_net_support`, `usage_operation`.
