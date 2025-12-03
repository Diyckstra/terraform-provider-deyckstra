---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_listener"
description: |-
  Provides information about a listener.
---

# Data Source: aws_lb_listener

Provides information about a listener.

## Example Usage

### Get Listener by ARN

```terraform
data "aws_lb_listener" "selected" {
  arn = "arn:c2:elasticloadbalancing::project-name@customer-name:listener/app/lb-12345678/li-12345678"
}
```

### Get Listener by Load Balancer ARN and Port

```terraform
data "aws_lb" "selected" {
  name = "lb-name"
}

data "aws_lb_listener" "selected" {
  load_balancer_arn = data.aws_lb.selected.arn
  port              = 443
}
```

## Argument Reference

The following arguments are supported:

* `arn` - (Optional) The Amazon Resource Name (ARN) of the listener.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:listener/<app|net>/lb-12345678/li-12345678`
    * _Constraints:_ `arn` is required if `load_balancer_arn` and `port` are not specified
* `load_balancer_arn` - (Optional) The ARN of the load balancer.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:loadbalancer/<app|net>/lb-12345678`
    * _Constraints:_ `load_balancer_arn` is required if `arn` is not specified
* `port` - (Optional) The port of the listener.
    * _Constraints:_ `port` is required if `arn` is not specified

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `certificate_arn` - The ARN of the IAM server certificate.
* `default_action` - The default action that is applied to incoming requests.
  The structure of this block is [described below](#default_action).
* `id` - The ARN of the listener.
* `load_balancer_arn` - The ARN of the load balancer.
* `port` - The port on which the listener receives requests.
* `protocol` - The protocol for a client connection to the load balancer.
* `tags` - Map of tags assigned to the listener.

#### default_action

The following arguments are required:

* `forward` - The block with information about forwarding requests to target groups.
  The structure of this block is [described below](#forward).
* `order` - The sequential number of the action.
* `target_group_arn` - The ARN of the target group to forward traffic to.
* `type` - The type of the routing action.

##### forward

The `forward` block has the following structure:

* `target_group` - List of target groups to forward traffic to.
  The structure of this block is [described below](#target_group).

###### target_group

The `target_group` block has the following structure:

* `arn` - The ARN of the target group to forward traffic to.
* `weight` - The weight of the target group.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`alpn_policy`, `default_action.authenticate_cognito`, `default_action.authenticate_oidc`, `default_action.fixed_response`, `default_action.redirect`, `forward.stickiness`, `ssl_policy`.
