---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route"
description: |-
    Provides information about a specific route.
---

# Data Source: aws_route

Provides information about a specific route.

This resource can be used when finding the resource associated with a CIDR. For example, finding the peering connection associated with a CIDR value.

## Example Usage

The following example shows how one might use a CIDR value to find a network interface id and use this to create a data source of that network interface.

```terraform
variable "subnet_id" {}

data "aws_route_table" "selected" {
  subnet_id = var.subnet_id
}

data "aws_route" "route" {
  route_table_id         = data.aws_route_table.selected.id
  destination_cidr_block = "10.0.1.0/24"
}

data "aws_network_interface" "interface" {
  id = data.aws_route.route.network_interface_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available Route in the current region. The given filters must match exactly oneRoute whose data will be exported as attributes.

The following arguments are required:

* `route_table_id` - (Required) ID of the specific Route Table containing the Route entry.

The following arguments are optional:

* `destination_cidr_block` - (Optional) CIDR block of the Route belonging to the Route Table.
* `gateway_id` - (Optional) Gateway ID of the Route belonging to the Route Table.
* `instance_id` - (Optional) Instance ID of the Route belonging to the Route Table.
* `network_interface_id` - (Optional) Network Interface ID of the Route belonging to the Route Table.
* `transit_gateway_id` - (Optional) The ID of the transit gateway.

->  **Unsupported arguments**
These arguments are currently unsupported:

* `carrier_gateway_id` - ID of a carrier gateway. Always `""`.
* `core_network_arn` - ARN of a core network. Always `""`.
* `destination_ipv6_cidr_block` - The destination IPv6 CIDR block. Always `""`.
* `destination_prefix_list_id` - ID of a managed prefix list destination of the route. Always `""`.
* `egress_only_gateway_id` - ID of a VPC Egress Only Internet Gateway. Always `""`.
* `local_gateway_id` - ID of an Outpost local gateway. Always `""`.
* `nat_gateway_id` - ID of a VPC NAT gateway. Always `""`.
* `vpc_endpoint_id` - ID of a VPC Endpoint. Always `""`.
* `vpc_peering_connection_id` - ID of a VPC peering connection. Always `""`.

## Attributes Reference

All the argument attributes are also exported as result attributes when there is data available. For example, the `vpc_peering_connection_id` field will be empty when the route is attached to a Network Interface.
