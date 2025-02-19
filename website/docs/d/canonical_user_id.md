---
subcategory: "S3 (Simple Storage Service)"
layout: "aws"
page_title: "aws_canonical_user_id"
description: |-
  Provides information about the canonical user ID (S3 User ID) associated with the used account.
---

# Data Source: aws_canonical_user_id

Provides information about the canonical user ID (S3 User ID) associated with the used account.

~> **Note** To use this data source, you must have the `s3:ListAllMyBuckets` permission.

## Example Usage

```terraform
data "aws_canonical_user_id" "current" {}

output "canonical_user_id" {
  value = data.aws_canonical_user_id.current.id
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The canonical user ID (S3 User ID) associated with the used account.
* `display_name` - The human-friendly name linked to the canonical user ID. The bucket owner's display name.
