---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "aws_s3_bucket_cors_configuration"
description: |-
  Manages an S3 bucket CORS configuration.
---

[cors]: https://docs.k2.cloud/en/services/object_storage/operations.html#cors

# Resource: aws_s3_bucket_cors_configuration

Manages an S3 bucket CORS configuration. For more information about CORS, go to [Cross-Origin Resource Sharing][cors].

~> **Note** S3 Buckets only support a single CORS configuration. Declaring multiple `aws_s3_bucket_cors_configuration` resources to the same S3 Bucket will cause a perpetual difference in configuration.

## Example Usage

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_cors_configuration" "example" {
  bucket = aws_s3_bucket.example.bucket

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }

  cors_rule {
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, Forces new resource) The name of the bucket.
* `cors_rule` - (Required) Set of origins and methods (cross-origin access that you want to allow) [documented below](#cors_rule). You can configure up to 100 rules.

### cors_rule

The `cors_rule` configuration block supports the following arguments:

* `allowed_headers` - (Optional) Set of Headers that are specified in the `Access-Control-Request-Headers` header.
* `allowed_methods` - (Required) Set of HTTP methods that you allow the origin to execute. Valid values are `GET`, `PUT`, `HEAD`, `POST`, and `DELETE`.
* `allowed_origins` - (Required) Set of origins you want customers to be able to access the bucket from.
* `expose_headers` - (Optional) Set of headers in the response that you want customers to be able to access from their applications (for example, from a JavaScript `XMLHttpRequest` object).
* `id` - (Optional) Unique identifier for the rule. The value cannot be longer than 255 characters.
* `max_age_seconds` - (Optional) The time in seconds that your browser is to cache the preflight response for the specified resource.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The `bucket`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported: `expected_bucket_owner`.

## Import

S3 bucket CORS configuration can be imported using the `bucket` e.g.,

```
$ terraform import aws_s3_bucket_cors_configuration.example bucket-name
```
