---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot"
description: |-
  Provides information about an EBS snapshot.
---

# Data Source: aws_ebs_snapshot

Provides information about an EBS snapshot for use when provisioning EBS volumes.

## Example Usage

```terraform
data "aws_ebs_snapshot" "ebs_snapshot" {
  most_recent = true
  owners      = ["self"]

  filter {
    name   = "volume-size"
    values = ["40"]
  }

  filter {
    name   = "tag:Name"
    values = ["Example"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `most_recent` - (Optional) If more than one result is returned, use the most recent snapshot.
* `owners` - (Optional) List of the snapshot owners.
    * _Valid values_: The project ID (`project@customer`) or `self`
* `snapshot_ids` - (Optional) Returns information on a specific snapshot ID.
* `restorable_by_user_ids` - (Optional) List of the project IDs (`project@customer`).
  that can create volumes from the snapshot.
* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-snapshots]

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the EBS snapshot.
* `id` - The snapshot ID (e.g., snap-12345678).
* `snapshot_id` - The snapshot ID (e.g., snap-12345678).
* `description` - A description for the snapshot
* `owner_id` - The project ID.
* `owner_alias` - The alias of the EBS snapshot owner.
* `volume_id` - The volume ID (e.g., vol-12345678).
* `volume_size` - The size of the drive in GiB.
* `state` - The snapshot state.
* `tags` - Map of tags assigned to the resource.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`data_encryption_key_id`, `encrypted`, `kms_key_id`, `outpost_arn`, `storage_tier`.

[describe-snapshots]: https://docs.k2.cloud/en/api/ec2/snapshots/DescribeSnapshots.html
