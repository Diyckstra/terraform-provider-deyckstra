---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami"
description: |-
  Creates and manages a custom Amazon Machine Image (AMI).
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[images]: https://docs.k2.cloud/en/services/storage/images.html
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_ami

Creates and manages images.

If you just want to share an existing image with another project,
it's better to use [`aws_ami_launch_permission`](ami_launch_permission.md) instead.

For more information about images, see [user documentation][images].

## Example Usage

```terraform
# Create an image that will start a machine whose root device is backed by
# an EBS volume populated from a snapshot. It is assumed that such a snapshot
# already exists with the id "snap-12345678".
resource "aws_ami" "example" {
  name                = "tf-ami"
  virtualization_type = "hvm"
  root_device_name    = "disk1"

  ebs_block_device {
    device_name = "disk1"
    snapshot_id = "snap-12345678"
    volume_size = 8
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` – (Required) An unique name for the image.
* `architecture` – (Optional) Machine architecture for created instances.
    * _Valid values_: `x86_64`
* `description` – (Optional) A longer, human-readable description for the image.
* `ebs_block_device` – (Optional) An EBS block device that should be
  attached to created instances. The structure of this block is [described below](#ebs_block_device).
* `ephemeral_block_device` – (Optional) An ephemeral block device that
  should be attached to created instances. The structure of this block is [described below](#ephemeral_block_device).
* `root_device_name` – (Optional) The name of the root device.
    * _Valid values_: `disk<N>`, `cdrom<N>`, `floppy<N>`, `menu`, where `<N>` is a disk number
    *  _Default value_: `disk1`
* `virtualization_type` – (Optional) Keyword to choose what virtualization mode created instances will use.
    * _Valid values_: `hvm`, `hvm-legacy`
    * _Default value_: `hvm`
* `tags` – (Optional) Map of tags to assign to the resource. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

### ebs_block_device

The `ebs_block_device` blocks has the following structure:

* `device_name` – (Required) The device name of one or more block device mapping entries.
    * _Valid values_: `disk<N>`, `cdrom<N>`, `floppy<N>`, where `<N>` is a disk number
* `delete_on_termination` – (Optional) Boolean controlling whether the EBS volumes created to
  support each created instance will be deleted once that instance is terminated.
    * _Default value_: `true`
* `iops` – (Required only when `volume_type` is `io2`) Number of I/O operations per second the
  created volumes will support.
* `snapshot_id` – (Optional) The ID of an EBS snapshot that will be used to initialize the created
  EBS volumes.
    * _Constraints_:  If set, the `volume_size` attribute must be at least as large as the referenced
  snapshot
* `volume_size` – (Required unless `snapshot_id` is set) The size of created volumes in GiB.
    * _Constraints_:  If `snapshot_id` is set and `volume_size` is omitted then the volume will have the same size as the selected snapshot.
* `volume_type` – (Optional) The type of EBS volume to create.
    * _Valid values_: `st2`, `gp2`, `io2`
    * _Default value_: `st2`

### ephemeral_block_device

The `ephemeral_block_device` block has the following structure:

* `device_name` – (Required) The device name of one or more block device mapping entries.
    * _Valid values_: `cdrom<N>`, `floppy<N>`, where `<N>` is a disk number
* `virtual_name` – (Required) A name for the ephemeral device.
    * _Constraints_: Must match with the device name

### Timeouts

The `timeouts` block allows you to specify [timeouts][timeouts] for certain actions:

* `create` – (Default `40 minutes`) Used when creating the image
* `update` – (Default `40 minutes`) Used when updating the image
* `delete` – (Default `90 minutes`) Used when deregistering the image

## Attributes Reference

### Supported attributes

In addition to the [arguments above](#argument-reference), the following attributes are exported:

* `arn` – The ARN of the image.
* `id` – The ID of the created image.
* `root_snapshot_id` – The snapshot ID for the root volume (for EBS-backed images)
* `image_owner_alias` – The owner alias (for example, `self`) or the project ID.
* `image_type` – The type of image.
* `owner_id` – The project ID.
* `platform` – This value is set to windows for Windows images; otherwise, it is blank.
* `public` – Indicates whether the image has public launch permissions.
* `tags_all` – Map of tags to assign to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the ``terraform.tfstate`` file but their values are preset and cannot be specified in configuration files.

The following attributes are not currently supported:

`boot_mode`, `deprecation_time`, `ena_support`, `ebs_block_device.encrypted`, `ebs_block_device.kms_key_id`, `ebs_block_device.outpost_arn`, `ebs_block_device.throughput`, `hypervisor`, `image_location`, `kernel_id`, `platform_details`, `ramdisk_id`, `sriov_net_support`, `usage_operation`.

## Import

`aws_ami` can be imported using the ID of the image, e.g.,

```
$ terraform import aws_ami.example cmi-12345678
```
