---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_volume_attachment"
description: |-
  Manages a volume attachment.
---

[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_volume_attachment

Manages a volume attachment as a top level resource to attach and detach volumes from instances.

~> **Note on block devices:** If you use `ebs_block_device` on an `aws_instance`, Terraform will assume management over the full set of non-root block devices for the instance, and treats additional block devices as drift. For this reason, `ebs_block_device` cannot be mixed with external `aws_ebs_volume` + `aws_volume_attachment` resources for a given instance.

## Example usage

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

## Argument reference

The following arguments are required:

* `instance_id` - (Required, Forces new resource, String) The ID of the instance to attach to.
* `volume_id` - (Required, Forces new resource, String) The ID of the volume to be attached.

The following arguments are optional:

* `device_name` - (Optional, Editable, String) The device name to expose to the instance.

    ~> **Note** This argument is deprecated. Its value is ignored. The device name will be generated during attaching and can be changed.

* `skip_destroy` - (Optional, Editable, Boolean) Indicates whether to detach the volume on resource destroy or skip the detachment, but still remove the volume's attachment from the Terraform state.
    * _Default value:_ `false`

    ~> **Note**  This can be useful if the volume was created outside of Terraform and you want to keep it attached to the instance.

* `stop_instance_before_detaching` - (Optional, Editable, Boolean) Indicates whether the instance should be stopped before detaching the volume.
    * _Default value:_ `false`

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `generated_device_name` - (String) The device name generated during attaching. The value can be changed.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `force_detach`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `5 minutes`) Used for attaching volume.
* `delete` - (Default `5 minutes`) Used for detaching volume.

## Import

The volume attachment can be imported using `DEVICE_NAME:VOLUME_ID:INSTANCE_ID` (the value of `DEVICE_NAME` is ignored), for example:

```
$ terraform import aws_volume_attachment.example disk1:vol-12345678:i-12345678
```
