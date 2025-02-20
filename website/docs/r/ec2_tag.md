---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ec2_tag"
description: |-
  Manages an individual EC2 resource tag
---

[ignore-tags]: https://www.terraform.io/docs/providers/aws/index.html#ignore_tags-configuration-block

# Resource: aws_ec2_tag

Manages an individual EC2 resource tag. This resource should only be used when Terraform was not used to create EC2 resources (e.g., AMIs).

~> **Note** This tagging resource should not be combined with the Terraform resource for managing the parent resource. For example, using `aws_vpc` and `aws_ec2_tag` to manage tags of the same VPC will cause a perpetual difference where the `aws_vpc` resource will try to remove the tag being added by the `aws_ec2_tag` resource.

~> **Note** This tagging resource does not use the [provider `ignore_tags` configuration][ignore-tags].

## Example Usage

```terraform
resource "aws_ec2_tag" "example" {
  resource_id = "vol-12345678"
  key         = "tag-from-tf"
  value       = "tf-tag"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` – (Required) The ID of the EC2 resource to manage the tag for.
* `key` – (Required) The tag name.
* `value` – (Required) The value of the tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` – EC2 resource identifier and key, separated by a comma (`,`)

## Import

`aws_ec2_tag` can be imported using the EC2 resource identifier and key separated by a comma (`,`), e.g.,

```
$ terraform import aws_ec2_tag.example tgw-attach-12345678,Name
```
