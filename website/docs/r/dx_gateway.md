---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_gateway"
description: |-
  Manages a Direct Connect gateway.
---

[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_dx_gateway

Manages a Direct Connect gateway.

## Example Usage

```terraform
resource "aws_dx_gateway" "example" {
  name            = "tf-dxgw-example"
  amazon_side_asn = "64512"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Direct Connect gateway.
* `amazon_side_asn` - (Required) The ASN to be configured on the cloud side of the connection. The ASN must be in the private range of 64,512 to 65,534 or 4,200,000,000 to 4,294,967,294.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Direct Connect gateway.
* `owner_account_id` - The ID of the project that owns the Direct Connect gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10 minutes`) Used for creating the gateway
- `delete` - (Default `10 minutes`) Used for destroying the gateway

## Import

Direct Connect gateways can be imported using `id`, e.g.,

```
$ terraform import aws_dx_gateway.example dxgw-12345678
```
