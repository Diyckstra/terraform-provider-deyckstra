---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_listener_rule"
description: |-
  Manages a listener rule for an application load balancer.
---

[host-header-config]: https://docs.k2.cloud/en/api/elb/datatypes/HostHeaderConditionConfig.html
[path-pattern-config]: https://docs.k2.cloud/en/api/elb/datatypes/PathPatternConditionConfig.html
[redirect-action-config]: https://docs.k2.cloud/en/api/elb/datatypes/RedirectActionConfig.html
[rule-conditions]: https://docs.k2.cloud/en/services/elb/alb.html#albconditiontypes
[rules]: https://docs.k2.cloud/en/services/elb/alb.html#albrules

# Resource: aws_lb_listener_rule

Manages a listener rule for an application load balancer.
For details about listener rules, see the [user documentation][rules].

## Example Usage

### Forward Action

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.1.1.0/24"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_lb" "example" {
  name               = "tf-lb"
  internal           = true
  load_balancer_type = "application"
  subnets            = [aws_subnet.example.id]

  tags = {
    Name = "tf-lb"
  }
}

resource "aws_lb_target_group" "default-tg" {
  name = "tf-lb-default-tg"

  target_type = "instance"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.example.id

  tags = {
    Name = "tf-lb-default-tg"
  }
}

resource "aws_lb_target_group" "another-tg" {
  name = "tf-lb-another-tg"

  target_type = "instance"
  port        = 81
  protocol    = "HTTP"
  vpc_id      = aws_vpc.example.id

  tags = {
    Name = "tf-lb-another-tg"
  }
}

resource "aws_lb_listener" "example" {
  load_balancer_arn = aws_lb.example.arn

  port     = 1222
  protocol = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.default-tg.arn
  }

  tags = {
    Name = "tf-lb-listener"
  }
}

resource "aws_lb_listener_rule" "forward-action" {
  listener_arn = aws_lb_listener.example.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.another-tg.arn
  }

  condition {
    path_pattern {
      values = ["/static/*"]
    }
  }

  condition {
    host_header {
      values = ["example.com"]
    }
  }
}
```

### Weighted Forward Action

~> **Note** This example uses the listener and target groups defined in the [Forward Action example](#forward-action).

```terraform
resource "aws_lb_listener_rule" "weighted-forward-action" {
  listener_arn = aws_lb_listener.example.arn
  priority     = 99

  action {
    type = "forward"
    forward {
      target_group {
        arn    = aws_lb_target_group.default-tg.arn
        weight = 80
      }

      target_group {
        arn    = aws_lb_target_group.another-tg.arn
        weight = 20
      }
    }
  }

  condition {
    host_header {
      values = ["example.com"]
    }
  }
}
```

### Redirect Action

~> **Note** This example uses the listener defined in the [Forward Action example](#forward-action).

```terraform
resource "aws_lb_listener_rule" "redirect-action" {
  listener_arn = aws_lb_listener.example.arn

  action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }

  condition {
    host_header {
      values = ["example.com"]
    }
  }
}
```

### Fixed-Response Action

~> **Note** This example uses the listener defined in the [Forward Action example](#forward-action).

```terraform
resource "aws_lb_listener_rule" "fixed-response-action" {
  listener_arn = aws_lb_listener.example.arn

  action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "HEALTHY"
      status_code  = "200"
    }
  }

  condition {
    host_header {
      values = ["example.com"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required, Editable) The action that will be applied to incoming requests.
  The structure of this block is [described below](#action).
* `condition` - (Required, Editable) List of conditions that will be applied to incoming requests.
  For details about conditions, see the [user documentation][rule-conditions].
  The structure of this block is [described below](#condition).
* `listener_arn` - (Required) The Amazon Resource Name (ARN) of the listener.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:listener/<app|net>/lb-12345678/li-12345678`
* `priority` - (Optional) Priority of the rule.
  If `priority` is not set, its value will be set to the next highest available priority.
    * _Valid values:_ From 1 to 50000

~> **Note** A listener cannot have multiple rules with the same priority.

### action

The `action` block has the following structure:

* `type` - (Required, Editable) The type of the routing action.
    * _Valid values:_ `fixed-response`, `forward`, `redirect`
* `fixed_response` - (Optional, Editable) The block with information about a fixed response to requests.
  The structure of this block is [described below](#fixed_response).
    * _Constraints:_ `fixed_response` can be specified only if `type` is `fixed-response`
* `forward` - (Optional, Editable) The block with information about forwarding requests to target groups.
  The structure of this block is [described below](#forward).
    * _Constraints_: `forward` can be specified only if `type` is `forward` and `target_group_arn` is not specified
* `order` - (Optional, Editable) The sequential number of the action.
    * _Valid values:_ From 1 to 50000
* `redirect` - (Optional, Editable) The block with information about redirecting requests to another URL.
  The structure of this block is [described below](#redirect).
    * _Constraints:_ `redirect` can be specified only if `type` is `redirect`
* `target_group_arn` - (Optional, Editable) The ARN of the target group to forward traffic to.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:targetgroup/tg-12345678`
    * _Constraints:_ `target_group_arn` can be specified only if `type` is `forward` and the `forward` block is not specified

-> **Note** Use `target_group_arn` if you want to forward requests to a single target group. Otherwise, use the `forward` block.

#### fixed_response

The `fixed_response` block has the following structure:

* `content_type` - (Required, Editable) The content type of the response.
    * _Valid values:_ `application/javascript`, `application/json`, `text/css`, `text/html`, `text/plain`
* `status_code` - (Required, Editable) The HTTP code of the response.
    * _Valid values:_ `2XX`, `4XX`, `5XX`, where `X` is a digit
* `message_body` - (Optional, Editable) The message of the response.
    * _Value length:_ From 0 to 1024 symbols
    * _Constraints:_ if `status_code` is `204`, `message_body` must be empty

#### forward

The `forward` block has the following structure:

* `target_group` - (Required, Editable) List of target groups to forward traffic to.
  The structure of this block is [described below](#target_group).
    * _List size:_ From 2 to 5 elements

##### target_group

The `target_group` block has the following structure:

* `arn` - (Required, Editable) The ARN of the target group to forward traffic to.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:targetgroup/tg-12345678`
* `weight` - (Optional, Editable) The weight of the target group.
    * _Valid values:_ From 0 to 256
    * _Default value:_ 1

#### redirect

-> **Note** You can reuse URI components using the following reserved keywords: `#{protocol}`, `#{host}`, `#{port}`, `#{path}` (the leading "/" is removed), and `#{query}`.

* `status_code` - (Required, Editable) The HTTP redirection code.
    * _Valid values:_ `HTTP_301`, `HTTP_302`
* `host` - (Optional, Editable) The name of the host. This component is not percent-encoded.
    * _Default value:_ `#{host}`
    * _Value length:_ From 3 to 255 symbols
    * _Constraints:_ See all the constraints in the [ELB API documentation][redirect-action-config]
* `path` - (Optional, Editable) The case-sensitive absolute path.
  The value can contain `#{host}`, `#{path}`, and `#{port}`.
  This component is not percent-encoded.
    * _Default value:_ `/#{path}`
    * _Value length:_ From 1 to 128 symbols
    * _Constraints:_ The path should start with the leading `/`
* `port` - (Optional, Editable) The port.
    * _Default value:_ `#{port}`
    * _Constraints:_ See all the constraints in the [ELB API documentation][redirect-action-config]
* `protocol` - (Optional, Editable) The protocol.
    * _Valid values:_ `#{protocol}`, `HTTP`, `HTTPS`
    * _Default value:_ `#{protocol}`
* `query` - (Optional, Editable) The query parameters.
  The `?` character is automatically added to the beginning of the line.
  This component is not percent-encoded.
    * _Default value:_ `#{query}`
    * _Value length:_ From 0 to 128 symbols

### condition

The `condition` block has the following structure:

* `host_header` - (Optional, Editable) The block with information about the host headers that the host name in the URL should match.
  The structure of this block is [described below](#host_header).
    * _Constraints:_ `host_header` can be specified only once per rule
* `path_pattern` - (Optional, Editable) The block with information about the path patterns that the path in the URL should match.
  The structure of this block is [described below](#path_pattern).
    * _Constraints:_ `path_pattern` can be specified only once per rule

~> **Note** Exactly one of `host_header` or `path_pattern` must be set per condition.

#### host_header

The `host_header` block has the following structure:

* `values` - (Required, Editable) List of the domain host names that the host name in the URL will be compared to.
  If multiple names are specified, a logical OR operator is applied for the condition.
    * _Value length:_ From 3 to 255 symbols
    * _Constraints:_ See all the constraints in the [ELB API documentation][host-header-config]

#### path_pattern

The `path_pattern` block has the following structure:

* `values` - (Required, Editable) List of the path patterns that the path in the URL will be compared to.
  If multiple values are specified, a logical OR operator is applied for the condition.
    * _Value length:_ From 1 to 128 symbols
    * _Constraints:_ See all the constraints in the [ELB API documentation][path-pattern-config]

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the rule.
* `id` - The ARN of the rule.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`action.authenticate_cognito`, `action.authenticate_oidc`, `condition.http_header`, `condition.http_request_method`, `condition.query_string`, `condition.source_ip`, `forward.stickiness`.

## Import

The listener rule can be imported using `arn`, e.g.,

```
$ terraform import aws_lb_listener_rule.forward-action arn:c2:elasticloadbalancing::project-name@cusomer-name:listener-rule/app/lb-12345678/li-12345678/rule-12345678
```
