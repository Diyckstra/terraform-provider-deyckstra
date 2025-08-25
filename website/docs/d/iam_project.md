---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_project"
description: |-
  Provides information about an IAM project.
---

[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Data Source: aws_iam_project

Provides information about an IAM project.

## Example Usage

```terraform
data "aws_iam_project" "selected" {
  name = "project-name"
}
```

## Argument Reference

* `name` - (Required) The name of the project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the project.
* `create_date` - The time in [RFC3339 format] when the project was created.
* `display_name` - The displayed name of the project.
* `id` - The name of the project.
* `project_id` - The ID of the project.
* `s3_email` - The email used to set S3 ACL for buckets and objects in the project.
* `state` - The state of the project.
