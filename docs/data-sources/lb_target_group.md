---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_target_group"
description: |-
  Provides information about a target group.
---

# Data Source: aws_lb_target_group

Provides information about a target group.

## Example Usage

```terraform
data "aws_lb_target_group" "selected" {
  name = "tg-name"
}
```

## Argument Reference

The following arguments are supported:

* `arn` - (Optional) The Amazon Resource Name (ARN) of the target group.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:targetgroup/tg-12345678`
* `name` - (Optional) The name of the target group.

~> **Note** When both `arn` and `name` are specified, `arn` takes precedence.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `health_check` - The health check configuration.
  The structure of this block is [described below](#health_check).
* `id` - The ARN of the target group.
* `port` - The port on which targets receive requests.
* `protocol` - The protocol that is used for routing traffic to the targets.
* `protocol_version` - The version of HTTP protocol.
* `tags` - Map of tags assigned to the target group.
* `target_type` - The type of the target.
* `vpc_id` - The ID of the VPC.

### health_check

The `health_check` block has the following structure:

* `enabled` - Indicates whether health check is enabled.
* `healthy_threshold` - Number of consecutive successful health checks after which the target status changes to healthy.
* `interval` - The amount of time, in seconds, between health checks on an individual target.
* `matcher` - The HTTP code used to check the target availability.
* `path` - The destination for the health check request.
* `port` - The port used to perform health checks on targets.
* `protocol` - The protocol used to perform health checks on targets.
* `timeout` - The amount of time, in seconds, after which no response indicates a failed health check.
* `unhealthy_threshold` - Number of consecutive failed health checks after which the target status changes to unhealthy.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`arn_suffix`, `connection_termination`, `deregistration_delay`, `lambda_multi_value_headers_enabled`, `load_balancing_algorithm_type`, `preserve_client_ip`, `proxy_protocol_v2`, `slow_start`, `stickiness`.
