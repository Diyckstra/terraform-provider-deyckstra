---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route_table"
description: |-
  Provides a resource to create a VPC routing table.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[route-tables]: https://docs.k2.cloud/en/services/networking/routetables.html
[tf-route]: route.html
[tf-vpn-gateway-route-propagation]: vpn_gateway_route_propagation.html
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_route_table

Provides a resource to create a VPC routing table.

~> **Note on Route Tables and Routes:** Terraform currently
provides both a standalone [`aws_route` resource][tf-route] and a Route Table resource with routes
defined in-line. At this time you cannot use a route table with in-line routes
in conjunction with any route resources. Doing so will cause
a conflict of rule settings and will overwrite rules.

~> **Note on `propagating_vgws` and the `aws_vpn_gateway_route_propagation` resource:**
If the `propagating_vgws` argument is present, it's not supported to _also_
define route propagations using [`aws_vpn_gateway_route_propagation`][tf-vpn-gateway-route-propagation], since
this resource will delete any propagating gateways not explicitly listed in
`propagating_vgws`. Omit this argument when defining route propagation using
the separate resource.

For more information, see the documentation on [Route Tables][route-tables].

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_subnet" "example" {
  availability_zone = "ru-msk-vol52"
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 1, 0)
}

resource "aws_network_interface" "example" {
  subnet_id = aws_subnet.example.id
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id

  route {
    cidr_block           = "10.0.1.0/24"
    network_interface_id = aws_network_interface.example.id
  }

  tags = {
    Name = "example"
  }
}
```

To subsequently remove all managed routes:

```terraform
resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id

  route = []

  tags = {
    Name = "example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the VPC.
* `route` - (Optional) A list of route objects. Their keys are documented below. This argument is processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).
This means that omitting this argument is interpreted as ignoring any existing routes. To remove all managed routes an empty list should be specified. See the example above.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.
* `propagating_vgws` - (Optional) A list of virtual gateways for propagation.

### route Argument Reference

This argument is processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).

One of the following destination arguments must be supplied:

* `cidr_block` - (Required) The CIDR block of the route.

One of the following target arguments must be supplied:

* `gateway_id` - (Optional) ID of an internet gateway or virtual private gateway.
* `instance_id` - (Optional) ID of an EC2 instance.
* `network_interface_id` - (Optional) ID of an EC2 network interface.
* `transit_gateway_id` - (Optional) The ID of the transit gateway.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

~> **Note** Only the target that is entered is exported as a readable
attribute once the route resource is created.

* `id` - ID of the route table.
* `arn` - ARN of the route table.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_gateway_id`, `destination_prefix_list_id`, `ipv6_cidr_block`, `owner_id`, `route.core_network_arn`, `route.egress_only_gateway_id`, `route.nat_gateway_id`, `route.vpc_endpoint_id`, `route.vpc_peering_connection_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `5 minutes`) Used for route creation
- `update` - (Default `2 minutes`) Used for route creation
- `delete` - (Default `5 minutes`) Used for route deletion

## Import

Route Tables can be imported using the route table `id`. For example, to import
route table `rtb-12345678`, use this command:

```
$ terraform import aws_route_table.example rtb-12345678
```

