---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_dx_gateway_attachment"
description: |-
  Provides information on the attachment of an EC2 transit gateway to a Direct Connect gateway.
---

[describe-transit-gateway-attachments]: https://docs.k2.cloud/en/api/ec2/transit_gateways/DescribeTransitGatewayAttachments.html

# Data Source: aws_ec2_transit_gateway_dx_gateway_attachment

Provides information on the attachment of an EC2 transit gateway to a Direct Connect gateway.

## Example Usage

### Using EC2 transit gateway and Direct Connect gateway identifiers to get information on the attachment

```terraform
resource "aws_dx_gateway" "example" {
  name            = "tf-dxgw-example"
  amazon_side_asn = "64512"
}

resource "aws_ec2_transit_gateway" "example" {
}

data "aws_ec2_transit_gateway_dx_gateway_attachment" "selected" {
  transit_gateway_id = aws_ec2_transit_gateway.example.id
  dx_gateway_id      = aws_dx_gateway.example.id
}
```

## Argument Reference

The following arguments are supported:

* `transit_gateway_id` - (Optional) The ID of the EC2 transit gateway.
* `dx_gateway_id` - (Optional) The ID of the Direct Connect gateway.
* `filter` - (Optional) One or more configuration blocks with name-value filters. The structure of a block is [described below](#filter-configuration-block).

### filter Configuration Block

The following arguments are supported by a `filter` configuration block:

* `name` - (Required) The name of the filter. Valid values can be found in [EC2 DescribeTransitGatewayAttachments API Reference][describe-transit-gateway-attachments].
* `values` - (Required) Set of valid values for a given filter field. Results will be selected if any given value matches a filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the EC2 transit gateway attachment.
* `tags` - Tags assigned to the EC2 transit gateway attachment.
