---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot_ids"
description: |-
  Provides a list of snapshot IDs.
---

[describe-snapshots]: https://docs.k2.cloud/en/api/ec2/actions/snapshots/DescribeSnapshots.html

# Data Source: aws_ebs_snapshot_ids

Provides a list of snapshot IDs matching the specified criteria.

## Example usage

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

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-snapshots]
* `owners` - (Optional, List of strings) The list of the snapshot owners.
    * _Valid values:_ `project@customer`, `self`
* `restorable_by_user_ids` - (Optional, List of strings) The list of the project IDs (`project@customer`), in which volumes can be created from snapshots.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The region.
* `ids` - (List of strings) The list of the snapshot IDs, sorted by creation time in descending order.
