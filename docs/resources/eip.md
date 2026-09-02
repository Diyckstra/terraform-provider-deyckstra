---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip"
description: |-
  Manages an Elastic IP.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[elastic-ips]: https://docs.k2.cloud/en/services/networking/addresses/operations.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_eip

Manages an Elastic IP (EIP). For more information about EIPs, see [user documentation][elastic-ips].

## Example usage

### EIP associated with an instance

```terraform
resource "aws_vpc" "example_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "example_subnet" {
  vpc_id     = aws_vpc.example_vpc.id
  cidr_block = "10.0.0.0/24"
}

resource "aws_internet_gateway" "example_igw" {
  vpc_id = aws_vpc.example_vpc.id
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.example_vpc.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.example_igw.id
}

resource "aws_instance" "example_instance" {
  ami           = "cmi-12345678"
  instance_type = "m1.micro"
  subnet_id     = aws_subnet.example_subnet.id
}

resource "aws_eip" "example_eip" {
  depends_on = [aws_internet_gateway.example_igw]

  vpc      = true
  instance = aws_instance.example_instance.id

  tags = {
    Name = "Example EIP"
  }
}
```

### EIP associated with a network interface

~> **Note**
This example uses the same VPC and subnet as in the [EIP associated with an instance example](#eip-associated-with-an-instance).

```terraform
resource "aws_network_interface" "example_interface" {
  subnet_id = aws_subnet.example_subnet.id
}

resource "aws_eip" "example_eip" {
  depends_on = [aws_internet_gateway.example_igw]

  vpc               = true
  network_interface = aws_network_interface.example_interface.id

  tags = {
    Name = "Example EIP"
  }
}
```

### Allocating EIP from the BYOIP pool

```terraform
resource "aws_eip" "example_eip" {
  vpc              = true
  public_ipv4_pool = "ipv4pool-ec2-012345"
}
```

## Argument reference

The following arguments are supported:

* `address` - (Optional) IP address from an EC2 BYOIP pool.
    _Constraints:_ This option is only available for VPC EIPs
* `instance` - (Optional) The ID of the EC2 instance.
* `network_interface` - (Optional) The ID of the network interface to associate with.
* `public_ipv4_pool` - (Optional) The ID of the EC2 IPv4 address pool.
    * _Constraints:_ This option is only available for VPC EIPs
* `tags` - (Optional) Map of tags to assign to the EIP. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.
    * _Constraints:_ Tags can only be applied to EIPs in a VPC
* `vpc` - (Optional) Boolean if the EIP is in a VPC or not.

~> **Note** You can specify either the ID of `instance` or the ID of `network_interface`, but not both.

~> **Note** If both `public_ipv4_pool` and `address` are specified, `address` will be used in the case both options are defined as API only requires one or the other.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `allocation_id` - (String) The ID representing the allocation of the IP address.
* `association_id` - (String) The ID representing the association of the allocation of the IP address with an instance or a private IP address.
* `domain` - (String) Indicates if this EIP is for use in VPC (`vpc`).
* `id` - (String) The ID of the EIP allocation.
* `private_ip` - (String) The private IP address.
* `public_ip` - (String) The public IP address.
* `tags_all` - (Map of strings) Key-value pairs assigned to the Elastic IP, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`associate_with_private_ip`, `carrier_ip`, `customer_owned_ip`, `customer_owned_ipv4_pool`, `network_border_group`, `private_dns`, `public_dns`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `read` - (Default `15 minutes`) Used when querying for information about EIPs.
* `update` - (Default `5 minutes`) Used when updating an EIP.
* `delete` - (Default `3 minutes`) Used when deleting an EIP.

## Import

EIPs in a VPC can be imported using their allocation ID, for example:

```
$ terraform import aws_eip.example_eip eipalloc-1234567
```

EIPs can be imported using their public IP, for example:

```
$ terraform import aws_eip.example_eip 1.1.1.1
```
