---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot"
description: |-
  Creates a snapshot of an EBS volume.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_ebs_snapshot

Creates a snapshot of an EBS volume.

## Example Usage

```terraform
resource "aws_ebs_volume" "example" {
  availability_zone = "ru-msk-vol52"
  size              = 40

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_ebs_snapshot" "example_snapshot" {
  volume_id = aws_ebs_volume.example.id

  tags = {
    Name = "HelloWorld_snap"
  }
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` - (Required) The volume ID of which to make a snapshot.
* `description` - (Optional) A description of what the snapshot is.
* `tags` - (Optional) Map of tags to assign to the snapshot. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

### Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10 minutes`) Used for creating the ebs snapshot
- `delete` - (Default `10 minutes`) Used for deleting the ebs snapshot

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the EBS snapshot.
* `id` - The snapshot ID (e.g., snap-12345678).
* `owner_id` - The project ID.
* `owner_alias` - The alias of the EBS snapshot owner.
* `volume_size` - The size of the drive in GiB.
* `tags_all` -Map of tags to assign to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`data_encryption_key_id`, `encrypted`, `kms_key_id`, `outpost_arn`, `permanent_restore`, `storage_tier`, `temporary_restore_days`.

## Import

EBS snapshot can be imported using `id`, e.g.,

```
$ terraform import aws_ebs_snapshot.id snap-12345678
```
