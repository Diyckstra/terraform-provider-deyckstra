---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_policy"
description: |-
  Provides information about an IAM policy.
---

[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Data Source: aws_iam_policy

Provides information about an IAM policy.

## Example Usage

### By ARN

```terraform
data "aws_iam_policy" "selected" {
  arn = "arn:c2:iam:::policy/EC2ReadOnlyAccess"
}
```

### By Name

```terraform
data "aws_iam_policy" "selected" {
  name = "EC2ReadOnlyAccess"
}
```

## Argument Reference

* `arn` - The Amazon Resource Name (ARN) of the policy (e.g. `arn:c2:iam::<customer-name>:policy/<policy-name>`).
    _Constraints_: Required if `name` is not specified
* `name` - The name of the policy.
    _Constraints_: Required if `arn` is not specified

~> **Note** Filtering by `name` is performed locally and can affect performance when the list of policies is large.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `create_date` - The time in [RFC3339 format] when the policy was created.
* `description` - The description of the policy.
* `id` - The ARN of the policy.
* `owner` - The owner of the policy.
* `policy` - Policy-defined access rules in JSON format.
* `policy_id` - The ID of the policy.
* `type` - The type of the policy.
* `update_date` - The time in [RFC3339 format] when the policy was last updated.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`path`, `path_prefix`, `tags`.
