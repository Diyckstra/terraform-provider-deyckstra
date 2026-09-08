---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_gateway"
description: |-
  Provides information about a VPN gateway.
---

[describe-vpn-gateways]: https://docs.k2.cloud/en/api/ec2/actions/vpn_gateways/DescribeVpnGateways.html

# Data Source: aws_vpn_gateway

Provides information about a VPN gateway.

-> **Note** For convenience, the ID of the VPN gateway is the same as the ID of the VPC, to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

data "aws_vpn_gateway" "selected" {
  id = aws_vpc.example.id # vpc_id can be used as vpn_gateway_id
}

output "vpn_gateway_state" {
  value = data.aws_vpn_gateway.selected.state
}
```

## Argument Reference

~> **Note** The platform supports the search by the gateway ID only, see [EC2 API documentation][describe-vpn-gateways].
The other arguments are not taken into account, so the search without the `id` argument fails with the "multiple EC2 VPN Gateways matched" error whenever the project has more than one VPC.

The following argument is supported:

* `id` - (Optional, String) The ID of the VPN gateway to retrieve.

### Unsupported arguments

The following arguments are not currently supported:

`attached_vpc_id`, `availability_zone`, `filter`, `state`, `tags`.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `attached_vpc_id` - (String) The ID of the VPC the gateway is attached to.
* `state` - (String) The state of the VPN gateway.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`amazon_side_asn`, `arn`, `availability_zone`, `tags`.
