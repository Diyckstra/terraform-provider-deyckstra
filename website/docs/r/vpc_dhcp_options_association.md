---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc_dhcp_options_association"
description: |-
  Manages a VPC DHCP Options Association.
---

# Resource: aws_vpc_dhcp_options_association

Manages a VPC DHCP Options Association.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the VPC to which we would like to associate a DHCP Options Set.
* `dhcp_options_id` - (Required) ID of the DHCP Options Set to associate to the VPC.

## Remarks

* You can only associate one DHCP Options Set to a given VPC ID.
* Removing the DHCP Options Association automatically sets `default` DHCP Options Set to the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the DHCP Options Set Association.

## Import

DHCP associations can be imported by providing the VPC ID associated with the options:

```
$ terraform import aws_vpc_dhcp_options_association.imported vpc-CFE7ADB5
```
