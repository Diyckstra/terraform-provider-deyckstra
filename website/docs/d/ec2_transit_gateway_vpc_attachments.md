---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_vpc_attachments"
description: |-
  Provides list of transit gateway VPC attachment IDs.
---

[describe-tgw-vpc-attachments]: https://docs.k2.cloud/en/api/ec2/transit_gateways/DescribeTransitGatewayVpcAttachments.html

# Data Source: aws_ec2_transit_gateway_vpc_attachments

Provides list of transit gateway VPC attachment IDs.

## Example Usage

### By Filter

```terraform
data "aws_ec2_transit_gateway_vpc_attachments" "selected" {
  filter {
    name   = "state"
    values = ["available"]
  }
}

data "aws_ec2_transit_gateway_vpc_attachment" "vpc-attachments" {
  count = length(data.aws_ec2_transit_gateway_vpc_attachments.selected.ids)
  id    = data.aws_ec2_transit_gateway_vpc_attachments.selected.ids[count.index]
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more configuration blocks containing name-values filters.
  The structure of this block is [described below](#filter).
* `id` - (Optional) The ID of the transit gateway VPC attachment.

### filter

* `name` - (Required) The name of the field to filter by it.
  Valid values can be found in the [EC2 API documentation][describe-tgw-vpc-attachments].
* `values` - (Required) List of one or more values for the filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - List of transit gateway attachment IDs.
