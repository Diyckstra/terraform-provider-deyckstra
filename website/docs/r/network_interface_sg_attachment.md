---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interface_sg_attachment"
description: |-
  Associates a security group with a network interface.
---

# Resource: aws_network_interface_sg_attachment

This resource attaches a security group to an elastic network interface (ENI).
It can be used to attach a security group to any existing ENI, be it a
secondary ENI or one attached as the primary interface on an instance.

~> **Note on instances, interfaces, and security groups:** Terraform currently
provides the capability to assign security groups via the [`aws_instance`](instance.md)
and the [`aws_network_interface`](network_interface.md) resources. Using this resource in
conjunction with security groups provided in-line in those resources will cause
conflicts, and will lead to spurious diffs and undefined behavior - please use
one or the other.

## Example Usage

The following provides a very basic example of setting up an instance (provided
by `instance`) in the default security group, creating a security group
(provided by `sg`) and then attaching the security group to the instance's
primary network interface via the `aws_network_interface_sg_attachment` resource,
named `sg_attachment`:

```terraform
resource "aws_instance" "instance" {
  instance_type = "m1.micro"
  ami           = "cmi-12345678" # add image id, change instance type if needed

  tags = {
    type = "terraform-test-instance"
  }
}

resource "aws_security_group" "sg" {
  tags = {
    type = "terraform-test-security-group"
  }
}

resource "aws_network_interface_sg_attachment" "sg_attachment" {
  security_group_id    = aws_security_group.sg.id
  network_interface_id = aws_instance.instance.primary_network_interface_id
}
```

## Argument Reference

* `security_group_id` - (Required) The ID of the security group.
* `network_interface_id` - (Required) The ID of the network interface to attach to.

## Attributes Reference

No additional attributes are exported.
