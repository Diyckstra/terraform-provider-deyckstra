---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_internet_gateway"
description: |-
  Manages an internet gateway.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block

# Resource: aws_internet_gateway

Manages an internet gateway.

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "15.0.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "tf-igw"
  }
}
```

## Argument Reference

The following arguments are supported:

* `tags` - (Optional) Map of tags to assign to the internet gateway.
  If configured with a provider [`default_tags` configuration block][default-tags] present,
  tags with matching keys will overwrite those defined at the provider level.
* `vpc_id` - (Optional) The ID of the VPC to which the internet gateway will be attached.
  See the [aws_internet_gateway_attachment](internet_gateway_attachment.md) resource for another way to attach an internet gateway to a VPC.

-> **Note** It's recommended to explicitly specify an internet gateway as a dependency
  for the resources that require Internet access. For example:

```terraform
resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "tf-igw"
  }
}

resource "aws_instance" "example" {
  # ... other arguments ...

  depends_on = [aws_internet_gateway.example]
}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the internet gateway.
* `id` - The ID of the internet gateway.
* `owner_id` - The ID of the project that the internet gateway belongs to.
* `tags_all` - Map of tags assigned to the internet gateway,
  including those inherited from the provider [`default_tags` configuration block][default-tags].

## Import

Internet gateways can be imported using `id`, e.g.,

```
$ terraform import aws_internet_gateway.example igw-12345678
```
