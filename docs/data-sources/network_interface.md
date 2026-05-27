---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interface"
description: |-
  Provides information about a network interface.
---

[describe-network-interfaces]: https://docs.k2.cloud/en/api/ec2/actions/network_interfaces/DescribeNetworkInterfaces.html

# aws_network_interface

Provides information about a network interface.

## Example usage

### Basic example

```terraform
data "aws_network_interface" "example" {
  id = "eni-xxxxxxxx"
}
```

### Specific Example

```terraform
data "aws_network_interface" "example" {
  filter {
    name   = "tag:Name"
    values = ["example"]
  }
}
```

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-network-interfaces]
* `id` - (Optional, String) The ID of the network interface.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the network interface.
* `association` - ([Block](#association)) The association information for an Elastic IP address (IPv4) associated with the network interface.
* `attachment` - ([Block](#attachment)) The list of network interface attachments.
* `availability_zone` - (String) The name of the availability zone.
* `description` - (String) The description of the network interface.
* `mac_address` - (String) The MAC address.
* `owner_id` - (String) The ID of the project that owns the network interface.
* `private_dns_name` - (String) The internal domain name of the network interface.
    * _Example:_ `ip-123-4-5-6.vpc-12345678.internal`
* `private_ip` - (String) The private IPv4 address of the network interface within the subnet.
* `private_ips` - (List of strings) The private IPv4 addresses associated with the network interface.
* `security_groups` - (List of strings) The list of security groups for the network interface.
* `subnet_id` - (String) The ID of the subnet, if the network interface is created in a subnet, or the ID of the switch, if the network interface is created in a switch.
* `tags` - (Map of strings) Key-value pairs assigned to the resource.
* `vpc_id` - (String) The ID of the VPC.

#### association

* `allocation_id` - (String) The ID of the allocation.
* `association_id` - (String) The ID of the association.
* `public_dns_name` - (String) The public DNS name.
* `public_ip` - (String) The Elastic IP address bound to the network interface.

#### attachment

* `attachment_id` - (String) The ID of the network interface attachment.
* `device_index` - (Integer) The index of the network interface.
* `instance_id` - (String) The ID of the instance.
* `instance_owner_id` - (String) The ID of the project that owns the instance.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`association.carrier_ip`, `association.customer_owned_ip`, `association.ip_owner_id`, `interface_type`, `ipv6_addresses`, `outpost_arn`, `requester_id`.
