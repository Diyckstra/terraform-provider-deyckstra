---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_user_policy_attachment"
description: |-
  Attaches an IAM policy to an IAM user.
---

# Resource: aws_iam_user_policy_attachment

Attaches an IAM policy to an IAM user.

## Example Usage

```terraform
resource "aws_iam_user" "example" {
  name = "tf-user"
}

resource "aws_iam_project" "example" {
  name = "tf-project"
}

resource "aws_iam_policy" "global-policy" {
  name = "tf-policy-global"
  type = "global"

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

resource "aws_iam_policy" "project-policy" {
  name = "tf-policy-project"
  type = "project"

  policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "ec2:DescribeVpcs",
            "ec2:DescribeVpcAttribute",
          ],
        },
      ]
    }
  )
}

resource "aws_iam_user_policy_attachment" "attach-global" {
  user       = aws_iam_user.example.name
  policy_arn = aws_iam_policy.global-policy.arn
}

resource "aws_iam_user_policy_attachment" "attach-project" {
  user       = aws_iam_user.example.name
  project    = aws_iam_project.example.name
  policy_arn = aws_iam_policy.project-policy.arn
}
```

## Argument Reference

The following arguments are supported:

* `policy_arn` - (Required) The Amazon Resource Name (ARN) of the attached policy
  (e.g. `arn:c2:iam::<customer-name>:policy/<policy-name>`).
* `project` - (Optional) The name of the project. Specified when a project policy is attached.
* `user` - (Required) The name of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - `user`, `project` (optionally), and `policy_arn` separated by a hash sign (`#`).

## Import

IAM user policy attachment can be imported using `id`, e.g.,

```
$ terraform import aws_iam_user_policy_attachment.attach-project user-name#project-name#arn:c2:iam:::policy/BackupFullAccess
```
