---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eips"
description: |-
    Provides a list of Elastic IPs in a region.
---

# Data Source: aws_eips

Provides a list of Elastic IPs in a region.

## Example Usage

The following shows all Elastic IPs with the specific tag value.

```terraform
data "aws_eips" "example" {
  tags = {
    Env = "dev"
  }
}

output "allocation_ids" {
  value = data.aws_eips.example.allocation_ids
}

output "public_ips" {
  value = data.aws_eips.example.public_ips
}
```

## Argument Reference

* `filter` - (Optional) Custom filter block as described below.
* `tags` - (Optional) Map of tags, each pair of which must exactly match a pair on the desired Elastic IPs.

More complex filters can be expressed using one or more `filter` sub-blocks, which take the following arguments:

* `name` - (Required) The name of the field to filter by it.
* `values` - (Required) Set of values that are accepted for the given field.
* An Elastic IP will be selected if any one of the given values matches.

For more information about filtering, see the [EC2 API documentation][describe-addresses].

[describe-addresses]: https://docs.k2.cloud/en/api/ec2/addresses/DescribeAddresses.html

## Attributes Reference

* `id` - The region.
* `allocation_ids` - A list of all the allocation IDs.
* `public_ips` - A list of all the Elastic IP addresses.
