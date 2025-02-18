---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_volumes"
description: |-
    Provides identifying information about EBS volumes matching given criteria.
---

# Data Source: aws_ebs_volumes

Provides identifying information about EBS volumes matching given criteria.
This data source can be used to get a list of volume IDs with (for example) matching tags.

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

### Filter example

If matching against the `size` filter, use:

```terraform
data "aws_ebs_volumes" "ten_or_twenty_gb_volumes" {
  filter {
    name   = "size"
    values = ["10", "20"]
  }
}
```

## Argument Reference

* `filter` - (Optional) One or more name/value pairs to use as filters.
	Valid names and values can be found in the [EC2 API documentation][describe-volumes].
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired volumes.

## Attributes Reference

* `id` - The region.
* `ids` - A set of all the EBS volume IDs found.

[describe-volumes]: https://docs.k2.cloud/en/api/ec2/volumes/DescribeVolumes.html
