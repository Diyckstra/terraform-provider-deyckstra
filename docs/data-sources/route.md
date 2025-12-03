---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route"
description: |-
  Provides information about a route.
---

# Data Source: aws_route

Provides information about a route.

This resource can be used when finding the resource associated with a CIDR. For example, finding the peering connection associated with a CIDR value.

## Example Usage

The following example shows how one might use a CIDR value to find the ID of a network interface and use this to create a data source of that network interface.

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

The arguments of this data source act as filters for querying the available route in the current region. The given filters must match exactly oneRoute whose data will be exported as attributes.

* `route_table_id` - (Required) ID of the specific route table containing the route entry.
* `destination_cidr_block` - (Optional) CIDR block of the route belonging to the route table.
* `gateway_id` - (Optional) Gateway ID of the route belonging to the route table.
* `instance_id` - (Optional) Instance ID of the route belonging to the route table.
* `network_interface_id` - (Optional) Network interface ID of the route belonging to the route table.
* `transit_gateway_id` - (Optional) The ID of the transit gateway.

## Attribute Reference

All the argument attributes are also exported as result attributes.

### Unsupported attributes

~> **Note** These arguments may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following arguments are not currently supported:

`carrier_gateway_id`, `core_network_arn`, `destination_ipv6_cidr_block`, `destination_prefix_list_id`, `egress_only_gateway_id`, `local_gateway_id`, `nat_gateway_id`, `vpc_endpoint_id`, `vpc_peering_connection_id`.
