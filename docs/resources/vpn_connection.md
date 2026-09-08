---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_connection"
description: |-
  Manages a Site-to-Site VPN connection. A Site-to-Site VPN connection is an Internet Protocol security (IPSec) VPN connection between a VPC and an on-premises network.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[sensitive-data]: https://www.terraform.io/docs/state/sensitive-data.html
[vpn-connections]: https://docs.k2.cloud/en/services/interconnect/vpn_connections/operations.html

# Resource: aws_vpn_connection

Manages a Site-to-Site VPN connection. A Site-to-Site VPN connection is an Internet Protocol security (IPSec) VPN connection between a VPC and an on-premises network.

For more information about VPN connections, see [user documentation][vpn-connections].

-> **Note** For convenience, the ID of the VPN gateway is the same as the ID of the VPC, to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).

~> **Note** Only one VPN connection can exist between a customer gateway and a VPC. Creating a second connection for the same pair returns the existing connection instead of a new one.

~> **Note** All arguments including `tunnel1_preshared_key` and `tunnel2_preshared_key` will be stored in the raw state as plain-text.
[Read more about sensitive data in state][sensitive-data].

## Example Usage

### Fault-tolerant connection

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.8.0/24"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_customer_gateway" "example" {
  bgp_asn    = 65000
  ip_address = "172.0.0.1"
  type       = "ipsec.1"

  tags = {
    Name = "tf-customer-gateway"
  }
}

resource "aws_vpn_connection" "example" {
  vpn_gateway_id      = aws_vpc.example.id # vpc_id can be used as vpn_gateway_id
  customer_gateway_id = aws_customer_gateway.example.id
  type                = aws_customer_gateway.example.type

  tags = {
    Name = "tf-vpn-connection"
  }
}
```

### Connection without fault tolerance

~> **Note** This example uses the VPC defined in the [fault-tolerant connection example](#fault-tolerant-connection).
Only one connection can exist between a customer gateway and a VPC, so the example creates its own customer gateway.

```terraform
resource "aws_customer_gateway" "single_tunnel" {
  bgp_asn    = 65001
  ip_address = "172.0.0.2"
  type       = "ipsec.1"

  tags = {
    Name = "tf-customer-gateway-single-tunnel"
  }
}

resource "aws_vpn_connection" "single_tunnel" {
  vpn_gateway_id      = aws_vpc.example.id # vpc_id can be used as vpn_gateway_id
  customer_gateway_id = aws_customer_gateway.single_tunnel.id
  type                = aws_customer_gateway.single_tunnel.type
  high_availability   = false

  tunnel1_inside_cidr   = "169.254.252.8/30"
  tunnel1_preshared_key = "tf_example_key"
  tunnel1_ike_versions  = ["ikev2"]

  tags = {
    Name = "tf-vpn-connection-single-tunnel"
  }
}
```

## Routing

Static routes to the networks behind the customer gateway are configured with [`aws_route`](route.md) or the `route` block of [`aws_route_table`](route_table.md), where the `gateway_id` argument is set to the ID of the VPN connection.

To install the routes received over BGP into a route table, use [`aws_vpn_gateway_route_propagation`](vpn_gateway_route_propagation.md).
A fault-tolerant connection supports BGP routing only, static routing is not available for it, see [user documentation][vpn-connections].

## Argument Reference

The following arguments are required:

* `customer_gateway_id` - (Required, Forces new resource, String) The ID of the customer gateway.
* `type` - (Required, Forces new resource, String) The type of VPN connection.
    * _Valid values:_ `ipsec.1`, `ipsec.legacy`
* `vpn_gateway_id` - (Required, Forces new resource, String) The ID of the VPN gateway.

The following arguments are optional:

* `high_availability` - (Optional, Forces new resource, Boolean) Indicates whether the connection is created with two tunnels terminated in different availability zones.
A connection without fault tolerance has a single tunnel, so the `tunnel2_*` arguments cannot be specified for it.
    * _Default value:_ `true`
* `local_ipv4_network_cidr` - (Optional, Forces new resource, String) The IPv4 CIDR on the customer gateway (on-premises) side of the VPN connection.
    * _Constraints:_ The value must not fall within the range of `169.254.0.0/16`
    * _Default value:_ `0.0.0.0/0`
* `remote_ipv4_network_cidr` - (Optional, Forces new resource, String) The IPv4 CIDR on the cloud side of the VPN connection.
    * _Constraints:_ The value must not fall within the range of `169.254.0.0/16`
    * _Default value:_ `0.0.0.0/0`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the connection. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

The following arguments configure the tunnels of the connection.
The `tunnel2_*` arguments are supported only for a fault-tolerant connection.

* `tunnel1_ike_versions`, `tunnel2_ike_versions` - (Optional, Forces new resource, Set of strings) The IKE version that is permitted for the VPN tunnel.
    * _Constraints:_ Only one version can be specified
    * _Valid values:_ `ikev1`, `ikev2`
    * _Default value:_ `ikev1`
* `tunnel1_inside_cidr`, `tunnel2_inside_cidr` - (Optional, Forces new resource, String) The CIDR block of the inside IP addresses for the VPN tunnel.
The first address of the block is assigned to the VPC, the second one is assigned to the customer gateway.
    * _Constraints:_ A `/30` CIDR block from the `169.254.252.0/22` range, unique among the connections of the VPN gateway
* `tunnel1_preshared_key`, `tunnel2_preshared_key` - (Optional, Forces new resource, String) The pre-shared key (PSK) used for the primary authentication between the VPN gateway and the customer gateway.
    * _Constraints:_ From 8 to 64 alphanumeric characters, periods (`.`) and underscores (`_`), the key cannot start with zero (`0`)
* `tunnel1_phase1_dh_group_numbers`, `tunnel2_phase1_dh_group_numbers` - (Optional, Forces new resource, Set of integers) The Diffie-Hellman group numbers that are permitted for the VPN tunnel for phase 1 IKE negotiations.
    * _Constraints:_ Up to six values
    * _Valid values:_ `2`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Default value:_ `5`, `14`, `15`, `16`, `17`, `18`
* `tunnel1_phase2_dh_group_numbers`, `tunnel2_phase2_dh_group_numbers` - (Optional, Forces new resource, Set of integers) The Diffie-Hellman group number that is permitted for the VPN tunnel for phase 2 IKE negotiations.
The `0` value disables the Perfect Forward Secrecy (PFS) mode. In order to prevent the session encryption key from being compromised, do not disable PFS.
    * _Constraints:_ Only one number can be specified, the `19`, `20` and `21` numbers are not supported for `ikev1`
    * _Valid values:_ `0`, `2`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Default value:_ `14`
* `tunnel1_phase1_encryption_algorithms`, `tunnel2_phase1_encryption_algorithms` - (Optional, Forces new resource, Set of strings) The encryption algorithms that are permitted for the VPN tunnel for phase 1 IKE negotiations.
    * _Constraints:_ The `aes_gcm128`, `aes_gcm256` and `chacha20poly1305` algorithms are not supported for `ikev1`
    * _Valid values:_ `aes128`, `aes256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Default value:_ All the supported algorithms, depending on the IKE version
* `tunnel1_phase2_encryption_algorithms`, `tunnel2_phase2_encryption_algorithms` - (Optional, Forces new resource, Set of strings) The encryption algorithms that are permitted for the VPN tunnel for phase 2 IKE negotiations.
    * _Constraints:_ The `chacha20poly1305` algorithm is not supported for `ikev1`
    * _Valid values:_ `aes128`, `aes256`, `aes_ccm128`, `aes_ccm256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Default value:_ All the supported algorithms, depending on the IKE version
* `tunnel1_phase1_integrity_algorithms`, `tunnel2_phase1_integrity_algorithms` - (Optional, Forces new resource, Set of strings) The integrity algorithms that are permitted for the VPN tunnel for phase 1 IKE negotiations.
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel1_phase2_integrity_algorithms`, `tunnel2_phase2_integrity_algorithms` - (Optional, Forces new resource, Set of strings) The integrity algorithms that are permitted for the VPN tunnel for phase 2 IKE negotiations.
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel1_phase1_lifetime_seconds`, `tunnel2_phase1_lifetime_seconds` - (Optional, Forces new resource, Integer) The lifetime for phase 1 of the IKE negotiation, in seconds.
    * _Constraints:_ From `900` to `28800`
    * _Default value:_ `28800`
* `tunnel1_phase2_lifetime_seconds`, `tunnel2_phase2_lifetime_seconds` - (Optional, Forces new resource, Integer) The lifetime for phase 2 of the IKE negotiation, in seconds.
    * _Constraints:_ From `900` to `3600`, the value must not exceed the phase 1 lifetime
    * _Default value:_ `3600`
* `tunnel1_replay_window_size`, `tunnel2_replay_window_size` - (Optional, Forces new resource, Integer) The number of packets in an IKE replay window.
    * _Constraints:_ From `32` to `2048`
    * _Default value:_ `1024`

## Attribute Reference

~> **Note** The platform does not return the tunnel options, so their values are read from the customer gateway configuration.

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the VPN connection.
* `customer_gateway_configuration` - (String) The configuration information for the VPN connection's customer gateway (in the native XML format).
* `id` - (String) The ID of the VPN connection.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `tunnel1_address`, `tunnel2_address` - (String) The public IP address of the VPN tunnel.
* `tunnel1_bgp_asn`, `tunnel2_bgp_asn` - (String) The BGP ASN of the VPN tunnel.
* `tunnel1_bgp_holdtime`, `tunnel2_bgp_holdtime` - (Integer) The BGP hold time of the VPN tunnel.
* `tunnel1_cgw_inside_address`, `tunnel2_cgw_inside_address` - (String) The RFC 6890 link-local address of the VPN tunnel (customer gateway side).
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

`core_network_arn`, `core_network_attachment_arn`, `enable_acceleration`, `local_ipv6_network_cidr`, `remote_ipv6_network_cidr`, `routes`, `static_routes_only`, `transit_gateway_attachment_id`, `transit_gateway_id`, `tunnel_inside_ip_version`, `tunnel1_dpd_timeout_action`, `tunnel2_dpd_timeout_action`, `tunnel1_dpd_timeout_seconds`, `tunnel2_dpd_timeout_seconds`, `tunnel1_inside_ipv6_cidr`, `tunnel2_inside_ipv6_cidr`, `tunnel1_rekey_fuzz_percentage`, `tunnel2_rekey_fuzz_percentage`, `tunnel1_rekey_margin_time_seconds`, `tunnel2_rekey_margin_time_seconds`, `tunnel1_startup_action`, `tunnel2_startup_action`, `vgw_telemetry.certificate_arn`.

## Import

VPN connections can be imported using the ID of VPN connection, for example:

```
$ terraform import aws_vpn_connection.example vpn-12345678
```
