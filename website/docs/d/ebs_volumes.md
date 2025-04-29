---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_volumes"
description: |-
    Provides identifying information for EBS volumes matching given criteria
---

# Data Source: aws_ebs_volumes

`aws_ebs_volumes` provides identifying information for EBS volumes matching given criteria.

This data source can be useful for getting a list of volume IDs with (for example) matching tags.

## Example Usage

The following demonstrates obtaining a map of availability zone to EBS volume ID for volumes with a given tag value.

```terraform
data "aws_ebs_volumes" "example" {
  tags = {
    Name = "Example"
  }
}

data "aws_ebs_volume" "example" {
  for_each = toset(data.aws_ebs_volumes.example.ids)
  filter {
    name   = "volume-id"
    values = [each.value]
  }
}

output "availability_zone_to_volume_id" {
  value = { for s in data.aws_ebs_volume.example : s.id => s.availability_zone }
}
```

## Argument Reference

* `filter` - (Optional) Custom filter block as described below.
* `tags` - (Optional) A map of tags, each pair of which must exactly match
  a pair on the desired volumes.

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:

* `name` - (Required) The name of the field to filter by it.
For example, if matching against the `size` filter, use:

```terraform
data "aws_ebs_volumes" "ten_or_twenty_gb_volumes" {
  filter {
    name   = "size"
    values = ["10", "20"]
  }
}
```

* `values` - (Required) Set of values that are accepted for the given field.
  EBS volume IDs will be selected if any one of the given values match.

For more information about filtering, see the [EC2 API documentation][describe-volumes].

## Attributes Reference

* `id` - The region.
* `ids` - A set of all the EBS volume IDs found.

[describe-volumes]: https://docs.k2.cloud/en/api/ec2/volumes/DescribeVolumes.html
