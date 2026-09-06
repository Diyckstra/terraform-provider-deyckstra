---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip_association"
description: |-
  Manages an EIP association.
---

# Resource: aws_eip_association

Manages an Elastic IP (EIP) association as a top level resource, to associate and
disassociate EIPs from instances and network interfaces.

-> **Note** `aws_eip_association` is useful in scenarios where EIPs are either
pre-existing or distributed to customers or users and therefore cannot be changed.

## Example usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.0.0.0/24"
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.example.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.example.id
}

resource "aws_instance" "example" {
  ami           = "cmi-12345678"
  instance_type = "m1.micro"
  subnet_id     = aws_subnet.example.id
}

resource "aws_eip" "example" {}

resource "aws_eip_association" "example" {
  depends_on = [aws_internet_gateway.example]

  instance_id   = aws_instance.example.id
  allocation_id = aws_eip.example.id
}
```

## Argument reference

The following arguments are supported:

* `allocation_id` - (Optional, Forces new resource, String) The ID of the allocation.
    * _Constraints:_ Required if the `public_ip` is not supplied
* `allow_reassociation` - (Optional, Forces new resource, Boolean) Indicates whether to allow an EIP to be re-associated. Reassociation is automatic, but you can specify `false` to ensure the operation fails if the EIP is already associated with another resource.
* `instance_id` - (Optional, Forces new resource, String) The ID of the instance.
    * _Constraints:_ Required if the `network_interface_id` is not supplied
* `network_interface_id` - (Optional, Forces new resource, String) The ID of the network interface.
    * _Constraints:_ Required if the `instance_id` is not supplied
* `public_ip` - (Optional, Forces new resource, String) The EIP address.
    * _Constraints:_ Required if the `allocation_id` is not supplied

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attribute is exported:

* `id` - (String) The ID of the association.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`private_ip_address`.

## Timeouts

Timeouts usage for EIP association is not currently supported.

## Import

EIP associations can be imported using IDs of their associations, for example:

```
$ terraform import aws_eip_association.example eipassoc-12345678
```
