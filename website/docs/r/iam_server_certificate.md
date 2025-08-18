---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_server_certificate"
description: |-
  Manages an IAM server certificate.
---

[iam-certificates]: https://docs.k2.cloud/en/services/iam/certificates.html
[lifecycle]: https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle
[sensitive-data]: https://www.terraform.io/docs/state/sensitive-data.html
[rfc3339]: https://tools.ietf.org/html/rfc3339#section-5.8

# Resource: aws_iam_server_certificate

Manages an IAM server certificate. For more information about IAM server certificates,
see [user documentation][iam-certificates].

~> **Note** All arguments including the private key will be stored in the raw state as plain text.
[Read more about sensitive data management][sensitive-data].

## Example Usage

**Using certificates from file:**

```terraform
resource "aws_iam_server_certificate" "example" {
  name             = "example"
  certificate_body = file("self-ca-cert.pem")
  private_key      = file("test-key.pem")
}
```

**Using certificates in-line:**

```terraform
resource "aws_iam_server_certificate" "example" {
  name = "example"

  certificate_body = <<EOF
-----BEGIN CERTIFICATE-----
[......] # cert contents
-----END CERTIFICATE-----
EOF

  private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
[......] # private key contents
-----END RSA PRIVATE KEY-----
EOF
}
```

**Using certificates in combination with an LB resource:**

Some properties of IAM server certificates cannot be updated while they are
in use. In order for Terraform to effectively manage a certificate in this situation, it is
recommended that you utilize the `name_prefix` attribute and enable the
`create_before_destroy` [lifecycle block][lifecycle]. This will allow Terraform
to create a new, updated `aws_iam_server_certificate` resource and replace it in
dependant resources before attempting to destroy the old version.

```terraform
resource "aws_iam_server_certificate" "example" {
  name_prefix      = "example-cert"
  certificate_body = file("self-ca-cert.pem")
  private_key      = file("test-key.pem")

  lifecycle {
    create_before_destroy = true
  }
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
  certificate_arn = aws_iam_server_certificate.example.arn

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

The following arguments are supported:

* `certificate_body` – (Required) The contents of the public key certificate in PEM-encoded format.
* `private_key` – (Required) The contents of the private key in PEM-encoded format.
* `certificate_chain` – (Optional) The contents of the certificate chain.
* `name` - (Optional) The name of the server certificate.
    * _Value length:_ From 1 to 128 symbols
    * _Constraints:_ `name` cannot be specified if `name_prefix` is set
* `name_prefix` - (Optional) Creates a unique name beginning with the specified prefix.
    * _Value length:_ From 1 to 102 symbols
    * _Constraints:_ `name_prefix` cannot be specified if `name` is set

~> **Note** If `name` and `name_prefix` are not specified, Terraform will autogenerate a name with the prefix `terraform-`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the server certificate.
* `expiration` - The expiration date in [RFC3339 format][rfc3339] of the IAM server certificate.
* `id` - The ID of the server certificate.
* `name` - The name of the server certificate.
* `upload_date` - The date in [RFC3339 format][rfc3339] when the IAM server certificate was uploaded.

## Import

IAM server certificates can be imported using the `name`, e.g.,

```
$ terraform import aws_iam_server_certificate.certificate example.com-certificate-until-2018
```
