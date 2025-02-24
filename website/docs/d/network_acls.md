---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_acls"
description: |-
    Provides a list of network ACL IDs for a VPC.
---

# Data Source: aws_network_acls

Provides a list of network ACL IDs for a VPC.

## Example Usage

The following example shows all network ACL ids in a vpc.

```terraform
variable vpc_id {}

data "aws_network_acls" "example" {
  vpc_id = var.vpc_id
}

output "example" {
  value = data.aws_network_acls.example.ids
}
```

The following example retrieves a list of all network ACL ids in a VPC with a custom
tag of `Tier` set to a value of "Private".

```terraform
variable vpc_id {}

data "aws_network_acls" "example" {
  vpc_id = var.vpc_id

  tags = {
    Tier = "Private"
  }
}
```

The following example retrieves a network ACL id in a VPC which associated
with specific subnet.

```terraform
variable vpc_id {}
variable subnet_id {}

data "aws_network_acls" "example" {
  vpc_id = var.vpc_id

  filter {
    name   = "association.subnet-id"
    values = [var.subnet_id]
  }
}
```

## Argument Reference

* `vpc_id` - (Optional) The VPC ID that you want to filter from.
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired network ACLs.
* `filter` - (Optional) One or more name/value pairs to use as filters.
  A network ACL will be selected if any one of the given values matches.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-network-acls]

## Attributes Reference

* `id` - The region.
* `ids` - A list of all the network ACL ids found.

[describe-network-acls]: https://docs.k2.cloud/en/api/ec2/network_acls/DescribeNetworkAcls.html
