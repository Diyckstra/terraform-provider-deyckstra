---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_target_group_attachment"
description: |-
  Attaches a target to a target group.
---

# Resource: aws_lb_target_group_attachment

Attaches a target to a target group.

## Example Usage

```terraform
variable ami {}
variable instance_type {}

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

resource "aws_instance" "example" {
  ami           = var.ami
  instance_type = var.instance_type
  subnet_id     = aws_subnet.example.id

  tags = {
    Name = "tf-instance"
  }
}

resource "aws_lb_target_group_attachment" "example" {
  target_group_arn = aws_lb_target_group.example.arn
  target_id        = aws_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `target_group_arn` - (Required) The Amazon Resource Name (ARN) of the target group.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:targetgroup/tg-12345678`
* `target_id` - (Required) The ID of the target.
* `port` - (Optional) The port on which the target receives traffic.
    * _Valid values:_ From 1 to 65535

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - A generated unique identifier for the attachment.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `availability_zone`.

## Import

Target group attachments cannot be imported.
