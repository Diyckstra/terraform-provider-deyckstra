---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_group"
description: |-
    Provides information about a specific Security Group.
---

# Data Source: aws_security_group

Provides information about a specific Security Group.

This resource can be used when a module accepts a Security Group id as
an input variable and needs to, for example, determine the id of the
VPC that the security group belongs to.

## Example Usage

The following example shows how one might accept a Security Group id as a variable
and use this data source to obtain the data necessary to create a subnet.

```terraform
variable "security_group_id" {}

data "aws_security_group" "selected" {
  id = var.security_group_id
}

resource "aws_subnet" "subnet" {
  vpc_id     = data.aws_security_group.selected.vpc_id
  cidr_block = "10.0.1.0/24"
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
security group in the current region. The given filters must match exactly one
security group whose data will be exported as attributes.


* `filter` - (Optional) Custom filter block as described below.
* `id` - (Optional) ID of the specific security group to retrieve.
* `name` - (Optional) The name that the desired security group must have.
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired security group.
* `vpc_id` - (Optional) The id of the VPC that the desired security group belongs to.

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:

* `name` - (Required) The name of the field to filter by it.
* `values` - (Required) Set of values that are accepted for the given field.
  A Security Group will be selected if any one of the given values matches.

For more information about filtering, see the [EC2 API documentation][describe-security-groups].

## Attributes Reference

All argument attributes except `filter` blocks are also exported as
result attributes. This data source will complete the data by populating
any fields that are not included in the configuration with the data for
the selected Security Group.

The following fields are also exported:

* `description` - The description of the security group.
* `arn` - The computed ARN of the security group.

~> **Note** The default security group for a VPC has the name `default`.

[describe-security-groups]: https://docs.k2.cloud/en/api/ec2/security_groups/DescribeSecurityGroups.html
