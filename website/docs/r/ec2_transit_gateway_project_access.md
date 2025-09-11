---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_project_access"
description: |-
  Grants project access to the transit gateway.
---

# Resource: aws_ec2_transit_gateway_project_access

Grants project access to the transit gateway.

~> **Note** Do not manage transit gateway sharing in more than one place. Use either: `shared_owners` on [`aws_ec2_transit_gateway`](ec2_transit_gateway.md), or [`aws_ec2_transit_gateway_project_access`](ec2_transit_gateway_project_access.md) resource. Using both for the same transit gateway will result in a conflict, as both attempt to enforce the sharing state.

## Example Usage

```terraform
resource "aws_iam_project" "example" {
  name = "tf-project"
}

resource "aws_ec2_transit_gateway" "example" {
  description = "tf example"

  tags = {
    Name = "tf-tgw"
  }
}

resource "aws_ec2_transit_gateway_project_access" "example" {
  transit_gateway_id = aws_ec2_transit_gateway.example.id

  # project@customer
  account_id = "${aws_iam_project.example.name}@${split(":", aws_iam_project.example.arn)[4]}"
}
```

## Argument Reference

* `account_id` - (Required) The ID of the project (`project@customer`) that will be granted access to the transit gateway.
* `transit_gateway_id` - (Required) The ID of the transit gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - `transit_gateway_id` and `account_id` separated by a forward slash (`/`). For example: `tgw-12345ABC/project@customer`.

## Import

In Terraform v1.5.0 or later, the resource can be imported by the ID of the transit gateway and the ID of the project separated by a forward slash (`/`) using the `import` block.

For example:

```terraform
import {
  to = aws_ec2_transit_gateway_project_access.example
  id = "tgw-12345ABC/project@customer"
}
```

In earlier Terraform versions, the resource can be imported by the ID of the transit gateway and the ID of the project separated by a forward slash (`/`) using `terraform import`, e.g.:

```console
% terraform import aws_ec2_transit_gateway_project_access.example 'tgw-12345ABC/project@customer'
```
