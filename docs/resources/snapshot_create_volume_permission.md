---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_snapshot_create_volume_permission"
description: |-
  Adds a permission to create volumes from a given snapshot.
---

# Resource: aws_snapshot_create_volume_permission

Adds a permission to create volumes from a given snapshot.

## Example usage

```terraform
resource "aws_snapshot_create_volume_permission" "example_perm" {
  snapshot_id = aws_ebs_snapshot.example_snapshot.id
  account_id  = "project@customer"
}

resource "aws_ebs_volume" "example" {
  availability_zone = "ru-msk-vol52"
  size              = 40
}

resource "aws_ebs_snapshot" "example_snapshot" {
  volume_id = aws_ebs_volume.example.id
}
```

## Argument reference

The following arguments are supported:

* `account_id` - (Required, Forces new resource, String) The ID of the project.
    * _Valid values:_ `project@customer`
* `snapshot_id` - (Required, Forces new resource, String) The ID of the snapshot.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - A combination of `snapshot_id` and `account_id` separated by a hyphen (`-`).

## Timeouts

Timeouts usage for creating volume permissions is not currently supported.

## Import

Import of the volume permissions is not currently supported.
