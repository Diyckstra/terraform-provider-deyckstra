---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_default_network_acl"
description: |-
  Manages the default network ACL of a VPC.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[network-acl]: https://docs.k2.cloud/en/services/security/networkacl.html

# Resource: aws_default_network_acl

Manages the default network ACL of a VPC. This resource can manage the default network ACL of the default or a non-default VPC.

~> **Note** This is an advanced resource with special caveats. Please read this document in its entirety before using this resource. The `aws_default_network_acl` behaves differently from normal resources. Terraform does not _create_ this resource but instead attempts to "adopt" it into management.

Every VPC has a default network ACL that can be managed but not destroyed. When Terraform first adopts the default network ACL, it **immediately removes all rules in the ACL**. It then proceeds to create any rules specified in the configuration. This step is required so that only the rules specified in the configuration are created.

This resource treats its inline rules as absolute; only the rules defined inline are created, and any additions/removals external to this resource will result in diffs being shown. For these reasons, this resource is incompatible with the `aws_network_acl_rule` resource.

For more information about network ACLs, see the documentation on [network ACL][network-acl].

## Example Usage

### Basic Example

The following config gives the default network ACL the same rules that the cloud includes but pulls the resource under management by Terraform. This means that any ACL rules added or changed will be detected as drift.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_network_acl" "default" {
  default_network_acl_id = aws_vpc.mainvpc.default_network_acl_id

  ingress {
    protocol   = -1
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }

  egress {
    protocol   = -1
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }
}
```

### Example: Deny All Egress Traffic, Allow Ingress

The following denies all egress traffic by omitting any `egress` rules, while including the default `ingress` rule to allow all traffic.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_network_acl" "default" {
  default_network_acl_id = aws_vpc.mainvpc.default_network_acl_id

  ingress {
    protocol   = -1
    rule_no    = 100
    action     = "allow"
    cidr_block = aws_vpc.mainvpc.cidr_block
    from_port  = 0
    to_port    = 0
  }
}
```

### Example: Deny All Traffic To Any Subnet In The Default Network ACL

This config denies all traffic in the default ACL. This can be useful if you want to lock down the VPC to force all resources to assign a non-default ACL.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_network_acl" "default" {
  default_network_acl_id = aws_vpc.mainvpc.default_network_acl_id

  # no rules defined, deny all traffic in this ACL
}
```

### Managing Subnets In A Default Network ACL

Within a VPC, all subnets must be associated with a network ACL. In order to "delete" the association between a subnet and a non-default network ACL, the association is destroyed by replacing it with an association between the subnet and the default ACL instead.

When managing the default network ACL, you cannot "remove" subnets. Instead, they must be reassigned to another network ACL, or the subnet itself must be destroyed. Because of these requirements, removing the `subnet_ids` attribute from the configuration of a `aws_default_network_acl` resource may result in a reoccurring plan, until the subnets are reassigned to another network ACL or are destroyed.

Because subnets are by default associated with the default network ACL, any non-explicit association will show up as a plan to remove the subnet. For example: if you have a custom `aws_network_acl` with two subnets attached, and you remove the `aws_network_acl` resource, after successfully destroying this resource future plans will show a diff on the managed `aws_default_network_acl`, as those two subnets have been orphaned by the now destroyed network acl and thus adopted by the default network ACL. In order to avoid a reoccurring plan, they will need to be reassigned, destroyed, or added to the `subnet_ids` attribute of the `aws_default_network_acl` entry.

As an alternative to the above, you can also specify the following lifecycle configuration in your `aws_default_network_acl` resource:

```terraform
resource "aws_default_network_acl" "default" {
  # ... other configuration ...

  lifecycle {
    ignore_changes = [subnet_ids]
  }
}
```

### Removing `aws_default_network_acl` From Your Configuration

Each VPC comes with a default network ACL that cannot be deleted. The `aws_default_network_acl` allows you to manage this network ACL, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your statefile and management, **but will not destroy the network ACL.** All subnets associations and ingress or egress rules will be left as they are at the time of removal. You can resume managing them via the cloud console.

## Argument Reference

The following arguments are required:

* `default_network_acl_id` - (Required) network ACL ID to manage.
  This attribute is exported from `aws_vpc`, or manually found via the cloud console.

The following arguments are optional:

* `egress` - (Optional) Configuration block for an egress rule. Detailed below.
* `ingress` - (Optional) Configuration block for an ingress rule. Detailed below.
* `subnet_ids` - (Optional) List of subnet IDs to apply the ACL to. See the notes below on managing subnets in the default network ACL.
* `tags` - (Optional) Map of tags to assign to the resource. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

### egress and ingress

Both the `egress` and `ingress` configuration blocks have the same arguments.

The following arguments are required:

* `action` - (Required) The action to take.
* `from_port` - (Required) The from port to match.
* `protocol` - (Required) The protocol to match. If using the -1 'all' protocol, you must specify a from and to port of 0.
* `rule_no` - (Required) The rule number. Used for ordering.
* `to_port` - (Required) The to port to match.

The following arguments are optional:

* `cidr_block` - (Optional) The CIDR block to match. This must be a valid network mask.
* `icmp_code` - (Optional) The ICMP type code to be used. Default 0.
* `icmp_type` - (Optional) The ICMP type to be used. Default 0.

-> For more information on ICMP types and codes, see [Internet Control Message Protocol (ICMP) Parameters](https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml).

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the default network ACL
* `id` - ID of the default network ACL
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `vpc_id` -  ID of the associated VPC

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`ipv6_cidr_block`, `owner_id`.

## Import

Default network ACLs can be imported using `id`, e.g.,

```
$ terraform import aws_default_network_acl.sample acl-12345678
```

