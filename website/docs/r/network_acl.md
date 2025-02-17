---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_acl"
description: |-
  Creates a network ACL.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[network-acl]: https://docs.k2.cloud/en/services/security/networkacl.html

# Resource: aws_network_acl

Creates a network ACL. You might set up network ACLs with rules similar
to your security groups in order to add an additional layer of security to your VPC.

~> **Note on Network ACLs and Network ACL Rules:** Terraform currently
provides both a standalone [`aws_network_acl_rule`](network_acl_rule.md) resource and a Network ACL resource with rules
defined in-line. At this time you cannot use a Network ACL with in-line rules
in conjunction with any Network ACL Rule resources. Doing so will cause
a conflict of rule settings and will overwrite rules.

~> **Note on Network ACLs and Network ACL Associations:** Terraform provides both a standalone [aws_network_acl_association](network_acl_association.md)
resource and a network ACL resource with a `subnet_ids` attribute. Do not use the same subnet ID in both a network ACL
resource and a network ACL association resource. Doing so will cause a conflict of associations and will overwrite the association.

For more information about network ACLs, see the documentation on [Network ACL][network-acl].

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_network_acl" "example" {
  vpc_id = aws_vpc.example.id

  egress {
    protocol   = "tcp"
    rule_no    = 200
    action     = "allow"
    cidr_block = "10.3.0.0/18"
    from_port  = 443
    to_port    = 443
  }

  ingress {
    protocol   = "tcp"
    rule_no    = 100
    action     = "allow"
    cidr_block = "10.3.0.0/18"
    from_port  = 80
    to_port    = 80
  }

  tags = {
    Name = "main"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the associated VPC.
* `subnet_ids` - (Optional) A list of subnet IDs to apply the ACL to.
* `ingress` - (Optional) Specifies an ingress rule. Parameters defined below.
  This argument is processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).
* `egress` - (Optional) Specifies an egress rule. Parameters defined below.
  This argument is processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.

### egress and ingress

Both arguments are processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).

Both `egress` and `ingress` support the following keys:

* `from_port` - (Required) The from port to match.
* `to_port` - (Required) The to port to match.
* `rule_no` - (Required) The rule number. Used for ordering.
* `action` - (Required) The action to take.
* `protocol` - (Required) The protocol to match. If using the -1 'all'
protocol, you must specify a from and to port of 0.
* `cidr_block` - (Optional) The CIDR block to match. This must be a
valid network mask.
* `icmp_type` - (Optional) The ICMP type to be used. Default 0.
* `icmp_code` - (Optional) The ICMP type code to be used. Default 0.

~> **Note** For more information on ICMP types and codes, see here: https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the network ACL.
* `arn` - ARN of the network ACL.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`ipv6_cidr_block`, `owner_id`.

## Import

Network ACLs can be imported using the `id`, e.g.,

```
$ terraform import aws_network_acl.main acl-12345678
```

