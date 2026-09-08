---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_connection"
description: |-
  Provides details about a specific Site-to-Site VPN connection.
---

[sensitive-data]: https://www.terraform.io/docs/state/sensitive-data.html
[vpn-connections]: https://docs.k2.cloud/en/services/interconnect/vpn_connections/operations.html

# Data Source: aws_vpn_connection

Provides details about a specific Site-to-Site VPN connection.

For more information about VPN connections, see [user documentation][vpn-connections].

~> **Note** The `customer_gateway_configuration`, `tunnel1_preshared_key` and `tunnel2_preshared_key` attributes will be stored in the raw state as plain-text.
[Read more about sensitive data in state][sensitive-data].

## Example usage

### Search by ID

```terraform
data "aws_vpn_connection" "selected" {
  id = "vpn-12345678"
}
```

### Search by the customer gateway

```terraform
data "aws_vpn_connection" "selected" {
  customer_gateway_id = "cgw-12345678"
}
```

### Search by filters

```terraform
data "aws_vpn_connection" "selected" {
  filter {
    name   = "tag:Name"
    values = ["tf-vpn-connection"]
  }
}
```

## Argument reference

The arguments of this data source act as filters for querying the VPN connections.
The given filters must match exactly one connection.

* `customer_gateway_id` - (Optional, String) The ID of the customer gateway of the connection.
* `filter` - (Optional) One or more name/value pairs to use as filters, see [below](#filter).
* `id` - (Optional, String) The ID of the VPN connection.
* `type` - (Optional, String) The type of the VPN connection.
    * _Valid values:_ `ipsec.1`, `ipsec.legacy`
* `vpn_gateway_id` - (Optional, String) The ID of the VPN gateway of the connection.

### filter

The following arguments are required:

* `name` - (Required, String) The name of the filter.
    * _Valid values:_ `customer-gateway-id`, `state`, `tag-key`, `tag:<tag-key>`, `type`, `vpn-connection-id`, `vpn-gateway-id`
* `values` - (Required, List of strings) The values of the filter.

## Attribute reference

~> **Note** The platform does not return the tunnel options, so their values are read from the customer gateway configuration.
It is available only if the connection is in the `pending` or `available` state.

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the VPN connection.
* `customer_gateway_configuration` - (String) The configuration information for the VPN connection's customer gateway (in the native XML format).
* `high_availability` - (Boolean) Indicates whether the connection has two tunnels terminated in different availability zones.
* `local_ipv4_network_cidr` - (String) The IPv4 CIDR on the customer gateway (on-premises) side of the VPN connection.
* `remote_ipv4_network_cidr` - (String) The IPv4 CIDR on the cloud side of the VPN connection.
* `state` - (String) The state of the VPN connection.
* `tags` - (Map of strings) Key-value pairs assigned to the connection.
* `tunnel1_address`, `tunnel2_address` - (String) The public IP address of the VPN tunnel.
* `tunnel1_bgp_asn`, `tunnel2_bgp_asn` - (String) The BGP ASN of the VPN tunnel.
* `tunnel1_bgp_holdtime`, `tunnel2_bgp_holdtime` - (Integer) The BGP hold time of the VPN tunnel.
* `tunnel1_cgw_inside_address`, `tunnel2_cgw_inside_address` - (String) The RFC 6890 link-local address of the VPN tunnel (customer gateway side).
* `tunnel1_inside_cidr`, `tunnel2_inside_cidr` - (String) The CIDR block of the inside IP addresses for the VPN tunnel.
* `tunnel1_preshared_key`, `tunnel2_preshared_key` - (String) The pre-shared key (PSK) of the VPN tunnel.
* `tunnel1_vgw_inside_address`, `tunnel2_vgw_inside_address` - (String) The RFC 6890 link-local address of the VPN tunnel (VPN gateway side).
* `vgw_telemetry` - (Set of objects) Telemetry for the VPN tunnels, see [below](#vgw_telemetry).

The `tunnel2_*` attributes are empty for a connection without fault tolerance.

#### vgw_telemetry

The following attributes are exported inside the block:

* `accepted_route_count` - (Integer) The number of accepted routes.
* `last_status_change` - (String) The date and time of the last change in status.
* `outside_ip_address` - (String) The internet-routable IP address of the VPN gateway's outside interface.
* `status` - (String) The status of the VPN tunnel.
* `status_message` - (String) If an error occurs, a description of the error.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`core_network_arn`, `core_network_attachment_arn`, `enable_acceleration`, `local_ipv6_network_cidr`, `remote_ipv6_network_cidr`, `routes`, `static_routes_only`, `transit_gateway_attachment_id`, `transit_gateway_id`, `tunnel_inside_ip_version`, `vgw_telemetry.certificate_arn`.
