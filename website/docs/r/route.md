---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route"
description: |-
  Creates a routing entry in a VPC routing table.
---

[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_route

Creates a routing table entry (a route) in a VPC routing table.

~> **Note on route tables and routes:** Terraform currently provides both a standalone route resource and a [`aws_route_table`][tf-route-table] resource with routes defined in-line. At this time you cannot use a route table with in-line routes in conjunction with any route resources. Doing so will cause a conflict of rule settings and will overwrite rules.

## Example Usage

```terraform
variable route_table_id {}
variable instance_id {}

resource "aws_route" "example" {
  route_table_id         = var.route_table_id
  destination_cidr_block = "10.0.0.0/22"
  instance_id            = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required) ID of the routing table.

One of the following destination arguments must be supplied:

* `destination_cidr_block` - (Required) The destination CIDR block.

One of the following target arguments must be supplied:

* `gateway_id` - (Optional) ID of an internet gateway or virtual private gateway.
* `instance_id` - (Optional) ID of an EC2 instance.
* `network_interface_id` - (Optional) ID of an EC2 network interface.
* `transit_gateway_id` - (Optional) The ID of the transit gateway.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

~> **Note** Only the arguments that are configured (one of the above) will be exported as an attribute once the resource is created.

* `id` - Route identifier computed from the routing table identifier and route destination.
* `instance_owner_id` - The project ID that owns the EC2 instance.
* `origin` - How the route was created - `CreateRouteTable`, `CreateRoute` or `EnableVgwRoutePropagation`.
* `state` - The state of the route - `active` or `blackhole`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_gateway_id`, `core_network_arn`, `destination_ipv6_cidr_block`, `destination_prefix_list_id`, `egress_only_gateway_id`, `local_gateway_id`, `nat_gateway_id`, `vpc_endpoint_id`, `vpc_peering_connection_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `5 minutes`) Used for route creation
- `update` - (Default `2 minutes`) Used for route creation
- `delete` - (Default `5 minutes`) Used for route deletion

## Import

Individual routes can be imported using `ROUTETABLEID_DESTINATION`.

For example, import a route in route table `rtb-12345678` with an IPv4 destination CIDR of `10.1.0.0/16` like this:

```console
$ terraform import aws_route.my_route rtb-12345678_10.1.0.0/16
```
