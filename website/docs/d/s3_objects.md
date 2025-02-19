---
subcategory: "S3 (Simple Storage Service)"
layout: "aws"
page_title: "aws_s3_objects"
description: |-
    Provides information about keys and metadata of S3 objects.
---

# Data Source: aws_s3_objects

~> **Note on `max_keys`:** Retrieving very large numbers of keys can adversely affect Terraform's performance.

Provides information about keys (i.e., file names) and other metadata of objects in an S3 bucket.

## Example Usage

The following example retrieves a list of all object keys in an S3 bucket and creates corresponding Terraform object data sources:

```terraform
data "aws_s3_objects" "example" {
  bucket = "tf-example"
}

data "aws_s3_object" "example" {
  count  = length(data.aws_s3_objects.example.keys)
  key    = element(data.aws_s3_objects.example.keys, count.index)
  bucket = data.aws_s3_objects.example.bucket
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) Lists object keys in this S3 bucket.
* `prefix` - (Optional) Limits results to object keys with this prefix.
* `delimiter` - (Optional) A character used to group keys.
* `encoding_type` - (Optional) Encodes keys using this method. Valid value is `url`.
* `max_keys` - (Optional) Maximum object keys to return. Defaults to `1000`.
* `start_after` - (Optional) Returns key names lexicographically after a specific object key in your bucket. S3 lists object keys in UTF-8 character encoding in lexicographical order.
* `fetch_owner` - (Optional) Boolean specifying whether to populate the owner list. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `keys` - List of strings representing object keys.
* `common_prefixes` - List of any keys between `prefix` and the next occurrence of `delimiter` (i.e., similar to subdirectories of the `prefix` "directory"); the list is only returned when you specify `delimiter`.
* `id` - S3 bucket.
* `owners` - List of strings representing object owner IDs (see `fetch_owner` above).
