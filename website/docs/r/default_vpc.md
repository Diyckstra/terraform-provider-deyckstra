---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_default_vpc"
description: |-
  Manage the default VPC.
---

# Resource: aws_default_vpc

Manages the default VPC.

**This is an advanced resource** and has special caveats to be aware of when using it. Please read this document in its entirety before using this resource.

The `aws_default_vpc` resource behaves differently from normal resources in that if a default VPC exists, Terraform does not _create_ this resource, but instead "adopts" it into management.
If no default VPC exists, Terraform creates a new default VPC, which leads to the implicit creation of other resources.
By default, `terraform destroy` does not delete the default VPC but does remove the resource from Terraform state.
Set the `force_destroy` argument to `true` to delete the default VPC.

## Example Usage

Basic usage with tags:

```terraform
resource "aws_default_vpc" "default" {
  tags = {
    Name = "Default VPC"
  }
}
```

## Argument Reference

The arguments of an `aws_default_vpc` differ slightly from those of [`aws_vpc`](vpc.md):

* The `cidr_block` and `instance_tenancy` arguments become computed attributes
* The default value for `enable_dns_hostnames` is `true`

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `cidr_block` - The primary IPv4 CIDR block for the VPC
* `instance_tenancy` - The allowed tenancy of instances launched into the VPC

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `force_destroy`.

## Import

Default VPCs can be imported using the `vpc id`, e.g.,

```
$ terraform import aws_default_vpc.default vpc-12345678
```
