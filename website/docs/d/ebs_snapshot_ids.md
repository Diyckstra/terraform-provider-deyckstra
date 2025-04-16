---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot_ids"
description: |-
  Provides a list of EBS snapshot IDs.
---

[describe-snapshots]: https://docs.k2.cloud/en/api/ec2/snapshots/DescribeSnapshots.html

# Data Source: aws_ebs_snapshot_ids

Use this data source to get a list of EBS snapshot IDs matching the specified criteria.

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

* `owners` - (Optional) List of the snapshot owners. Valid items are the project ID (`project@customer`) or `self`.
* `restorable_by_user_ids` - (Optional) List of the project IDs (`project@customer`).
  that can create volumes from the snapshot.
* `filter` - (Optional) One or more name/value pairs to filter.

For more information about filtering, see the [EC2 API documentation][describe-snapshots].

## Attributes Reference

* `id` - The region (e.g., `region-1`).
* `ids` - Set of EBS snapshot IDs, sorted by creation time in descending order.
