---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_internet_gateway_attachment"
description: |-
  Attaches an internet gateway to a VPC.
---

# Resource: aws_internet_gateway_attachment

Attaches an internet gateway to a VPC.

## Example usage

### Basic example

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

## Argument reference

The following arguments are supported:

* `internet_gateway_id` - (Required, Forces new resource, String) The ID of the internet gateway.
* `vpc_id` - (Required, Forces new resource, String) The ID of the VPC.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The ID of the internet gateway and VPC separated by a colon (`:`).

## Timeouts

Timeouts usage for the internet gateway attachments is not currently supported.

## Import

Internet gateway attachments can be imported using `id`, for example:

```
$ terraform import aws_internet_gateway_attachment.example igw-12345678:vpc-12345678
```
