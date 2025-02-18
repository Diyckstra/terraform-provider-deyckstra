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
* `filter` - (Optional) Custom filter block as described below.

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:

* `name` - (Required) The name of the field to filter by it.
* `values` - (Required) Set of values that are accepted for the given field.
  A network ACL will be selected if any one of the given values matches.

For more information about filtering, see the [EC2 API documentation][describe-network-acls].

## Attributes Reference

* `id` - The region.
* `ids` - A list of all the network ACL ids found.

[describe-network-acls]: https://docs.k2.cloud/en/api/ec2/network_acls/DescribeNetworkAcls.html
