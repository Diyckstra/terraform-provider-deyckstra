---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_listener"
description: |-
  Provides information about a load balancer listener.
---

# Data Source: aws_lb_listener

~> **Note** `aws_alb_listener` is known as `aws_lb_listener`. The functionality is identical.

Provides information about a load balancer listener.

This data source can prove useful when a module accepts an LB listener as an input variable and needs to know the LB it is attached to, or other information specific to the listener in question.

## Example Usage

```terraform
# get listener from listener arn

variable "listener_arn" {
  type = string
}

data "aws_lb_listener" "listener" {
  arn = var.listener_arn
}

# get listener from load_balancer_arn and port

data "aws_lb" "selected" {
  name = "default-public"
}

data "aws_lb_listener" "selected443" {
  load_balancer_arn = data.aws_lb.selected.arn
  port              = 443
}
```

## Argument Reference

The following arguments are supported:

* `arn` - (Optional) The Amazon Resource Name (ARN) of the listener.
    * _Constraints_: Required if `load_balancer_arn` and `port` are not set.
* `load_balancer_arn` - (Optional) The Amazon Resource Name (ARN) of the load balancer.
    * _Constraints_: Required if `arn` is not set.
* `port` - (Optional) Port of the listener.
    * _Constraints_: Required if `arn` is not set.

## Attribute Reference

See the [LB Listener Resource](../r/lb_listener.md) for details on the returned attributes - they are identical.
