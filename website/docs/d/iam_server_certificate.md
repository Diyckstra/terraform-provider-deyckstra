---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_server_certificate"
description: |-
  Provides information about an IAM server certificate.
---

# Data Source: aws_iam_server_certificate

Provides information about an IAM server certificate.
Use this data source to lookup information about IAM server certificates.

## Example Usage

```terraform
data "aws_iam_server_certificate" "example" {
  name_prefix = "example.org"
  latest      = true
}

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
  certificate_arn = data.aws_iam_server_certificate.example.arn

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
```

## Argument Reference

* `latest` - (Optional) Indicates whether to return the server certificate with the latest expiration date.
    * _Default value:_ `false`
* `name` - (Optional) The name of the server certificate.
    * _Constraints:_ `name` cannot be specified if `name_prefix` is set
* `name_prefix` - (Optional) The prefix of the server certificate name.
    * _Constraints:_ `name_prefix` cannot be specified if `name` is set

## Attribute Reference

* `arn` - The Amazon Resource Name (ARN) of the IAM server certificate.
* `certificate_body` - The public key certificate in PEM-encoded format.
* `certificate_chain` - The public key certificate chain in PEM-encoded format if exists, empty otherwise.
* `expiration_date` - The expiration date of the IAM server certificate.
* `id` - The ID of the server certificate.
* `upload_date` - The date when the server certificate was uploaded.
