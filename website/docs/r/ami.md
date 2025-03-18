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

The AMI (Amazon Machine Image) resource allows the creation and management of images.

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

* `name` - (Required) An unique name for the image.
* `description` - (Optional) A longer, human-readable description for the image.
* `root_device_name` - (Optional) The name of the root device. Valid values are `disk<N>`, `cdrom<N>`, `floppy<N>`, `menu` (N - disk number). Defaults to `disk1`.
* `virtualization_type` - (Optional) Keyword to choose what virtualization mode created instances will use. Valid values are `hvm`, `hvm-legacy`. Defaults to `hvm`.
* `architecture` - (Optional) Machine architecture for created instances. Defaults to `x86_64`.
* `ebs_block_device` - (Optional) Nested block describing an EBS block device that should be
  attached to created instances. The structure of this block is described below.
* `ephemeral_block_device` - (Optional) Nested block describing an ephemeral block device that
  should be attached to created instances. The structure of this block is described below.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.

Nested `ebs_block_device` blocks have the following structure:

* `device_name` - (Required) The device name of one or more block device mapping entries. Valid values are `disk<N>`, `cdrom<N>`, `floppy<N>`, (N - disk number).
* `delete_on_termination` - (Optional) Boolean controlling whether the EBS volumes created to
  support each created instance will be deleted once that instance is terminated. Defaults to `true`.
* `iops` - (Required only when `volume_type` is `io2`) Number of I/O operations per second the
  created volumes will support.
* `snapshot_id` - (Optional) The ID of an EBS snapshot that will be used to initialize the created
  EBS volumes. If set, the `volume_size` attribute must be at least as large as the referenced
  snapshot.
* `volume_size` - (Required unless `snapshot_id` is set) The size of created volumes in GiB.
  If `snapshot_id` is set and `volume_size` is omitted then the volume will have the same size
  as the selected snapshot.
* `volume_type` - (Optional) The type of EBS volume to create. Valid values are `st2`, `gp2`, `io2`. Defaults to `st2`.

Nested `ephemeral_block_device` blocks have the following structure:

* `device_name` - (Required) The device name of one or more block device mapping entries. Valid values are `cdrom<N>`, `floppy<N>` (N - disk number).
* `virtual_name` - (Required) A name for the ephemeral device. Must match with the device name.

### Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `40 minutes`) Used when creating the image
* `update` - (Default `40 minutes`) Used when updating the image
* `delete` - (Default `90 minutes`) Used when deregistering the image

## Attributes Reference

### Supported attributes

In addition to the [arguments above](#argument-reference), the following attributes are exported:

* `arn` - The ARN of the image.
* `id` - The ID of the created image.
* `root_snapshot_id` - The Snapshot ID for the root volume (for EBS-backed images)
* `image_owner_alias` - The owner alias (for example, `self`) or the project ID.
* `image_type` - The type of image.
* `owner_id` - The project ID.
* `platform` - This value is set to windows for Windows images; otherwise, it is blank.
* `public` - Indicates whether the image has public launch permissions.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the ``terraform.tfstate`` file but their values are preset and cannot be specified in configuration files.

The following attributes are not currently supported:

`boot_mode`, `deprecation_time`, `ena_support`, `ebs_block_device.encrypted`, `ebs_block_device.kms_key_id`, `ebs_block_device.outpost_arn`, `ebs_block_device.throughput`, `hypervisor`, `image_location`, `kernel_id`, `platform_details`, `ramdisk_id`, `sriov_net_support`, `usage_operation`.

## Import

`aws_ami` can be imported using the ID of the image, e.g.,

```
$ terraform import aws_ami.example cmi-12345678
```
