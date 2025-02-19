---
subcategory: "S3 (Simple Storage Service)"
layout: "aws"
page_title: "aws_s3_bucket_policy"
description: |-
    Provides information about a policy of an S3 bucket.
---

# Data Source: aws_s3_bucket_policy

Provides information about a policy of an S3 bucket.

## Example Usage

The following example retrieves the policy of the specified S3 bucket.

```terraform
data "aws_s3_bucket_policy" "example" {
  bucket = "tf-example"
}

output "bucket_policy" {
  value = data.aws_s3_bucket_policy.example.policy
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The bucket name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy` - The bucket policy.
