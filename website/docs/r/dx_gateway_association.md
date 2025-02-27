---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_gateway_association"
description: |-
  Manages a Direct Connect gateway association with a transit gateway.
---

[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_dx_gateway_association

Manages a Direct Connect gateway association with a transit gateway.

## Example Usage

```terraform
resource "aws_dx_gateway" "example" {
  name            = "tf-dxassoc-example"
  amazon_side_asn = "64512"
}

resource "aws_ec2_transit_gateway" "example" {
}

resource "aws_dx_gateway_association" "example" {
  dx_gateway_id         = aws_dx_gateway.example.id
  associated_gateway_id = aws_ec2_transit_gateway.example.id

  allowed_prefixes = [
    "10.255.255.0/30",
    "10.255.255.8/30",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `allowed_prefixes` - (Required, Editable) VPC prefixes (CIDRs) to advertise to the Direct Connect gateway.
* `associated_gateway_id` - (Required) The ID of the transit gateway that the Direct Connect gateway must be associated with.
* `dx_gateway_id` - (Required) The ID of the Direct Connect gateway.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `dx_gateway_owner_account_id` - The ID of the project that owns the Direct Connect gateway.
* `id` - The ID of the Direct Connect gateway association.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`associated_gateway_type`, `associated_gateway_owner_account_id`, `proposal_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `30 minutes`) Timeout for creating the association
- `update` - (Default `30 minutes`) Timeout for updating the association
- `delete` - (Default `30 minutes`) Timeout for destroying the association

## Import

Direct Connect gateway associations can be imported using `dx_gateway_id` together with `associated_gateway_id`,
e.g.,

```
$ terraform import aws_dx_gateway_association.example dxgw-12345678/tgw-12345678
```
