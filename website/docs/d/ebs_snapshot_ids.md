---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot_ids"
description: |-
  Provides a list of EBS snapshot IDs.
---

[describe-snapshots]: https://docs.k2.cloud/en/api/ec2/snapshots/DescribeSnapshots.html

# Data Source: aws_ebs_snapshot_ids

Provides the list of EBS snapshot IDs matching the specified criteria.

## Example Usage

```terraform
data "aws_ebs_snapshot_ids" "ebs_snapshot_ids" {
  owners = ["self"]

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

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-snapshots]
* `owners` - (Optional) List of the snapshot owners.
    * _Valid values_: Project ID (`project@customer`) or `self`
* `restorable_by_user_ids` - (Optional) List of the project IDs (`project@customer`).
  that can create volumes from the snapshot.

## Attribute Reference

* `id` - The region.
* `ids` - Set of EBS snapshot IDs, sorted by creation time in descending order.
