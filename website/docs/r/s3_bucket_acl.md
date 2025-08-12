---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "aws_s3_bucket_acl"
description: |-
  Manages an S3 bucket ACL.
---

[access-rights]: https://docs.k2.cloud/en/services/object_storage/operations.html#s3accessrules
[canned-acl]: https://docs.k2.cloud/en/api/s3/acl.html#cannedacl

# Resource: aws_s3_bucket_acl

Manages an S3 bucket ACL.
For more information about access rights for buckets, see [user documentation][access-rights].

~> **Note** `terraform destroy` does not delete the S3 bucket ACL but does remove the resource from Terraform state.

## Example Usage

### With ACL

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_acl" "example_bucket_acl" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}
```

### With Grants

```terraform
data "aws_canonical_user_id" "current" {}

resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  # Use the predefined provider configuration to connect to object storage
  # https://docs.k2.cloud/en/api/tools/terraform.html#providers-tf
  provider = aws.noregion
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  access_control_policy {
    grant {
      grantee {
        id   = data.aws_canonical_user_id.current.id
        type = "CanonicalUser"
      }
      permission = "READ"
    }

    grant {
      grantee {
        type = "Group"
        uri  = "http://acs.amazonaws.com/groups/global/AllUsers"
      }
      permission = "READ_ACP"
    }

    owner {
      id = data.aws_canonical_user_id.current.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional, Conflicts with `access_control_policy`) The [canned ACL][canned-acl] to apply to the bucket.
    * _Valid values_: `private`, `public-read`, `public-read-write`, `authenticated-read`
* `access_control_policy` - (Optional, Conflicts with `acl`) A configuration block that sets the ACL permissions for an object per grantee [documented below](#access_control_policy).
* `bucket` - (Required, Forces new resource) The name of the bucket.

### access_control_policy

The `access_control_policy` configuration block supports the following arguments:

* `grant` - (Required) Set of `grant` configuration blocks [documented below](#grant).
* `owner` - (Required) Configuration block of the bucket owner's display name and ID [documented below](#owner).

### grant

The `grant` configuration block supports the following arguments:

* `grantee` - (Required) Configuration block for the person being granted permissions [documented below](#grantee).
* `permission` - (Required) Logging permissions assigned to the grantee for the bucket.
    * _Valid values_: `READ`, `WRITE`, `READ_ACP`, `WRITE_ACP`, `FULL_CONTROL`

### owner

The `owner` configuration block supports the following arguments:

* `id` - (Required) The ID of the owner.
* `display_name` - (Optional) The display name of the owner.

### grantee

The `grantee` configuration block supports the following arguments:

* `email_address` - (Optional) Email address of the grantee (S3 Project email). Used only when `type` is `AmazonCustomerByEmail`.
* `id` - (Optional) The canonical user ID of the grantee (S3 User ID). Used only when `type` is `CanonicalUser`.
* `type` - (Required) Type of grantee.
    * _Valid values_:  `CanonicalUser`, `AmazonCustomerByEmail`, `Group`
* `uri` - (Optional) URI of the grantee group. Supported groups are `http://acs.amazonaws.com/groups/global/AllUsers` and `http://acs.amazonaws.com/groups/global/AuthenticatedUsers`. Used only when `type` is `Group`.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The `bucket` and `acl` (if configured) separated by commas (`,`).

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `expected_bucket_owner`.

## Import

S3 bucket ACL can be imported in one of two ways.


If the source bucket is **not configured** with a [canned ACL][canned-acl] (i.e. predefined grant),
the S3 bucket ACL resource should be imported using the `bucket` e.g.,

```
$ terraform import aws_s3_bucket_acl.example bucket-name
```

If the source bucket is **configured** with a [canned ACL][canned-acl] (i.e. predefined grant),
the S3 bucket ACL resource should be imported using the `bucket` and `acl` separated by a comma (`,`), e.g.

```
$ terraform import aws_s3_bucket_acl.example bucket-name,private
```
