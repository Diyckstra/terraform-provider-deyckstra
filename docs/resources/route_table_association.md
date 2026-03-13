---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route_table_association"
description: |-
  Creates an association between a route table and a subnet.
---

# Resource: aws_route_table_association

Creates an association between a route table and a subnet.

## Example usage

### Basic example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_subnet" "example" {
  availability_zone = "ru-msk-vol52"
  vpc_id            = aws_vpc.example.id

  cidr_block = cidrsubnet(aws_vpc.example.cidr_block, 1, 0)
}

resource "aws_route_table_association" "example" {
  subnet_id      = aws_subnet.example.id
  route_table_id = aws_route_table.example.id
}
```

## Argument reference

The following arguments are supported:

* `subnet_id` - (Required, Forces new resource, String) The ID of the subnet to create an association with.
* `route_table_id` - (Required, Editable, String) The ID of the route table to associate with.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attribute is exported:

* `id` - (String) The ID of the association.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `gateway_id`.

## Timeouts

Timeouts usage for the route table associations is not currently supported.

## Import

Route table associations can be imported using the associated resource ID and route table ID separated by a forward slash (`/`).

```
$ terraform import aws_route_table_association.example subnet-12345678/rtb-12345678
```
