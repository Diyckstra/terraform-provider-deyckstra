---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interface"
description: |-
  Manages an elastic network interface (ENI).
---


[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[network-interfaces]: https://docs.k2.cloud/en/services/networking/interfaces/operations.html

# Resource: aws_network_interface

Manages an elastic network interface (ENI) resource.
For more information, see the documentation on [Network interfaces][network-interfaces].


## Example Usage

```terraform
resource "aws_network_interface" "example" {
  subnet_id   = "subnet-12345678"
  private_ips = ["10.0.31.50"]

  attachment {
    instance     = "i-12345678"
    device_index = 1
  }
}
```

## Argument Reference

The following arguments are required:

* `subnet_id` - (Required) Subnet ID to create the ENI in.

The following arguments are optional:

* `attachment` - (Optional) Configuration block to define the attachment of the ENI. See [Attachment](#attachment) below for more details!
* `description` - (Optional) Description for the network interface.
* `private_ip_list` - (Optional) List of private IPs to assign to the ENI in sequential order. One value only.
* `private_ips` - (Optional) List of private IPs to assign to the ENI without regard to order. One value only.
* `security_groups` - (Optional) List of security group IDs to assign to the ENI.
* `source_dest_check` - (Optional) Whether to enable source destination checking for the ENI. Default true.
* `tags` - (Optional) Map of tags to assign to the resource. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

### Attachment

The `attachment` block supports the following:

* `instance` - (Required) ID of the instance to attach to.
* `device_index` - (Required) Integer to define the devices index.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the network interface.
* `id` - ID of the network interface.
* `mac_address` - MAC address of the network interface.
* `owner_id` - The project ID.
* `private_dns_name` - Private DNS name of the network interface (IPv4).
* `tags_all` -Map of tags to assign to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

### Supported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`interface_type`, `ipv4_prefix_count`, `ipv4_prefixes`, `ipv6_address_count`, `ipv6_address_list_enable`, `ipv6_address_list`, `ipv6_addresses`, `ipv6_prefix_count`, `ipv6_prefixes`, `private_ip_list_enable`, `private_ips_count`.

## Import

Network Interfaces can be imported using the `id`, e.g.,

```
$ terraform import aws_network_interface.test eni-12345678
```

