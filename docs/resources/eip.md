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

### Associating a single EIP with an instance

```terraform
resource "aws_eip" "example" {
  instance = "i-12345678"
}
```

### Attaching an EIP to an instance with a pre-assigned private IP

```terraform
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "tf_test_subnet" {
  vpc_id     = aws_vpc.default.id
  cidr_block = "10.0.0.0/24"
}

resource "aws_internet_gateway" "default" {
  vpc_id = aws_vpc.default.id
}

resource "aws_instance" "foo" {
  ami           = "cmi-12345678"
  instance_type = "m1.micro"

  private_ip = "10.0.0.12"
  subnet_id  = aws_subnet.tf_test_subnet.id
}

resource "aws_eip" "bar" {
  instance                  = aws_instance.foo.id
  associate_with_private_ip = "10.0.0.12"
  depends_on                = [aws_internet_gateway.default]
}
```

### Attaching an EIP to a network interface

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.1.0.0/24"
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_network_interface" "example" {
  subnet_id = aws_subnet.example.id
}

resource "aws_eip" "by_interface" {
  network_interface = aws_network_interface.example.id
  depends_on        = [aws_internet_gateway.example]
}
```

### Attaching an EIP to the second network interface of an instance

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.2.0.0/16"
}

resource "aws_subnet" "primary" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.2.0.0/24"
}

resource "aws_subnet" "secondary" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.2.1.0/24"
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_instance" "example" {
  ami           = "cmi-12345678"
  instance_type = "m1.micro"

  private_ip = "10.2.0.10"
  subnet_id  = aws_subnet.primary.id
}

resource "aws_network_interface" "secondary" {
  subnet_id   = aws_subnet.secondary.id
  private_ips = ["10.2.1.20"]
}

resource "aws_network_interface_attachment" "secondary" {
  instance_id          = aws_instance.example.id
  network_interface_id = aws_network_interface.secondary.id
  device_index         = 1
}

# The EIP is attached to 10.2.1.20, not to the instance's primary 10.2.0.10.
resource "aws_eip" "by_second_interface" {
  network_interface         = aws_network_interface.secondary.id
  associate_with_private_ip = "10.2.1.20"

  depends_on = [
    aws_internet_gateway.example,
    aws_network_interface_attachment.secondary,
  ]
}
```

### Allocating an EIP from the BYOIP pool

```terraform
resource "aws_eip" "byoip-ip" {
  public_ipv4_pool = "ipv4pool-ec2-012345"
}
```

## Argument reference

The following arguments are optional:

* `address` - (Optional, Forces new resource, String) An IP address from an EC2 BYOIP pool.
* `associate_with_private_ip` - (Optional, Editable, String) A user-specified primary or secondary private IP address to associate with the EIP.
    * _Constraints:_ If no private IP address is specified, the EIP is associated with the primary private IP address
* `instance` - (Optional, Editable, String) The ID of the EC2 instance.
    * _Constraints:_ Conflicts with the `network_interface` argument.
      The EIP is associated with the primary network interface of the instance.
      If the instance has more than one network interface, use the `network_interface` argument to specify another one
* `network_interface` - (Optional, Editable, String) The ID of the network interface to associate with.
    * _Constraints:_ Conflicts with the `instance` argument
* `public_ipv4_pool` - (Optional, Forces new resource, String) The ID of the EC2 IPv4 address pool.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the EIP. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `vpc` - (Optional, Editable, Boolean, **Deprecated**) Indicates whether the EIP is in a VPC.

~> **Note** The argument `vpc` is deprecated.
Its value is ignored: all EIPs are for use in a VPC.

~> **Note** If both `public_ipv4_pool` and `address` are specified, `address` will be used in the case both options are defined as API only requires one or the other.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `allocation_id` - (String) The ID of the allocation of the IP address.
* `association_id` - (String) The ID of the address association with an instance.
* `domain` - (String) The domain in which the EIP is used.
    * _Constraints:_ Always `vpc`
* `id` - (String) The ID of the EIP allocation.
* `private_dns` - (String) The private DNS name of the network interface this EIP is attached to.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `private_ip` - (String) The private IP address.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `public_dns` - (String) The public DNS name of the network interface this EIP is attached to.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `public_ip` - (String) The public IP address.
* `tags_all` - (Map of strings) Key-value pairs assigned to the EIP, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_ip`, `customer_owned_ip`, `customer_owned_ipv4_pool`, `network_border_group`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `read` - (Default `15 minutes`) Used when waiting for a newly created EIP to become visible.
* `update` - (Default `5 minutes`) Used when associating the EIP with an instance or a network interface.
* `delete` - (Default `3 minutes`) Used when releasing an EIP.

## Import

The EIP can be imported using the allocation ID, for example:

```
$ terraform import aws_eip.bar eipalloc-1234567
```

The EIP can also be imported using the public IP, for example:

```
$ terraform import aws_eip.bar 1.1.1.1
```
