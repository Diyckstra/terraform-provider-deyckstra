---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_policy"
description: |-
  Manages an IAM policy.
---

[iam-policies-and-groups]: https://docs.k2.cloud/en/services/iam/policies.html
[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Resource: aws_iam_policy

Manages an IAM policy. For details about IAM policies, see the [user documentation][iam-policies-and-groups].

## Example Usage

### Global Policy

```terraform
resource "aws_iam_policy" "example" {
  name        = "tf-policy-global"
  description = "tf-policy-global description"
  type        = "global"

  # Terraform's "jsonencode" function converts
  # the result of a Terraform expression into the correct JSON syntax.
  policy = jsonencode(
    {
      Statement = [
        {
          Action = ["iam:ListUsers"],
        },
      ]
    }
  )
}
```

### Project Policy

```terraform
resource "aws_iam_policy" "example" {
  name        = "tf-policy-project"
  description = "tf-policy-project description"
  type        = "project"

  policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "cloudwatch:DescribeAlarms",
          ]
        },
        {
          Action = [
            "ec2:DescribeVpcs",
            "ec2:DescribeVpcAttribute",
          ]
        },
      ]
    }
  )
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, Editable) The description of the policy.
  The value must be no longer than 1000 characters.
* `name` - (Optional, Conflicts with `name_prefix`) The name of the policy. The value can only contain Latin letters, numbers, underscores (_),
  plus (+) and equal (=) signs, commas (,), periods (.), at symbols (@) and hyphens (-) (`^[\w+=,.@-]*$`).
  The value must be 1 to 128 characters long.
* `name_prefix` - (Optional, Conflicts with `name`) Creates a unique name starting with the specified prefix.
  The value has the same character restrictions as `name`. The value must be 1 to 102 characters long.

~> **Note** If `name` and `name_prefix` are omitted, Terraform will assign a random unique name with the `terraform-` prefix.

* `policy` - (Required, Editable) A string with policy-defined access rules in JSON format.
* `type` - (Required) The type of the policy.
    * _Valid values_: `global`, `project`

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the policy.
* `create_date` - The time in [RFC3339 format] when the policy was created.
* `id` - The Amazon Resource Name (ARN) of the policy.
* `owner` - The owner of the policy.
* `policy_id` - The ID of the policy.
* `update_date` - The time in [RFC3339 format] when the policy was last updated.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `path`.

## Import

IAM policy can be imported using `arn`, e.g.,

* import a policy `policy-example` provided by a customer `test.customer`:

```
$ terraform import aws_iam_policy.example arn:c2:iam::test.customer:policy/policy-example
```
