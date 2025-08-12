---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "aws_s3_bucket_lifecycle_configuration"
description: |-
  Manages an S3 bucket lifecycle configuration.
---

[lifecycle-management]: https://docs.k2.cloud/en/services/object_storage/operations.html#id28
[RFC3339 format]: https://tools.ietf.org/html/rfc3339#section-5.8

# Resource: aws_s3_bucket_lifecycle_configuration

Manages an S3 bucket lifecycle configuration.
For more information about lifecycle management, see [user documentation][lifecycle-management].

An S3 lifecycle configuration consists of one or more lifecycle rules. Each rule consists of the following:

* Rule metadata (`id` and `status`)
* [Filter](#filter) identifying objects to which the rule applies
* One or more expiration actions

~> **Note** S3 buckets only support a single lifecycle configuration. Declaring multiple `aws_s3_bucket_lifecycle_configuration` resources to the same S3 bucket will cause a perpetual difference in configuration.

## Example Usage

### With neither a filter nor prefix specified

The lifecycle rule applies to a subset of objects based on the key name prefix (`""`).

```terraform
# This bucket is used in all examples below
resource "aws_s3_bucket" "bucket" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_lifecycle_configuration" "example" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    id = "rule-1"

    expiration {
      days = 365
    }

    status = "Enabled"
  }
}
```

### Specifying an empty filter

The lifecycle rule applies to all objects in the bucket.

```terraform
resource "aws_s3_bucket_lifecycle_configuration" "example" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    id = "rule-1"

    filter {}

    expiration {
      days = 365
    }

    status = "Enabled"
  }
}
```

### Specifying a filter using key prefixes

The lifecycle rule applies to a subset of objects based on the key name prefix (`logs/`).

```terraform
resource "aws_s3_bucket_lifecycle_configuration" "example" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    id = "rule-1"

    filter {
      prefix = "logs/"
    }

    expiration {
      days = 365
    }

    status = "Enabled"
  }
}
```

If you want to apply a lifecycle action to a subset of objects based on different key name prefixes, specify separate rules.

```terraform
resource "aws_s3_bucket_lifecycle_configuration" "example" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    id = "rule-1"

    filter {
      prefix = "logs/"
    }

    expiration {
      date = "2035-11-10T00:00:00Z"
    }

    status = "Enabled"
  }

  rule {
    id = "rule-2"

    filter {
      prefix = "tmp/"
    }


    expiration {
      days = 2
    }

    status = "Enabled"
  }
}
```

### Specifying a filter based on an object tag

The lifecycle rule specifies a filter based on a tag key and value. The rule then applies only to a subset of objects with the specific tag.

```terraform
resource "aws_s3_bucket_lifecycle_configuration" "example" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    id = "rule-1"

    filter {
      tag {
        key   = "Name"
        value = "Staging"
      }
    }

    expiration {
      days = 365
    }

    status = "Enabled"
  }
}
```

### Specifying a filter based on multiple tags

The lifecycle rule directs S3 to perform lifecycle actions on objects with two tags (with the specific tag keys and values). Notice `tags` is wrapped in the `and` configuration block.

```terraform
resource "aws_s3_bucket_lifecycle_configuration" "example" {
  bucket = aws_s3_bucket.bucket.id

  rule {
    id = "rule-1"

    filter {
      and {
        tags = {
          Key1 = "Value1"
          Key2 = "Value2"
        }
      }
    }

    expiration {
      days = 365
    }

    status = "Enabled"
  }
}
```

### Creating a lifecycle configuration for a bucket with versioning

```terraform
resource "aws_s3_bucket" "versioning_bucket" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_acl" "versioning_bucket_acl" {
  bucket = aws_s3_bucket.versioning_bucket.id
  acl    = "private"
}

resource "aws_s3_bucket_versioning" "versioning" {
  bucket = aws_s3_bucket.versioning_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "versioning-bucket-config" {
  # Must have bucket versioning enabled first
  depends_on = [aws_s3_bucket_versioning.versioning]

  bucket = aws_s3_bucket.versioning_bucket.bucket

  rule {
    id = "config"

    filter {
      prefix = "config/"
    }

    noncurrent_version_expiration {
      noncurrent_days = 90
    }

    status = "Enabled"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the source S3 bucket you want S3 to monitor.
* `rule` - (Required) List of configuration blocks describing the rules managing the replication [documented below](#rule).

### rule

~> **Note** The `filter` argument, while optional, is required if the `rule` configuration block does not contain a `prefix` **and** you intend to override the default behavior of setting the rule to filter objects with the empty string prefix (`""`).
Since `prefix` is deprecated by Amazon S3 , we recommend users either specify `filter` or leave both `filter` and `prefix` unspecified.

~> **Note** A rule cannot be updated from having a filter (via either the `rule.filter` parameter or when neither `rule.filter` and `rule.prefix` are specified) to only having a prefix via the `rule.prefix` parameter.

~> **Note** Terraform cannot distinguish a difference between configurations that use `rule.filter {}` and configurations that neither use `rule.filter` nor `rule.prefix`, so a rule cannot be updated from applying to all objects in the bucket via `rule.filter {}` to applying to a subset of objects based on the key prefix `""` and vice versa.

The `rule` configuration block supports the following arguments:

* `expiration` - (Optional) Configuration block that specifies the expiration for the lifecycle of the object in the form of days [documented below](#expiration).
* `filter` - (Optional) Configuration block used to identify objects that a lifecycle rule applies to [documented below](#filter). If not specified, the `rule` will default to using `prefix`.
* `id` - (Required) Unique identifier for the rule. The value cannot be longer than 255 characters.
* `noncurrent_version_expiration` - (Optional) Configuration block that specifies when noncurrent object versions expire [documented below](#noncurrent_version_expiration).
* `prefix` - (Optional) **DEPRECATED** Use `filter` instead. This has been deprecated by Amazon S3. Prefix identifying one or more objects to which the rule applies.
    * _Default value_: An empty string (`""`) if `filter` is not specified
* `status` - (Required) Whether the rule is currently being applied.
    * _Valid values_:  `Enabled` or `Disabled`

### expiration

The `expiration` configuration block supports the following arguments:

* `date` - (Optional) The date the object is to be moved or deleted. Should be in [RFC3339 format]. The time is always midnight UTC, for example, `2015-11-10T00:00:00.000Z`.
* `days` - (Optional) The lifetime, in days, of the objects that are subject to the rule. The value must be a non-zero positive integer.
* `expired_object_delete_marker` - (Optional, Conflicts with `date` and `days`) Indicates whether S3 will remove a delete marker with no noncurrent versions. If set to `true`, the delete marker will be expired; if set to `false` the policy takes no action.

### filter

~> **Note** The `filter` configuration block must either be specified as the empty configuration block (`filter {}`) or with exactly one of `prefix`, `tag`, `and`, `object_size_greater_than` or `object_size_less_than` specified.

The `filter` configuration block supports the following arguments:

* `and`- (Optional) Configuration block used to apply a logical `AND` to two or more predicates [documented below](#and). The lifecycle rule will apply to any object matching all the predicates configured inside the `and` block.
* `prefix` - (Optional) Prefix identifying one or more objects to which the rule applies.
    * _Default value_: An empty string (`""`) if not specified
* `tag` - (Optional) A configuration block for specifying a tag key and value [documented below](#tag).

### noncurrent_version_expiration

The `noncurrent_version_expiration` configuration block supports the following arguments:

* `noncurrent_days` - (Required) The number of days an object is noncurrent before S3 can perform the associated action. Must be a positive integer.

### and

The `and` configuration block supports the following arguments:

* `tags` - (Required) Map of tags to assign to the resource. All of these tags must exist in the object's tag set in order for the rule to apply.

### tag

The `tag` configuration block supports the following arguments:

* `key` - (Required) Name of the object key.
* `value` - (Required) Value of the tag.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The `bucket`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`expected_bucket_owner`, `rule.abort_incomplete_multipart_upload`, `rule.filter.and.object_size_greater_than`, `rule.filter.and.object_size_less_than`, `rule.filter.and.prefix`, `rule.filter.object_size_greater_than`, `rule.filter.object_size_less_than`, `rule.noncurrent_version_expiration.newer_noncurrent_versions`, `rule.noncurrent_version_transition`, `rule.transition`.

## Import

S3 bucket lifecycle configuration can be imported using the `bucket` e.g.,

```sh
$ terraform import aws_s3_bucket_lifecycle_configuration.example bucket-name
```
