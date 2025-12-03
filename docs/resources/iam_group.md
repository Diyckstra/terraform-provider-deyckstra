---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_group"
description: |-
  Manages an IAM group.
---

[iam-policies-and-groups]: https://docs.k2.cloud/en/services/iam/policies.html
[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Resource: aws_iam_group

Manages an IAM group. For details about IAM groups, see the [user documentation][iam-policies-and-groups].

~> **User management in groups**

Manually managing user/group membership via the cloud console alongside using
the [`aws_iam_group_membership`](iam_group_membership.md) or
[`aws_iam_user_group_membership`](iam_user_group_membership.md) resources may result in configuration drift or conflicts.
For this reason, it's recommended to manage membership either entirely using Terraform or entirely in the cloud console.

## Example Usage

```terraform
resource "aws_iam_group" "example" {
  name = "tf-group"
  type = "project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the group. The value can only contain Latin letters, numbers, underscores (_),
  plus (+) and equal (=) signs, commas (,), periods (.), at symbols (@) and hyphens (-) (`^[\w+=,.@-]*$`).
  The value must be 1 to 128 characters long.
* `type` - (Required) The type of the group.
    * _Valid values:_ `global`, `project`

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the group.
* `create_date` - The time in [RFC3339 format] when the group was created.
* `group_id` - The ID of the group.
* `id` - The ARN of the group.
* `owner` - The owner of the group.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `path`.

## Import

IAM groups can be imported using `arn`, e.g.,

* import a group `group-example` provided by a customer `test.customer`:

```
$ terraform import aws_iam_group.example arn:c2:iam::test.customer:group/group-example
```
