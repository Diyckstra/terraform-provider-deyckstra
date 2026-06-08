---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_volumes"
description: |-
  Provides a list of volume IDs.
---

[describe-volumes]: https://docs.k2.cloud/en/api/ec2/actions/volumes/DescribeVolumes.html

# Data Source: aws_ebs_volumes

Provides a list of volume IDs matching the specified criteria.

## Example usage

The following example demonstrates obtaining a map of availability zone to volume ID for volumes with a given tag value.

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

### Filter example

For matching against the `size` filter, use:

```terraform
data "aws_ebs_volumes" "ten_or_twenty_gb_volumes" {
  filter {
    name   = "size"
    values = ["10", "20"]
  }
}
```

## Argument reference

In addition to all arguments above, the following attributes are exported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-volumes]
* `tags` - (Optional, Map of strings) Key-value pairs assigned to the volume.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

* `id` - (String) The region.
* `ids` - (List of strings) The list of all the volume IDs found.
