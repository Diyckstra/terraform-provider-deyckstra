---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "aws_s3_bucket_policy"
description: |-
  Attaches a policy to an S3 bucket resource.
---

[policy-restrictions]: https://docs.k2.cloud/en/api/s3/features.html#bucket-policy

# Resource: aws_s3_bucket_policy

Attaches a policy to an S3 bucket resource.

## Example Usage

### Basic Usage

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_policy" "example" {
  bucket = aws_s3_bucket.example.id
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ReadFromSite",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::my-tf-test-bucket/*",
            "Condition":{
               "StringLike":{"aws:Referer":["http://www.site.com/*","http://site.com/*"]}
            }
        }
    ]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to which to apply the policy.
* `policy` - (Required) The text of the policy. Bucket policies are limited to 20 KB in size.

~> **Note** The S3 API supports Bucket Policy with some limitations.
In particular, you cannot specify a user as Principal, but only the project that owns the bucket.
Accordingly, all project users will be granted the same permissions.
For more information about Bucket Policy restrictions, see [user documentation][policy-restrictions].

## Attributes Reference

No additional attributes are exported.

## Import

S3 bucket policies can be imported using the bucket name, e.g.,

```
$ terraform import aws_s3_bucket_policy.example tf-example
```
