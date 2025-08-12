---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_acl_rule"
description: |-
  Creates a network ACL rule.
---

# Resource: aws_network_acl_rule

Creates an entry (a rule) in a network ACL with the specified rule number.

~> **Note on Network ACLs and Network ACL Rules:** Terraform currently
provides both a standalone network ACL rule resource and an [aws_network_acl](network-acl.md) resource with rules
defined in-line. At this time you cannot use a network ACL with in-line rules
in conjunction with any network ACL rule resources. Doing so will cause
a conflict of rule settings and will overwrite rules.

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_network_acl" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_network_acl_rule" "example" {
  network_acl_id = aws_network_acl.example.id
  rule_number    = 200
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = aws_vpc.example.cidr_block
  from_port      = 22
  to_port        = 22
}
```

~> **Note** One of either `cidr_block` or `ipv6_cidr_block` is required.

## Argument Reference

The following arguments are supported:

* `network_acl_id` - (Required) ID of the network ACL.
* `rule_number` - (Required) The rule number for the entry (for example, 100). ACL entries are processed in ascending order by rule number.
* `egress` - (Optional, bool) Indicates whether this is an egress rule (rule is applied to traffic leaving the subnet).
    * _Default value:_ `false`
* `protocol` - (Required) The protocol. A value of -1 means all protocols.
* `rule_action` - (Required) Indicates whether to allow or deny the traffic that matches the rule.
    * _Valid values:_ `allow` or `deny`
* `cidr_block` - (Optional) The network range to allow or deny, in CIDR notation (for example 172.16.0.0/24).
* `ipv6_cidr_block` - (Optional) The IPv6 CIDR block to allow or deny.
* `from_port` - (Optional) The from port to match.
* `to_port` - (Optional) The to port to match.
* `icmp_type` - (Optional) ICMP protocol: The ICMP type, e.g., -1
    * _Constraints_: Required if specifying ICMP for the protocol
* `icmp_code` - (Optional) ICMP protocol: The ICMP code, e.g., -1
    * _Constraints_: Required if specifying ICMP for the protocol

~> **Note** If the value of `protocol` is `-1` or `all`, the `from_port` and `to_port` values will be ignored and the rule will apply to all ports.

~> **Note** If the value of `icmp_type` is `-1` (which results in a wildcard ICMP type), the `icmp_code` must also be set to `-1` (wildcard ICMP code).

~> **Note** For more information on ICMP types and codes, see here: https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the network ACL rule

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `ipv6_cidr_block`.

## Import

Individual rules can be imported using `NETWORK_ACL_ID:RULE_NUMBER:PROTOCOL:EGRESS`, where `PROTOCOL` can be a decimal (e.g., 6) or string (e.g., tcp) value.
If importing a rule previously provisioned by Terraform, the `PROTOCOL` must be the input value used at creation time.
For more information on protocol numbers and keywords, see here: https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml

For example, import a network ACL rule with an argument like this:

```console
$ terraform import aws_network_acl_rule.my_rule acl-12345678:100:tcp:false
```

Or by the protocol's decimal value:

```console
$ terraform import aws_network_acl_rule.my_rule acl-12345678:100:6:false
```
