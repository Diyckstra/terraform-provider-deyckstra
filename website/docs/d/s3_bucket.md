---
subcategory: "S3 (Simple Storage Service)"
layout: "aws"
page_title: "aws_s3_bucket"
description: |-
    Provides information about a specific S3 bucket.
---

# Data Source: aws_s3_bucket

Provides information about a specific S3 bucket.

## Example Usage

```terraform
data "aws_s3_bucket" "selected" {
  bucket = "tf-example"
}

output "bucket_arn" {
  value = data.aws_s3_bucket.selected.arn
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the bucket.
* `arn` - The ARN of the bucket. Will be of format `arn:aws:s3:::bucketname`.
* `region` - The region this bucket resides in.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`bucket_domain_name`, `bucket_regional_domain_name`, `hosted_zone_id`, `website_domain`, `website_endpoint`.
