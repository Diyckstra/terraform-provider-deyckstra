---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_group"
description: |-
  Manages a security group.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[security-groups]: https://docs.k2.cloud/en/services/security/securitygroups.html
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_security_group

Manages a security group.

~> **Note on Security Groups and Security Group Rules:** Terraform currently
provides both a standalone [`aws_security_group_rule`](security_group_rule.md) (a single `ingress` or
`egress` rule), and a security group resource with `ingress` and `egress` rules
defined in-line. At this time you cannot use a security group with in-line rules
in conjunction with any security group rule resources. Doing so will cause
a conflict of rule settings and will overwrite rules.

For more information, see the documentation on [Security groups][security-groups].

## Example Usage

### Basic Usage

```terraform
resource "aws_vpc" "main" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_security_group" "allow_tls" {
  name        = "allow_tls"
  description = "Allow TLS inbound traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "TLS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "allow_tls"
  }
}
```

~> **Note on Egress rules:** By default, the cloud creates an `ALLOW ALL` egress rule when creating a new security group inside a VPC. When creating a new Security Group inside a VPC, **Terraform will remove this default rule**, and require you specifically re-create it if you desire that rule. We feel this leads to fewer surprises in terms of controlling your egress rules. If you desire this rule to be in place, you can use this `egress` block:

```terraform
resource "aws_security_group" "example" {
  # ... other configuration ...

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

### Change of name or name-prefix value

Security Group's Name cannot be edited after the resource is created. In fact, the `name` and `name-prefix` arguments force the creation of a new security group resource when they change value. In that case, Terraform first deletes the existing Security Group resource and then it creates a new one. If the existing Security Group is associated to a Network Interface resource, the deletion cannot complete. The reason is that Network Interface resources cannot be left with no Security Group attached and the new one is not yet available at that point.

It is required to invert the default behavior of Terraform. That is, first the new security group resource must be created, then associated to possible network interface resources and finally the old Security Group can be detached and deleted. To force this behavior, you must set the [create_before_destroy](https://www.terraform.io/language/meta-arguments/lifecycle#create_before_destroy) property:

```terraform
resource "aws_security_group" "sg_with_changeable_name" {
  name = "changeable-name"
  # ... other configuration ...

  lifecycle {
    # Necessary if changing 'name' or 'name_prefix' properties.
    create_before_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, Forces new resource) Security group description. Defaults to `Managed by Terraform`.
* `egress` - (Optional, VPC only) Configuration block for egress rules. Can be specified multiple times for each egress rule. Each egress block supports fields documented below. This argument is processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).
* `ingress` - (Optional) Configuration block for ingress rules. Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below. This argument is processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `name` - (Optional, Forces new resource) Name of the security group. If omitted, Terraform will assign a random, unique name.
* `revoke_rules_on_delete` - (Optional) Instruct Terraform to revoke all the security groups attached ingress and egress rules before deleting the rule itself. Default `false`.
* `tags` - (Optional) Map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.
* `vpc_id` - (Optional, Forces new resource) VPC ID.

### ingress

This argument is processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).

The following arguments are required:

* `from_port` - (Required) Start port (or ICMP type number if protocol is `icmp` or `icmpv6`).
* `to_port` - (Required) End range port (or ICMP code if protocol is `icmp`).
* `protocol` - (Required) Protocol. If you select a protocol of "-1" (semantically equivalent to `all`, which is not a valid value here), you must specify a `from_port` and `to_port` equal to `0`. If not `icmp`, `tcp`, `udp`, or `-1` use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).

The following arguments are optional:

* `cidr_blocks` - (Optional) List of CIDR blocks.
* `description` - (Optional) Description of this ingress rule.
* `ipv6_cidr_blocks` - (Optional) List of IPv6 CIDR blocks.
* `security_groups` - (Optional) List of security group names or group IDs.
* `self` - (Optional) Whether the security group itself will be added as a source to this ingress rule.

### egress

This argument is processed in [attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html).

The following arguments are required:

* `from_port` - (Required) Start port (or ICMP type number if protocol is `icmp`)
* `to_port` - (Required) End range port (or ICMP code if protocol is `icmp`).

The following arguments are optional:

* `cidr_blocks` - (Optional) List of CIDR blocks.
* `description` - (Optional) Description of this egress rule.
* `ipv6_cidr_blocks` - (Optional) List of IPv6 CIDR blocks.
* `protocol` - (Required) Protocol. If you select a protocol of "-1" (semantically equivalent to `all`, which is not a valid value here), you must specify a `from_port` and `to_port` equal to `0`. If not `icmp`, `tcp`, `udp`, or `-1` use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).
* `security_groups` - (Optional) List of security group names or group IDs.
* `self` - (Optional) Whether the security group itself will be added as a source to this egress rule.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the security group.
* `id` - ID of the security group.
* `owner_id` - The project ID.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported: `prefix_list_ids`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10m`) How long to wait for a security group to be created.
- `delete` - (Default `15m`) How long to wait for a security group to be deleted.

## Import

Security Groups can be imported using the `security group id`, e.g.,

```
$ terraform import aws_security_group.elb_sg sg-12345678
```


