---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_route_table_propagation"
description: |-
  Manages a transit gateway route table propagation.
---

# Resource: aws_ec2_transit_gateway_route_table_propagation

Manages a transit gateway route table propagation.

## Example Usage

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

  transit_gateway_default_route_table_association = false

  tags = {
    Name = "tf-tgw-vpc-attachment"
  }
}

resource "aws_ec2_transit_gateway_route_table" "example" {
  transit_gateway_id = aws_ec2_transit_gateway.example.id

  tags = {
    Name = "tf-rtb"
  }
}

resource "aws_ec2_transit_gateway_route_table_propagation" "example" {
  transit_gateway_attachment_id  = aws_ec2_transit_gateway_vpc_attachment.example.id
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.example.id
}
```

## Argument Reference

The following arguments are supported:

* `transit_gateway_attachment_id` - (Required) The ID of the transit gateway attachment.
* `transit_gateway_route_table_id` - (Required) The ID of the transit gateway route table.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the transit gateway route table combined with the ID of transit gateway attachment (e.g. `tgw-rtb-12345678_tgw-attach-87654321`).
* `resource_id` - The ID of the resource.
* `resource_type` - The type of the resource.

## Import

The transit gateway route table propagation can be imported using `id`, e.g.,

```
$ terraform import aws_ec2_transit_gateway_route_table_propagation.example tgw-rtb-12345678_tgw-attach-87654321
```
