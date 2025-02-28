---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interfaces"
description: |-
    Provides a list of network interface ids
---

# Data Source: aws_network_interfaces

## Example Usage

The following shows all network interface IDs.

```terraform
data "aws_network_interfaces" "example" {}

output "example" {
  value = data.aws_network_interfaces.example.ids
}
```

The following example retrieves a list of all network interface IDs with a custom tag of `Name` set to a value of `test`.

```terraform
data "aws_network_interfaces" "example1" {
  tags = {
    Name = "test"
  }
}

output "example1" {
  value = data.aws_network_interfaces.example.ids
}
```

The following example retrieves a network interface IDs which associated with specific subnet.

```terraform
data "aws_network_interfaces" "example2" {
  filter {
    name   = "subnet-id"
    values = ["subnet-xxxxxxxx"]
  }
}

output "example2" {
  value = data.aws_network_interfaces.example.ids
}
```

## Argument Reference

* `tags` - (Optional) A map of tags, each pair of which must exactly match
  a pair on the desired network interfaces.
* `filter` - (Optional) Custom filter block as described below.

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:

* `name` - (Required) The name of the field to filter by it.
* `values` - (Required) Set of values that are accepted for the given field.

For more information about filtering, see the [EC2 API documentation][describe-network-interfaces].

## Attributes Reference

* `id` - The region (e.g., `region-1`).
* `ids` - A list of all the network interface IDs found.

[describe-network-interfaces]: https://docs.k2.cloud/en/api/ec2/network_interfaces/DescribeNetworkInterfaces.html
