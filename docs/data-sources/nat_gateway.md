---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_nat_gateway"
description: |-
  Provides information about a NAT gateway.
---

[describe-nat-gateways]: https://docs.k2.cloud/en/api/ec2/actions/nat_gateways/DescribeNatGateways.html

# Data source: aws_nat_gateway

Provides information about a NAT gateway.

## Example usage

### Specific example: get a NAT gateway by the ID of its VPC

```terraform
variable vpc_id {}

data "aws_nat_gateway" "selected" {
  vpc_id = var.vpc_id
}
```

### Specific example: get a NAT gateway by its tags

```terraform
data "aws_nat_gateway" "selected" {
  tags = {
    Name = "gw NAT"
  }
}
```

## Argument reference

The arguments of this data source act as filters for querying the available NAT gateway.

~> **Note** The given filters must exactly match the resource whose data will be exported as attributes.

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-nat-gateways]
* `id` - (Optional, String) The ID of the NAT gateway.
* `state` - (Optional, String) The current state of the NAT gateway.
    * _Valid values:_ `available`, `deleting`, `failed`, `pending`
* `tags` - (Optional, Map of strings) Key-value pairs.
  Must exactly match pairs on the required resource.
* `vpc_id` - (Optional, String) The ID of the VPC in which the NAT gateway is located.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

If any fields are missing from the configuration, then this data source will populate them with data for the selected NAT gateway.

In addition to all arguments above, the following attributes are exported:

* `auto_provision_zones` - (String) The state of automatic Elastic IP address allocation for each availability zone.
    * _Valid values:_ `disabled`, `enabled`
* `availability_mode` - (String) The availability mode of the NAT gateway.
    * _Valid values:_ `regional`
* `connectivity_type` - (String) The connectivity type of the NAT gateway.
    * _Valid values:_ `public`
* `nat_gateway_addresses` - ([Block](#nat_gateway_addresses)) The set of addresses of the NAT gateway.

### nat_gateway_addresses

Each block has the following structure:

* `allocation_id` - (String) The allocation ID of the Elastic IP address.
* `association_id` - (String) The association ID of the Elastic IP address.
* `availability_zone` - (String) The availability zone of the Elastic IP address.
* `is_primary` - (Boolean) Indicates whether the Elastic IP address is the primary address of the NAT gateway.
* `private_ip` - (String) The private IP address.
* `public_ip` - (String) The public IP address.
* `status` - (String) The status of the Elastic IP address.
