---
subcategory: "S3 (Simple Storage Service)"
layout: "aws"
page_title: "aws_s3_bucket_website_configuration"
description: |-
  Manages an S3 bucket website configuration.
---

[hosting-website]: https://docs.k2.cloud/en/services/object_storage/operations.html#objectstoragestaticwebsitesmanual

# Resource: aws_s3_bucket_website_configuration

Manages an S3 bucket website configuration.
For more information about hosting websites on S3, see [user documentation][hosting-website].

## Example Usage

```terraform
resource "aws_s3_bucket_website_configuration" "example" {
  bucket = aws_s3_bucket.example.bucket

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }

  routing_rule {
    condition {
      key_prefix_equals = "docs/"
    }
    redirect {
      replace_key_prefix_with = "documents/"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, Forces new resource) The name of the bucket.
* `error_document` - (Optional, Conflicts with `redirect_all_requests_to`) The name of the error document for the website [detailed below](#error_document).
* `index_document` - (Optional) The name of the index document for the website [detailed below](#index_document).
    * _Constraints_: Required if `redirect_all_requests_to` is not specified
* `redirect_all_requests_to` - (Optional) The redirect behavior for every request to this bucket's website endpoint [detailed below](#redirect_all_requests_to).
    * _Constraints_: Required if `index_document` is not specified. Conflicts with `error_document`, `index_document`, and `routing_rule`.
* `routing_rule` - (Optional, Conflicts with `redirect_all_requests_to`) List of rules that define when a redirect is applied and the redirect behavior [detailed below](#routing_rule).

### error_document

The `error_document` configuration block supports the following arguments:

* `key` - (Required) The object key name to use when a 4XX class error occurs.

### index_document

The `index_document` configuration block supports the following arguments:

* `suffix` - (Required) A suffix that is appended to a request that is for a directory on the website endpoint.
For example, if the suffix is `index.html` and you make a request to `samplebucket/images/`, the data that is returned will be for the object with the key name `images/index.html`.
The suffix must not be empty and must not include a slash character.

### redirect_all_requests_to

The `redirect_all_requests_to` configuration block supports the following arguments:

* `host_name` - (Required) Name of the host where requests are redirected.
* `protocol` - (Optional) Protocol to use when redirecting requests. The default is the protocol that is used in the original request.
    * _Valid values_:  `http`, `https`

### routing_rule

The `routing_rule` configuration block supports the following arguments:

* `condition` - (Optional) A configuration block for describing a condition that must be met for the specified redirect to apply [detailed below](#condition).
* `redirect` - (Required) A configuration block for redirect information [detailed below](#redirect).

### condition

The `condition` configuration block supports the following arguments:

* `http_error_code_returned_equals` - (Optional) The HTTP error code when the redirect is applied. If specified with `key_prefix_equals`, then both must be true for the redirect to be applied.
    * _Constraints_: Required if `key_prefix_equals` is not specified
* `key_prefix_equals` - (Optional) The object key name prefix when the redirect is applied. If specified with `http_error_code_returned_equals`, then both must be true for the redirect to be applied.
    * _Constraints_: Required if `http_error_code_returned_equals` is not specified

### redirect

The `redirect` configuration block supports the following arguments:

* `host_name` - (Optional) The host name to use in the redirect request.
* `protocol` - (Optional) Protocol to use when redirecting requests. The default is the protocol that is used in the original request.
    * _Valid values_:  `http`, `https`
* `replace_key_prefix_with` - (Optional, Conflicts with `replace_key_with`) The object key prefix to use in the redirect request. For example, to redirect requests for all pages with prefix `docs/` (objects in the `docs/` folder) to `documents/`, you can set a `condition` block with `key_prefix_equals` set to `docs/` and in the `redirect` set `replace_key_prefix_with` to `/documents`.
* `replace_key_with` - (Optional, Conflicts with `replace_key_prefix_with`) The specific object key to use in the redirect request. For example, redirect request to `error.html`.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The `bucket`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`expected_bucket_owner`, `redirect.http_redirect_code`, `website_domain`, `website_endpoint`.

## Import

S3 bucket website configuration can be imported using the `bucket` e.g.,

```
$ terraform import aws_s3_bucket_website_configuration.example bucket-name
```
