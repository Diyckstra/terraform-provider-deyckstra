---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_route"
description: |-
  Manages a transit gateway route.
---

# Resource: aws_ec2_transit_gateway_route

Manages a transit gateway route.

## Example Usage

### Standard Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 1)
  availability_zone = "ru-msk-vol52"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_ec2_transit_gateway" "example" {
  description = "tf example"

  tags = {
    Name = "tf-tgw"
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "example" {
  subnet_ids         = [aws_subnet.example.id]
  transit_gateway_id = aws_ec2_transit_gateway.example.id
  vpc_id             = aws_vpc.example.id

  tags = {
    Name = "tf-tgw-vpc-attachment"
  }
}

resource "aws_ec2_transit_gateway_route" "example" {
  destination_cidr_block         = "0.0.0.0/0"
  transit_gateway_attachment_id  = aws_ec2_transit_gateway_vpc_attachment.example.id
  transit_gateway_route_table_id = aws_ec2_transit_gateway.example.association_default_route_table_id
}
```

### Blackhole Route

```terraform
resource "aws_ec2_transit_gateway" "example" {
  description = "tf example"

  tags = {
    Name = "tf-tgw"
  }
}

resource "aws_ec2_transit_gateway_route" "example" {
  destination_cidr_block         = "0.0.0.0/0"
  blackhole                      = true
  transit_gateway_route_table_id = aws_ec2_transit_gateway.example.association_default_route_table_id
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required) The CIDR address block used for the destination match.
* `transit_gateway_attachment_id` - The ID of the transit gateway attachment.
    _Constraints_: Required if `blackhole` is `false`
* `blackhole` - (Optional) Indicates whether to drop traffic that matches this route.
    * _Default value_: `false`
* `transit_gateway_route_table_id` - (Required) The ID of the transit gateway route table.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the transit gateway route table combined with destination (e.g. `tgw-rtb-12345678_0.0.0.0/0`).

## Import

~> **Note** Only transit gateway routes with `static` type can be imported.

The transit gateway route can be imported using `id`, e.g.,

```
$ terraform import aws_ec2_transit_gateway_route.example tgw-rtb-12345678_0.0.0.0/0
```
