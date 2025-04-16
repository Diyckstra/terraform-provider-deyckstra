---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_group_membership"
description: |-
  Adds IAM users to an IAM group.
---

# Resource: aws_iam_group_membership

Adds IAM users to an IAM group.

~> **Note** Use no more than one resource per group to avoid inconsistent behavior.

How to manage a specific user in multiple groups,
see the [`aws_iam_user_group_membership` resource](iam_user_group_membership.md).

## Example Usage

```terraform
resource "aws_iam_group" "group1" {
  name = "tf-group-1"
  type = "global"
}

resource "aws_iam_group" "group2" {
  name = "tf-group-2"
  type = "project"
}

resource "aws_iam_user" "user1" {
  name = "tf-user-1"
}

resource "aws_iam_user" "user2" {
  name = "tf-user-2"
}

resource "aws_iam_project" "project" {
  name = "tf-project"
}

resource "aws_iam_group_membership" "global-group" {
  name = "tf-global-group-membership"

  group_arn = aws_iam_group.group1.arn

  users = [
    aws_iam_user.user1.name,
    aws_iam_user.user2.name,
  ]
}

resource "aws_iam_group_membership" "project-group" {
  name = "tf-project-group-membership"

  group_arn = aws_iam_group.group2.arn
  project   = aws_iam_project.project.name

  users = [
    aws_iam_user.user2.name,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to identify the group membership.
* `group_arn` – (Required) The Amazon Resource Name (ARN) of the group to which users are added
  (e.g. `arn:c2:iam::<customer-name>:group/<group-name>`).
* `project` - (Optional) The name of the project. Specified when users are added to a project group.
* `users` - (Required) List of user names.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the group membership.
