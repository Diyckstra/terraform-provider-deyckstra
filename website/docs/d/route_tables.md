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

* `filter` - (Optional) Custom filter block as described below.
* `vpc_id` - (Optional) The VPC ID that you want to filter from.
* `tags` - (Optional) A map of tags, each pair of which must exactly match
  a pair on the desired route tables.

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:

* `name` - (Required) The name of the field to filter by it.
* `values` - (Required) Set of values that are accepted for the given field.
  A Route Table will be selected if any one of the given values matches.

For more information about filtering, see the [EC2 API documentation][describe-route-tables].

## Attributes Reference

* `id` - The region.
* `ids` - A list of all the route table ids found.

[describe-route-tables]: https://docs.k2.cloud/en/api/ec2/routes/DescribeRouteTables.html
