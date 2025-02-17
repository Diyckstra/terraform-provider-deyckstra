---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_snapshot_create_volume_permission"
description: |-
  Adds the permission to create volumes from a given EBS snapshot.
---

# Resource: aws_snapshot_create_volume_permission

Adds the permission to create volumes from a given EBS snapshot.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `snapshot_id` - (required) A snapshot ID.
* `account_id` - (required) The project ID (`project@customer`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A combination of "`snapshot_id`-`account_id`".
