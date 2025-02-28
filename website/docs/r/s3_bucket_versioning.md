---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "aws_s3_bucket_versioning"
description: |-
  Provides an S3 bucket versioning resource.
---

[s3-versioning]: https://docs.k2.cloud/en/services/object_storage/operations.html#s3versioningmanual

# Resource: aws_s3_bucket_versioning

Provides a resource for controlling versioning on an S3 bucket.
Deleting this resource will either suspend versioning on the associated S3 bucket or
simply remove the resource from Terraform state if the associated S3 bucket is unversioned.

For more information about S3 versioning, see [user documentation][s3-versioning].

## Example Usage

### With Versioning Enabled

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}

resource "aws_s3_bucket_versioning" "versioning_example" {
  bucket = aws_s3_bucket.example.id
  versioning_configuration {
    status = "Enabled"
  }
}
```

### With Versioning Disabled

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}

resource "aws_s3_bucket_versioning" "versioning_example" {
  bucket = aws_s3_bucket.example.id
  versioning_configuration {
    status = "Disabled"
  }
}
```

### Object Dependency On Versioning

When you create an object whose `version_id` you need and an `aws_s3_bucket_versioning` resource in the same configuration, you are more likely to have success by ensuring the `s3_object` depends either implicitly (see below) or explicitly (i.e., using `depends_on = [aws_s3_bucket_versioning.example]`) on the `aws_s3_bucket_versioning` resource.

~> **Note** For critical and/or production S3 objects, do not create a bucket, enable versioning, and create an object in the bucket within the same configuration.

This example shows the `aws_s3_object.example` depending implicitly on the versioning resource through the reference to `aws_s3_bucket_versioning.example.bucket` to define `bucket`:

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.example.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_object" "example" {
  bucket = aws_s3_bucket_versioning.example.bucket
  key    = "droeloe"
  source = "example.txt"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, Forces new resource) The name of the S3 bucket.
* `versioning_configuration` - (Required) Configuration block for the versioning parameters [detailed below](#versioning_configuration).

### versioning_configuration

~> **Note** While the `versioning_configuration.status` parameter supports `Disabled`, this value is only intended for _creating_ or _importing_ resources that correspond to unversioned S3 buckets.
Updating the value from `Enabled` or `Suspended` to `Disabled` will result in errors as the S3 API does not support returning buckets to an unversioned state.

The `versioning_configuration` configuration block supports the following arguments:

* `status` - (Required) The versioning state of the bucket. Valid values: `Enabled`, `Suspended`, `Disabled`. `Disabled` should only be used when creating or importing resources that correspond to unversioned S3 buckets.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The `bucket`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`expected_bucket_owner`, `mfa`, `versioning_configuration.mfa_delete`.

## Import

S3 bucket versioning can be imported using the `bucket` e.g.,

```
$ terraform import aws_s3_bucket_versioning.example bucket-name
```
