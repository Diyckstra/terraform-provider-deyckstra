---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_default_route_table"
description: |-
  Manages the default route table of a VPC.
---

[attribute-as-blocks]: https://www.terraform.io/docs/configuration/attr-as-blocks.html
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[route-tables]: https://docs.k2.cloud/en/services/networking/routetables.html
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_default_route_table

Manages the default route table of a VPC. This resource can manage the default route table of the default or a non-default VPC.

~> **Note** This is an advanced resource with special caveats. Please read this document in its entirety before using this resource. The `aws_default_route_table` resource behaves differently from normal resources. Terraform does not _create_ this resource but instead attempts to "adopt" it into management. **Do not** use both `aws_default_route_table` to manage a default route table **and** `aws_main_route_table_association` with the same VPC due to possible route conflicts. See [aws_main_route_table_association](main_route_table_association.md) documentation for more details.

Every VPC has a default route table that can be managed but not destroyed. When Terraform first adopts a default route table, it **immediately removes all defined routes**. It then proceeds to create any routes specified in the configuration. This step is required so that only the routes specified in the configuration exist in the default route table.

For more information, see the documentation on [route tables][route-tables]. For information about managing normal route tables in Terraform, see [`aws_route_table`](route_table.md).

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

resource "aws_default_route_table" "example" {
  default_route_table_id = aws_vpc.example.default_route_table_id

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
resource "aws_default_route_table" "example" {
  default_route_table_id = aws_vpc.example.default_route_table_id

  route = []

  tags = {
    Name = "example"
  }
}
```

## Argument Reference

The following arguments are required:

* `default_route_table_id` - (Required) ID of the default route table.

The following arguments are optional:

* `propagating_vgws` - (Optional) List of virtual gateways for propagation.
* `route` - (Optional) Configuration block of routes. Detailed below. This argument is processed in [attribute-as-blocks mode][attribute-as-blocks]. This means that omitting this argument is interpreted as ignoring any existing routes. To remove all managed routes an empty list should be specified. See the example above.
* `tags` - (Optional) Map of tags to assign to the route table. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

### route

This argument is processed in [attribute-as-blocks mode][attribute-as-blocks].

One of the following destination arguments must be supplied:

* `cidr_block` - (Required) The CIDR block of the route.

One of the following target arguments must be supplied:

* `gateway_id` - (Optional) The ID of an internet gateway.
* `instance_id` - (Optional) The ID of an EC2 instance.
* `network_interface_id` - (Optional) The ID of an EC2 network interface.
* `transit_gateway_id` - (Optional) The ID of the transit gateway.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the route table.
* `arn` - The Amazon Resource Name (ARN) of the route table.
* `tags_all` - Map of tags assigned to the route table, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `vpc_id` - The ID of the VPC.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`destination_prefix_list_id`, `ipv6_cidr_block`, `owner_id`, `route.core_network_arn`, `route.egress_only_gateway_id`, `route.nat_gateway_id`, `route.vpc_endpoint_id`, `route.vpc_peering_connection_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `2 minutes`) Used for route creation.
- `update` - (Default `2 minutes`) Used for route creation.

## Import

Default VPC route tables can be imported using the `vpc_id`, e.g.,

```
$ terraform import aws_default_route_table.example vpc-12345678
```

