---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route_table"
description: |-
  Creates a VPC routing table.
---

[attribute-as-blocks]: https://www.terraform.io/docs/configuration/attr-as-blocks.html
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[route-tables]: https://docs.k2.cloud/en/services/networking/routetables.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts
[sdk-diff-issue]: https://github.com/hashicorp/terraform/issues/21901

# Resource: aws_route_table

Creates a VPC routing table.

~> **Note on route tables and routes** Terraform currently
provides both a standalone [`aws_route` resource](route.md) and a route table resource with routes
defined inline. At this time you cannot use a route table with inline routes
in conjunction with any route resources. Doing so will cause
a conflict of rule settings and will overwrite rules.

~> **Note on `propagating_vgws` and the `aws_vpn_gateway_route_propagation` resource**
If the `propagating_vgws` argument is present, it's not supported to _also_
define route propagations using [`aws_vpn_gateway_route_propagation`](vpn_gateway_route_propagation.md), since
this resource will delete any propagating gateways not explicitly listed in
`propagating_vgws`. Omit this argument when defining route propagation using
the separate resource.

For more information about route tables, see the documentation on [Route tables][route-tables].

## Example usage

### Basic example

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

### Specific example: removing all managed routes subsequently

~> **Note** This example deletes routes created in the [example above](#basic-example).

```terraform
resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id

  route = []

  tags = {
    Name = "example"
  }
}
```

## Argument reference

The following arguments are required:

* `vpc_id` - (Required, Forces new resource, String) The ID of the VPC.

The following arguments are optional:

* `propagating_vgws` - (Optional, Editable, List of strings) The list of virtual gateways for propagation.
* `route` - (Optional, Editable, [Block](#route)) One or more route objects. This argument is processed in [attribute-as-blocks mode][attribute-as-blocks].
It means that omitting this argument is interpreted as ignoring any existing routes. To remove all managed routes an empty list should be specified. See the [example above](#specific-example-removing-all-managed-routes-subsequently).
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource. If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

### route

The following destination argument must be supplied:

* `cidr_block` - (Required, Editable, String) The CIDR block of the route.

One of the following target arguments must be supplied:

* `gateway_id` - (Optional, Editable, String) The ID of the internet gateway.
* `nat_gateway_id` - (Optional, Editable, String) The ID of the NAT gateway.
* `network_interface_id` - (Optional, Editable, String) The ID of the network interface.
* `transit_gateway_id` - (Optional, Editable, String) The ID of the transit gateway.

The following attribute is exported inside the block:

* `instance_id` - (String) The ID of the instance the target network interface is attached to.

~> **Note** The `route.instance_id` attribute cannot be specified in configuration files.
To route traffic to an instance, use `network_interface_id` with the instance's `primary_network_interface_id` attribute.

~> **Note** Changing any argument inside a `route` block makes `terraform plan` report that every route of the route table is replaced.
This is a [known Terraform SDK issue][sdk-diff-issue]: an argument omitted in the configuration is treated as `null` there, while the state keeps the default value for its type, and the difference is highlighted as soon as the block is edited.
The plan is misleading, not the apply: the provider updates only the routes that actually changed and leaves the rest alone.

## Attribute reference

### Supported attributes

~> **Note** Only the target that is entered is exported as a readable attribute once the route resource is created.
The `route.instance_id` attribute is an exception: it is exported for any route that targets a network interface attached to an instance.

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the route table.
* `id` - (String) The ID of the route table.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`owner_id`, `route.carrier_gateway_id`, `route.core_network_arn`, `route.destination_prefix_list_id`, `route.egress_only_gateway_id`, `route.ipv6_cidr_block`, `route.local_gateway_id`, `route.vpc_endpoint_id`, `route.vpc_peering_connection_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `5 minutes`) Used for route creation.
- `update` - (Default `2 minutes`) Used for updating the route.
- `delete` - (Default `5 minutes`) Used for route deletion.

## Import

Route tables can be imported using the route table `id`, for example:

```
$ terraform import aws_route_table.example rtb-12345678
```

