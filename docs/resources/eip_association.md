---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip_association"
description: |-
  Manages an EIP association.
---

# Resource: aws_eip_association

Manages an EIP association as a top level resource, to associate and
disassociate Elastic IPs from instances and network interfaces.

~> **Note** `aws_eip_association` is useful in scenarios where EIPs are either
pre-existing or distributed to customers or users and therefore cannot be changed.

## Example usage

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
  vpc = true
}

resource "aws_eip_association" "eip_assoc" {
  depends_on = [aws_internet_gateway.example_igw]

  instance_id   = aws_instance.example_instance.id
  allocation_id = aws_eip.example_eip.id
}
```

## Argument reference

The following arguments are supported:

* `allocation_id` - (Optional, Forces new resource, String) The ID of the allocation.
    * _Constraints:_ Required if the `public_ip` is not supplied
* `allow_reassociation` - (Optional, Forces new resource, Boolean) Indicates whether to allow an Elastic IP to be re-associated.
    * _Default value:_ `true`
* `instance_id` - (Optional, Forces new resource, String) The ID of the instance.
    * _Constraints:_ Required if the `network_interface_id` is not supplied
* `network_interface_id` - (Optional, Forces new resource, String) The ID of the network interface.
    * _Constraints:_ Required if the `instance_id` is not supplied
* `public_ip` - (Optional, Forces new resource, String) The Elastic IP address.
    * _Constraints:_ Required if the `allocation_id` is not supplied

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `allocation_id` - (String) The ID of the allocation.
* `id` - (String) The ID of the association.
* `instance_id` - (String) The ID of the instance that the address is associated with.
* `network_interface_id` - (String) The ID of the network interface.
* `private_ip_address` - (String) The private IP address associated with the Elastic IP address.
* `public_ip` - (String) The public IP address of the Elastic IP.

~> **Note** The `private_ip_address` value is set by the platform when the Elastic IP is associated and cannot be specified in configuration files.

## Timeouts

Timeouts usage for EIP association is not currently supported.

## Import

EIP associations can be imported using IDs of their associations, for example:

```
$ terraform import aws_eip_association.test eipassoc-12345678
```
