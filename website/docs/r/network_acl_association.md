---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_acl_association"
description: |-
  Manages a network ACL association.
---

# Resource: aws_network_acl_association

Manages a network ACL association. This allows you to associate your network ACL with any subnet(s).

~> **Note on Network ACLs and Network ACL Associations:** Terraform provides both a standalone network ACL association resource
and an [aws_network_acl](network_acl.md) resource with a `subnet_ids` attribute. Do not use the same subnet ID in both a network ACL
resource and a network ACL association resource. Doing so will cause a conflict of associations and will overwrite the association.

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_subnet" "example" {
  availability_zone = "ru-msk-vol52"
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 1, 0)
}

resource "aws_network_acl" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_network_acl_association" "example" {
  network_acl_id = aws_network_acl.example.id
  subnet_id      = aws_subnet.example.id
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_id` - (Required) ID of the network ACL.
* `subnet_id` - (Required) ID of the associated subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the network ACL association

