---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_group"
description: |-
  Provides information about an IAM group.
---

[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Data Source: aws_iam_group

Provides information about an IAM group.

## Example Usage

### By ARN

```terraform
data "aws_iam_group" "selected" {
  arn = "arn:c2:iam:::group/InstanceViewers"
}
```

### By Name

```terraform
data "aws_iam_group" "selected" {
  name = "InstanceViewers"
}
```

## Argument Reference

* `arn` - (Optional) The Amazon Resource Name (ARN) of the group (e.g. `arn:c2:iam::<customer-name>:group/<group-name>`).
    * _Constraints_: Required if `name` is not specified
* `name` - (Optional) The name of the group.
    * _Constraints_: Required if `arn` is not specified

~> **Note** Filtering by `name` is performed locally and can affect performance when the list of groups is large.

## Attribute Reference

### Supported attributes

* `create_date` - The time in [RFC3339 format] when the group was created.
* `group_id` - The ID of the group.
* `id` - The ARN of the group.
* `owner` - The owner of the group.
* `type` - The type of the group.
* `users` - List of group members. The structure of this block is [described below](#users).

#### users

* `arn` - The ARN of a user.
* `user_id` - The ID of a user.
* `user_name` - The name of a user.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`path`, `users.path`.
