---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_groups"
description: |-
  Get information about a set of Security Groups.
---

# Data Source: aws_security_groups

Use this data source to get IDs and VPC membership of Security Groups that are created outside of Terraform.

## Example Usage

```terraform
data "aws_security_groups" "test" {
  tags = {
    Application = "k8s"
    Environment = "dev"
  }
}
```

```terraform
variable vpc_id {}

data "aws_security_groups" "test" {
  filter {
    name   = "group-name"
    values = ["nodes"]
  }

  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}
```

## Argument Reference

* `tags` - (Optional) A map of tags, each pair of which must exactly match for desired security groups.
* `filter` - (Optional) One or more name/value pairs to use as filters.

For more information about filtering, see the [EC2 API documentation][describe-security-groups].

## Attributes Reference

* `arns` - ARNs of the matched security groups.
* `id` - The region.
* `ids` - IDs of the matches security groups.
* `vpc_ids` - The VPC IDs of the matched security groups. The data source's tag or filter *will span VPCs* unless the `vpc-id` filter is also used.

[describe-security-groups]: https://docs.k2.cloud/en/api/ec2/security_groups/DescribeSecurityGroups.html
