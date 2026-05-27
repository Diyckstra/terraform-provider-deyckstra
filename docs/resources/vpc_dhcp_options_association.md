---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc_dhcp_options_association"
description: |-
  Manages a VPC DHCP options association.
---

# Resource: aws_vpc_dhcp_options_association

Manages a VPC DHCP options association.

## Example usage

### Basic example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_vpc_dhcp_options" "example" {
  domain_name_servers = ["8.8.8.8", "8.8.4.4"]
}

resource "aws_vpc_dhcp_options_association" "example" {
  vpc_id          = aws_vpc.example.id
  dhcp_options_id = aws_vpc_dhcp_options.example.id
}
```

## Argument reference

The following arguments are supported:

* `dhcp_options_id` - (Required, Editable, String) The ID of the DHCP options set to associate to the VPC.
* `vpc_id` - (Required, Forces new resource, String) The ID of the VPC to associate with the DHCP options set.

## Notes

* You can only associate one DHCP options set to a given VPC ID.
* Removing the DHCP options association automatically sets the `default` DHCP options set to the VPC.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) Consists of the ID of the DHCP options set and VPC ID separated by the hyphen.
    * _Example:_ `dopt-12345678-vpc-12345678`

## Timeouts

Timeouts usage for the DHCP options associations is not currently supported.

## Import

DHCP associations can be imported by providing the VPC ID associated with the DHCP options set:

```
$ terraform import aws_vpc_dhcp_options_association.example vpc-12345678
```
