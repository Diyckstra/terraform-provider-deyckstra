---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_group_rule"
description: |-
  Manages a security group rule.
---

# Resource: aws_security_group_rule

Manages a security group rule.
Represents a single `ingress` or `egress` group rule, which can be added to external security groups.

~> **Note on security groups and security group rules:** Terraform currently
provides both a standalone security group rule resource (a single `ingress` or
`egress` rule), and a [`aws_security_group`](security_group.md) with `ingress` and `egress` rules
defined in-line. At this time you cannot use a security group with in-line rules
in conjunction with any security group rule resources. Doing so will cause
a conflict of rule settings and will overwrite rules.

~> **Note** Setting `protocol = "all"` or `protocol = -1` with `from_port` and `to_port` will result in the EC2 API creating a security group rule with all ports open. This API behavior cannot be controlled by Terraform and may generate warnings in the future.

## Example Usage

Basic usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_security_group" "example" {
  name        = "test_security_group"
  description = "test_security_group"
  vpc_id      = aws_vpc.example.id
}

resource "aws_security_group_rule" "example" {
  type              = "ingress"
  from_port         = 0
  to_port           = 65535
  protocol          = "tcp"
  cidr_blocks       = [aws_vpc.example.cidr_block]
  security_group_id = aws_security_group.example.id
}
```

## Argument Reference

The following arguments are required:

* `from_port` - (Required) Start port (or ICMP type number if protocol is "icmp" or "icmpv6").
* `protocol` - (Required) Protocol. If not icmp, icmpv6, tcp, udp, or all use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)
* `security_group_id` - (Required) Security group to apply this rule to.
* `to_port` - (Required) End port (or ICMP code if protocol is "icmp").
* `type` - (Required) Type of rule being created. Valid options are `ingress` (inbound)
or `egress` (outbound).

The following arguments are optional:

* `cidr_blocks` - (Optional) List of CIDR blocks. Cannot be specified with `source_security_group_id` or `self`.
* `description` - (Optional) Description of the rule.
* `ipv6_cidr_blocks` - (Optional) List of IPv6 CIDR blocks. Cannot be specified with `source_security_group_id` or `self`.
* `self` - (Optional) Whether the security group itself will be added as a source to this ingress rule. Cannot be specified with `cidr_blocks`, `ipv6_cidr_blocks`, or `source_security_group_id`.
* `source_security_group_id` - (Optional) ID of the security group to allow access to/from, depending on the `type`. Cannot be specified with `cidr_blocks`, `ipv6_cidr_blocks`, or `self`.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the security group rule.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `prefix_list_ids`.

## Import

-> **Unsupported operation**
Import security group rules is currently unsupported.
