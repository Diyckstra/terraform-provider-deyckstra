---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_user_group_membership"
description: |-
  Adds an IAM user to IAM groups.
---

# Resource: aws_iam_user_group_membership

Adds an IAM user to IAM groups.

~> **Note** This resource can be used multiple times with the same user for non-overlapping groups.

How to manage users in a specific group, see the [`aws_iam_group_membership` resource](iam_group_membership.md).

## Example Usage

```terraform
resource "aws_iam_user" "user" {
  name = "tf-user"
}

resource "aws_iam_project" "project" {
  name = "tf-project"
}

resource "aws_iam_group" "group1" {
  name = "tf-group1"
  type = "global"
}

resource "aws_iam_group" "group2" {
  name = "tf-group2"
  type = "project"
}

resource "aws_iam_group" "group3" {
  name = "tf-group3"
  type = "project"
}

resource "aws_iam_user_group_membership" "global-groups" {
  user = aws_iam_user.user.name

  # global groups only
  group_arns = [
    aws_iam_group.group1.arn,
  ]
}

resource "aws_iam_user_group_membership" "project-groups" {
  user    = aws_iam_user.user.name
  project = aws_iam_project.project.name

  # project groups only
  group_arns = [
    aws_iam_group.group2.arn,
    aws_iam_group.group3.arn,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `group_arns` - (Required) List of Amazon Resource Names (ARNs) of the groups to which the user is added
  (e.g. `arn:c2:iam::<customer-name>:group/<group-name>`).

~> **Note** All groups in `group_arns` must be of the same type: **global** or **project**.
If groups are of the **project** type, `project` must be specified.

* `project` - (Optional) The name of the project. Specified when the user is added to project groups.
* `user` - (Required) The name of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Terraform-generated unique ID with the `terraform-` prefix.

## Import

IAM user group membership can be imported using `user`, `project` (optionally), and `group_arns` separated by a hash sign (`#`).

~> **Note** If the user isn't a member of some of the specified groups, these groups will be ignored during import.

Examples:

* import user global group membership:

```
$ terraform import aws_iam_user_group_membership.global-groups user-name#arn:c2:iam:::group/IAMAdministrators
```

* import user project group membership for project `project-example`:

```
$ terraform import aws_iam_user_group_membership.project-groups user-name#project-example#arn:c2:iam:::group/InstanceViewers#arn:c2:iam:::group/BackupOperators
```
