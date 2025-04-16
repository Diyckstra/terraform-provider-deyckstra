---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_project"
description: |-
  Manages an IAM project.
---

[iam-users-and-projects]: https://docs.k2.cloud/en/services/iam/iam.html
[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Resource: aws_iam_project

Manages an IAM project. For details about IAM projects, see the [user documentation][iam-users-and-projects].

## Example Usage

```terraform
resource "aws_iam_project" "example" {
  name = "tf-project"
}
```

## Argument Reference

* `display_name` - (Optional, Editable) The displayed name of the project.
  If no value is specified, `name` will be used as the displayed name.
* `name` - (Required) The name of the project. The value must start with a Latin letter and
  can only contain Latin letters, numbers, underscores (_), periods (.) and hyphens (-) (`^[a-zA-Z][a-zA-Z0-9_.-]*$`).
  The value must be 1 to 40 characters long.

## Attribute Reference

* `arn` - The Amazon Resource Name (ARN) of the project.
* `create_date` - The time in [RFC3339 format] when the project was created.
* `id` - The name of the project.
* `project_id` - The ID of the project.
* `s3_email` - The email used to set S3 ACL for buckets and objects in the project.
* `state` - The state of the project.

## Import

IAM project can be imported using `name`, e.g.,

```
$ terraform import aws_iam_project.example project-name
```
