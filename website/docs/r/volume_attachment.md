---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_volume_attachment"
description: |-
  Manages an EBS volume attachment
---

# Resource: aws_volume_attachment

Manages an EBS volume attachment as a top level resource, to attach and detach volumes from instances.

~> **Note on EBS block devices:** If you use `ebs_block_device` on an `aws_instance`, Terraform will assume management over the full set of non-root EBS block devices for the instance, and treats additional block devices as drift. For this reason, `ebs_block_device` cannot be mixed with external `aws_ebs_volume` + `aws_volume_attachment` resources for a given instance.

## Example Usage

```terraform
variable instance_id {}

resource "aws_ebs_volume" "example" {
  availability_zone = "ru-msk-vol52"
  size              = 1
}

resource "aws_volume_attachment" "example" {
  volume_id   = aws_ebs_volume.example.id
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `device_name` - (Optional) The device name to expose to the instance.

~> **Note** The parameter `device_name` is deprecated. Its value is ignored.
The device name will be generated during attaching and can be changed.

* `instance_id` - (Required) ID of the instance to attach to.
* `volume_id` - (Required) ID of the volume to be attached.
* `skip_destroy` - (Optional) Set this to `true` if you do not wish to detach the volume from the instance
  to which it is attached at destroy time, and instead just remove the attachment from Terraform state.
  This is useful when destroying an instance which has volumes created by some other means attached.
* `stop_instance_before_detaching` - (Optional) Set this to `true` to ensure
  that the target instance is stopped before trying to detach the volume.
  Stops the instance, if it is not already stopped.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `generated_device_name` - The device name generated during attaching. Value can be changed.
* `instance_id` - ID of the instance.
* `volume_id` - ID of the volume.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported: `force_detach`.

## Import

EBS volume attachments can be imported using `DEVICE_NAME:VOLUME_ID:INSTANCE_ID` (the value of `DEVICE_NAME` is ignored), e.g.,

```
$ terraform import aws_volume_attachment.example disk1:vol-12345678:i-12345678
```
