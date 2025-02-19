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

## Example Usage

```terraform
resource "aws_eip_association" "eip_assoc" {
  instance_id   = aws_instance.web.id
  allocation_id = aws_eip.example.id
}

resource "aws_instance" "web" {
  ami               = "cmi-12345678" # add image id, change instance type if needed
  availability_zone = "ru-msk-vol52"
  instance_type     = "m1.micro"

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_eip" "example" {
  vpc = true
}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Optional) The allocation ID. Required, if `public_ip` is not supplied.
* `allow_reassociation` - (Optional, Boolean) Whether to allow an Elastic IP to
be re-associated. Defaults to `true`.
* `instance_id` - (Optional) The ID of the instance. Required, if `network_interface_id` is not supplied.
* `network_interface_id` - (Optional) The ID of the network interface. Required, if `instance_id` is not supplied.
* `public_ip` - (Optional) The Elastic IP address. Required, if `allocation_id` is not supplied.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `association_id` - The ID that represents the association of the Elastic IP address with an instance.
* `allocation_id` - The allocation ID.
* `instance_id` - The ID of the instance that the address is associated with.
* `network_interface_id` - The ID of the network interface.
* `private_ip_address` - The private IP address associated with the Elastic IP address.
* `public_ip` - Public IP address of Elastic IP.

## Import

EIP associations can be imported using their association ID.

```
$ terraform import aws_eip_association.test eipassoc-12345678
```
