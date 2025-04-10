---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_internet_gateway_attachment"
description: |-
  Attaches an internet gateway to a VPC.
---

# Resource: aws_internet_gateway_attachment

Attaches an internet gateway to a VPC.

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_internet_gateway" "example" {
  tags = {
    Name = "tf-igw"
  }
}

resource "aws_internet_gateway_attachment" "example" {
  internet_gateway_id = aws_internet_gateway.example.id
  vpc_id              = aws_vpc.example.id
}
```

## Argument Reference

The following arguments are supported:

* `internet_gateway_id` - (Required) The ID of the internet gateway.
* `vpc_id` - (Required) The ID of the VPC.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPC and internet gateway separated by a colon (`:`).

## Import

Internet gateway attachments can be imported using `id`, e.g.

```
$ terraform import aws_internet_gateway_attachment.example igw-12345678:vpc-12345678
```
