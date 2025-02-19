---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route_tables"
description: |-
    Provides information about route tables.
---

# Data Source: aws_route_tables

Provides information about route tables.

## Example Usage

```terraform
variable vpc_id {}

data "aws_route_tables" "rts" {
  vpc_id = var.vpc_id

  filter {
    name   = "tag:kubernetes.io/kops/role"
    values = ["private", "public"]
  }
}
```

## Argument Reference

* `filter` - (Optional) One or more name/value pairs to use as filters.
  A route table will be selected if any one of the given values matches.
	Valid names and values can be found in the [EC2 API documentation][describe-route-tables].
* `vpc_id` - (Optional) The VPC ID that you want to filter from.
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired route tables.

## Attributes Reference

* `id` - The region.
* `ids` - A list of all the route table ids found.

[describe-route-tables]: https://docs.k2.cloud/en/api/ec2/routes/DescribeRouteTables.html
