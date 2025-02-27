---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc"
description: |-
  Manages a VPC.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[vpc]: https://docs.k2.cloud/en/services/networking/privatecloud.html

# Resource: aws_vpc

Manages a VPC.
For more information, see the documentation on [VPC][vpc].

## Example Usage

Basic usage:

```terraform
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
```

Basic usage with tags:

```terraform
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "main"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) The IPv4 CIDR block for the VPC.
* `enable_dns_support` - (Optional) A boolean flag to enable/disable DNS support in the VPC. Defaults true.
* `tags` - (Optional) Map of tags to assign to the resource. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of VPC
* `id` - ID of the VPC
* `main_route_table_id` - ID of the main route table associated with
     this VPC. Note that you can change a VPC's main route table by using an
     [`aws_main_route_table_association`](main_route_table_association.md).
* `default_network_acl_id` - ID of the network ACL created by default on VPC creation
* `default_security_group_id` - ID of the security group created by default on VPC creation
* `default_route_table_id` - ID of the route table created by default on VPC creation
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`assign_generated_ipv6_cidr_block`, `enable_classiclink`, `enable_classiclink_dns_support`, `enable_dns_hostnames`, `instance_tenancy`, `ipv4_ipam_pool_id`, `ipv4_netmask_length`, `ipv6_association_id`, `ipv6_cidr_block`, `ipv6_cidr_block_network_border_group`, `ipv6_ipam_pool_id`, `ipv6_netmask_length`, `owner_id`.

## Import

VPCs can be imported using the `vpc id`, e.g.,

```
$ terraform import aws_vpc.test_vpc vpc-12345678
```

