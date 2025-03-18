---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_subnet"
description: |-
  Manages a subnet.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[subnets]: https://docs.k2.cloud/en/services/networking/subnets.html
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_subnet

Manages a subnet.
For more information, see the documentation on [Subnets][subnets].

## Example Usage

### Basic Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.1.1.0/24"

  tags = {
    Name = "Main"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) AZ for the subnet.
* `cidr_block` - (Required) The IPv4 CIDR block for the subnet.
* `map_public_ip_on_launch` - (Optional) Indicates whether public IP addresses will be associated with instances created in this subnet. Addresses are associated only if there are available allocated Elastic IP addresses. Default is `false`.
* `tags` - (Optional) Map of tags to assign to the subnet. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.
* `vpc_id` - (Required) The VPC ID.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the subnet
* `tags_all` - Map of tags assigned to the subnet, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`arn`, `assign_ipv6_address_on_creation`, `availability_zone_id`, `customer_owned_ipv4_pool`, `enable_dns64`, `enable_resource_name_dns_aaaa_record_on_launch`, `enable_resource_name_dns_a_record_on_launch`, `ipv6_cidr_block_association_id`, `ipv6_native`, `map_customer_owned_ip_on_launch`, `outpost_arn`, `owner_id`, `private_dns_hostname_type_on_launch`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10m`) How long to wait for a subnet to be created.
- `delete` - (Default `20m`) How long to wait for a subnet to be deleted.

## Import

Subnets can be imported using the `subnet id`, e.g.,

```
$ terraform import aws_subnet.public_subnet subnet-12345678
```
