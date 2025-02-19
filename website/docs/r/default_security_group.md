---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_default_security_group"
description: |-
  Manage the default security group resource of a VPC.
---

[default-security-groups]: https://docs.k2.cloud/en/services/security/securitygroups.html#id3
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block

# Resource: aws_default_security_group

Manages the default security group of a VPC. This resource can manage the default security group of the default or a non-default VPC.

~> **Note** This is an advanced resource with special caveats. Please read this document in its entirety before using this resource. The `aws_default_security_group` resource behaves differently from normal resources. Terraform does not _create_ this resource but instead attempts to "adopt" it into management.

When Terraform first adopts the default security group, it **immediately removes all ingress and egress rules in the security group**. It then creates any rules specified in the configuration. This way only the rules specified in the configuration are created.

This resource treats its inline rules as absolute; only the rules defined inline are created, and any additions/removals external to this resource will result in diff shown. For these reasons, this resource is incompatible with the [`aws_security_group_rule`](security_group_rule.md) resource.

For more information about default security groups, see the documentation on [default security groups][default-security-groups]. To manage normal security groups, see the [`aws_security_group`][tf-security-group] resource.

## Example Usage

The following config gives the default security group the same rules that the cloud provides by default but under management by Terraform. This means that any ingress or egress rules added or changed will be detected as drift.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_security_group" "example" {
  vpc_id = aws_vpc.mainvpc.id

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

### Example Config To Deny All Egress Traffic, Allowing Ingress

The following denies all egress traffic by omitting any `egress` rules, while including the default `ingress` rule to allow all traffic.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_security_group" "example" {
  vpc_id = aws_vpc.mainvpc.id

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }
}
```

### Removing `aws_default_security_group` From Your Configuration

Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the security group.
All ingress or egress rules will be left as they are at the time of removal. You can resume managing them via the cloud console.

## Argument Reference

The following arguments are optional:

* `egress` - (Optional, VPC only) Configuration block. Detailed below.
* `ingress` - (Optional) Configuration block. Detailed below.
* `tags` - (Optional) Map of tags to assign to the resource. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.
* `vpc_id` - (Optional, Forces new resource) VPC ID. **Note that changing the `vpc_id` will _not_ restore any default security group rules that were modified, added, or removed.** It will be left in its current state.

### egress and ingress

Both arguments are processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).

Both `egress` and `ingress` objects have the same arguments.

* `cidr_blocks` - (Optional) List of CIDR blocks.
* `description` - (Optional) Description of this rule.
* `from_port` - (Required) Start port (or ICMP type number if protocol is `icmp`)
* `ipv6_cidr_blocks` - (Optional) List of IPv6 CIDR blocks.
* `protocol` - (Required) Protocol. If you select a protocol of "-1" (semantically equivalent to `all`, which is not a valid value here), you must specify a `from_port` and `to_port` equal to `0`. If not `icmp`, `tcp`, `udp`, or `-1` use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).
* `security_groups` - (Optional) List of security group names or IDs.
* `self` - (Optional) Whether the security group itself will be added as a source to this egress rule.
* `to_port` - (Required) End range port (or ICMP code if protocol is `icmp`).

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the security group.
* `description` - Description of the security group.
* `id` - ID of the security group.
* `name` - Name of the security group.
* `owner_id` - The project ID.
* `tags_all` -Map of tags to assign to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `prefix_list_ids`.

## Import

Security groups can be imported using the `security group id`, e.g.,

```
$ terraform import aws_default_security_group.default_sg sg-12345678
```
