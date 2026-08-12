---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip"
description: |-
  Provides information about an Elastic IP.
---

[describe-addresses]: https://docs.k2.cloud/en/api/ec2/actions/addresses/DescribeAddresses.html

# Data Source: aws_eip

Provides information about an Elastic IP (EIP).

## Example usage

### Searching by allocation ID

```terraform
data "aws_eip" "by_allocation_id" {
  id = "eipalloc-12345678"
}
```

### Searching by filters

```terraform
data "aws_eip" "by_filter" {
  filter {
    name   = "tag:Name"
    values = ["exampleNameTagValue"]
  }
}
```

### Searching by public IP

```terraform
data "aws_eip" "by_public_ip" {
  public_ip = "1.2.3.4"
}
```

### Searching by tags

```terraform
data "aws_eip" "by_tags" {
  tags = {
    Name = "exampleNameTagValue"
  }
}
```

## Argument reference

The arguments of this data source act as filters for querying the available EIPs.
The given filters must match exactly one EIP whose data will be exported as attributes.

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in the [EC2 API documentation][describe-addresses]
* `id` - (Optional, String) The allocation ID of the specific EIP to retrieve.
* `public_ip` - (Optional, String) The public IP address of the specific EIP to retrieve.
* `tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the required resource.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `association_id` - (String) The ID of the address association with an instance.
* `domain` - (String) The domain in which the EIP is used.
    * _Constraints:_ Always `vpc`
* `id` - (String) The ID of the EIP allocation.
* `instance_id` - (String) The ID of the instance that the address is associated with (if any).
* `network_interface_id` - (String) The ID of the network interface.
* `network_interface_owner_id` - (String) The ID of the project that the network interface belongs to.
* `private_dns` - (String) The private DNS name of the network interface this EIP is attached to.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `private_ip` - (String) The private IP address.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `public_dns` - (String) The public DNS name of the network interface this EIP is attached to.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `public_ip` - (String) The public IP address of the EIP.
* `public_ipv4_pool` - (String) The ID of the EC2 IPv4 address pool.
* `tags` - (Map of strings) Key-value pairs assigned to the EIP.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_ip`, `customer_owned_ip`, `customer_owned_ipv4_pool`.
