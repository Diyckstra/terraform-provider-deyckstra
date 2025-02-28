---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_route_tables"
description: |-
  Provides list of transit gateway route table IDs.
---

[describe-tgw-rtb]: https://docs.k2.cloud/en/api/ec2/transit_gateways/DescribeTransitGatewayRouteTables.html

# Data Source: aws_ec2_transit_gateway_route_tables

Provides list of transit gateway route table IDs.

## Example Usage

```terraform
data "aws_ec2_transit_gateway_route_tables" "selected" {}

output "tgw-rtb-ids" {
  value = data.aws_ec2_transit_gateway_route_tables.selected.ids
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more configuration blocks containing name-values filters.
  The structure of this block is [described below](#filter).

### filter

* `name` - (Required) The name of the field to filter by it.
  Valid values can be found in the [EC2 API documentation][describe-tgw-rtb].
* `values` - (Required) List of one or more values for the filter.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The region (e.g., `region-1`).
* `ids` - List of transit gateway route table IDs.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported: `tags`.
