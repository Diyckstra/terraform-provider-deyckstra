---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_target_group"
description: |-
  Provides information about a load balancer target group.
---

# Data Source: aws_lb_target_group

~> **Note** `aws_alb_target_group` is known as `aws_lb_target_group`. The functionality is identical.

Provides information about a load balancer target group.

This data source can be used when a module accepts an LB target group as an
input variable and needs to know its attributes. It can also be used to get the ARN of
an LB target group for use in other resources, given LB target group name.

## Example Usage

```terraform
variable "lb_tg_arn" {
  type    = string
  default = ""
}

variable "lb_tg_name" {
  type    = string
  default = ""
}

data "aws_lb_target_group" "test" {
  arn  = var.lb_tg_arn
  name = var.lb_tg_name
}
```

## Argument Reference

The following arguments are supported:

* `arn` - (Optional) The full ARN of the target group.
* `name` - (Optional) The unique name of the target group.

~> **Note**: When both `arn` and `name` are specified, `arn` takes precedence.

## Attributes Reference

See the [LB target group resource](../r/lb_target_group.md) for details
on the returned attributes - they are identical.
