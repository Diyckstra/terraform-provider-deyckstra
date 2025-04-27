---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb"
description: |-
  Provides information about a load balancer.
---

# Data Source: aws_lb

Provides information about a load balancer.

## Example Usage

```terraform
data "aws_lb" "selected" {
  name = "lb-name"
}
```

## Argument Reference

The following arguments are supported:

* `arn` - (Optional) The Amazon Resource Name (ARN) of the load balancer.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:loadbalancer/<app|net>/lb-12345678`
* `name` - (Optional) The name of the load balancer.
* `tags` - (Optional) Map of tags. All specified tags must match tags on the desired load balancer.

~> **Note** When both `arn` and `name` are specified, `arn` takes precedence.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `dns_name` - The DNS name of the load balancer.
* `id` - The ARN of the load balancer.
* `internal` - Indicates whether the load balancer is internal or internet-facing.
* `load_balancer_type` - The type of the load balancer.
* `subnet_mapping` - List of subnet-ID-to-IP-address mappings.
  The structure of this block is [described below](#subnet_mapping).
* `subnets` - List of subnet IDs.
* `vpc_id` - The ID of the VPC.
* `zone_id` - The ID of the Route53 hosted zone associated with the load balancer.

#### subnet_mapping

The `subnet_mapping` block has the following structure:

* `subnet_id` - The ID of the subnet.
* `allocation_id` - The ID of the Elastic IP address allocation.
* `private_ipv4_address` - The private IP address.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`access_logs`, `arn_suffix`, `customer_owned_ipv4_pool`, `desync_mitigation_mode`, `drop_invalid_header_fields`, `enable_cross_zone_load_balancing`, `enable_deletion_protection`, `enable_http2`, `enable_waf_fail_open`, `idle_timeout`, `ip_address_type`, `security_groups`, `subnet_mapping.ipv6_address`, `subnet_mapping.outpost_id`.
