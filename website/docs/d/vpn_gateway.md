---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_gateway"
description: |-
    Provides information about a specific VPN gateway.
---

# Data Source: aws_vpn_gateway

Provides information about a specific VPN gateway.

-> **Note** For convenience, the ID of the VPN gateway is the same as the ID of the VPC, to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).

## Example Usage

```terraform
data "aws_vpn_gateway" "selected" {
  filter {
    name   = "tag:Name"
    values = ["vpn-gw"]
  }
}

output "vpn_gateway_id" {
  value = data.aws_vpn_gateway.selected.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available VPN gateways.
The given filters must match exactly one VPN gateway whose data will be exported as attributes.

* `id` - (Optional) ID of the specific VPN gateway to retrieve.
* `state` - (Optional) The state of the specific VPN gateway to retrieve.
* `availability_zone` - (Optional) The availability zone of the specific VPN gateway to retrieve.
* `attached_vpc_id` - (Optional) ID of a VPC attached to the specific VPN gateway to retrieve.
* `filter` - (Optional) One or more name/value pairs to use as filters.
  A VPN gateway will be selected if any one of the given values matches.
	Valid names and values can be found in the [EC2 API documentation][describe-vpn-gateways].
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired VPN gateway.

## Attributes Reference

All the argument attributes are also exported as result attributes.

[describe-vpn-gateways]: https://docs.k2.cloud/en/api/ec2/vpn_gateways/DescribeVpnGateways.html
