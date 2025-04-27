---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_listener_certificate"
description: |-
  Adds a certificate to a listener for an application load balancer.
---

[certificates]: https://docs.k2.cloud/en/services/iam/certificates.html#certificatesmanual

# Resource: aws_lb_listener_certificate

Adds an IAM server certificate to a listener for an application load balancer.
For details about certificates, see the [user documentation][certificates].

This resource adds additional certificates, it does not replace the default listener certificate.

## Example Usage

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

resource "aws_lb_target_group" "example" {
  name = "tf-lb-tg"

  target_type = "instance"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.example.id

  tags = {
    Name = "tf-lb-tg"
  }
}

resource "aws_lb_listener" "example" {
  load_balancer_arn = aws_lb.example.arn

  port            = 1222
  protocol        = "HTTPS"
  certificate_arn = "arn:c2:iam::customer_name:server-certificate/default"

  default_action {
    type = "forward"

    forward {
      target_group {
        arn = aws_lb_target_group.example.arn
      }
    }
  }

  tags = {
    Name = "tf-lb-listener"
  }
}

resource "aws_lb_listener_certificate" "example" {
  listener_arn    = aws_lb_listener.example.arn
  certificate_arn = "arn:c2:iam::customer_name:server-certificate/example"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_arn` - (Required) The Amazon Resource Name (ARN) of the IAM server certificate.
    * _ARN Format:_ `arn:c2:iam::<customer-name>:certificate/<certificate-name>`
* `listener_arn` - (Required) The ARN of the listener.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:listener/<app|net>/lb-12345678/li-12345678`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - `listener_arn` and `certificate_arn` separated by an underscore (`_`).

## Import

This resource can be imported using `id`, e.g.,

```
$ terraform import aws_lb_listener_certificate.example arn:c2:elasticloadbalancing::project-example@customer-name:listener/app/lb-12345678/li-12345678_arn:c2:iam::customer-name:server-certificate/example
```
