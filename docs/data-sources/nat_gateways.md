---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_nat_gateways"
description: |-
  Provides a list of NAT gateway IDs.
---

[describe-nat-gateways]: https://docs.k2.cloud/en/api/ec2/actions/nat_gateways/DescribeNatGateways.html

# Data source: aws_nat_gateways

Provides a list of NAT gateway IDs.
This data source can be useful for retrieving a list of NAT gateway IDs to be referenced elsewhere.

## Example usage

### Specific example: get all available NAT gateways in the region

```terraform
data "aws_nat_gateways" "ngws" {
  filter {
    name   = "state"
    values = ["available"]
  }
}

data "aws_nat_gateway" "ngw" {
  count = length(data.aws_nat_gateways.ngws.ids)
  id    = tolist(data.aws_nat_gateways.ngws.ids)[count.index]
}
```

## Argument reference

The arguments of this data source act as filters for querying the available NAT gateways.

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-nat-gateways]
* `tags` - (Optional, Map of strings) Key-value pairs.
  Must exactly match pairs on the required resources.
* `vpc_id` - (Optional, String) The ID of the VPC in which the NAT gateways are located.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The region.
    * _Example:_ `ru-msk`
* `ids` - (List of strings) The list of all the NAT gateway IDs found.
