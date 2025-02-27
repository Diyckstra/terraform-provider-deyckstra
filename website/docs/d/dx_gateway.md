---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_gateway"
description: |-
  Provides information about a Direct Connect gateway.
---

# Data Source: aws_dx_gateway

Provides information about a Direct Connect gateway.

## Example Usage

```terraform
data "aws_dx_gateway" "selected" {
  name = "tf-dxgw-example"
}
```

## Argument Reference

* `name` - (Required) The name of the Direct Connect gateway.

## Attribute Reference

* `amazon_side_asn` - The ASN for the cloud side of the connection.
* `id` - The ID of the Direct Connect gateway.
* `owner_account_id` - The ID of the project that owns the Direct Connect gateway.
