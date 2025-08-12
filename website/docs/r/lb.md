---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb"
description: |-
  Manages a load balancer.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_lb

Manages a load balancer.

~> **Note** `aws_alb` is known as `aws_lb`. The functionality is identical.

## Example Usage

### Application Load Balancer

```terraform
resource "aws_lb" "test" {
  name               = "test-lb-tf"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.lb_sg.id]
  subnets            = [for subnet in aws_subnet.public : subnet.id]

  enable_deletion_protection = true

  access_logs {
    bucket  = aws_s3_bucket.lb_logs.bucket
    prefix  = "test-lb"
    enabled = true
  }

  tags = {
    Environment = "production"
  }
}
```

### Network Load Balancer

```terraform
resource "aws_lb" "test" {
  name               = "test-lb-tf"
  internal           = false
  load_balancer_type = "network"
  subnets            = [for subnet in aws_subnet.public : subnet.id]

  enable_deletion_protection = true

  tags = {
    Environment = "production"
  }
}
```

### Specifying Elastic IPs

```terraform
resource "aws_lb" "example" {
  name               = "example"
  load_balancer_type = "network"

  subnet_mapping {
    subnet_id     = aws_subnet.example1.id
    allocation_id = aws_eip.example1.id
  }

  subnet_mapping {
    subnet_id     = aws_subnet.example2.id
    allocation_id = aws_eip.example2.id
  }
}
```

### Specifying private IP addresses for an internal-facing load balancer

```terraform
resource "aws_lb" "example" {
  name               = "example"
  load_balancer_type = "network"

  subnet_mapping {
    subnet_id            = aws_subnet.example1.id
    private_ipv4_address = "10.0.1.15"
  }

  subnet_mapping {
    subnet_id            = aws_subnet.example2.id
    private_ipv4_address = "10.0.2.15"
  }
}
```

## Argument Reference

~> **Note** Please note that internal LBs can only use `ipv4` as the ip_address_type. You can only change to `dualstack` ip_address_type if the selected subnets are IPv6 enabled.

~> **Note** Please note that one of either `subnets` or `subnet_mapping` is required.

The following arguments are supported:

* `name` - (Optional) The name of the LB. This name must be unique within your AWS account, can have a maximum of 32 characters,
must contain only alphanumeric characters or hyphens, and must not begin or end with a hyphen. If not specified,
Terraform will autogenerate a name beginning with `tf-lb`.
* `name_prefix` - (Optional) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `internal` - (Optional) If true, the LB will be internal.
* `load_balancer_type` - (Optional) The type of load balancer to create. Possible values are `application`, `gateway`, or `network`. The default value is `application`.
* `security_groups` - (Optional) A list of security group IDs to assign to the LB. Only valid for load balancers of type `application`.
* `drop_invalid_header_fields` - (Optional) Indicates whether HTTP headers with header fields that are not valid are removed by the load balancer (true) or routed to targets (false). The default is false. Elastic Load Balancing requires that message header names contain only alphanumeric characters and hyphens. Only valid for load balancers of type `application`.
* `access_logs` - (Optional) An access logs block. Access logs documented below.
* `subnets` - (Optional) A list of subnet IDs to attach to the LB. Subnets
cannot be updated for load balancers of type `network`. Changing this value
for load balancers of type `network` will force a recreation of the resource.
* `subnet_mapping` - (Optional) A subnet mapping block as documented below.
* `idle_timeout` - (Optional) The time in seconds that the connection is allowed to be idle. Only valid for load balancers of type `application`. Default: 60.
* `enable_deletion_protection` - (Optional) If true, deletion of the load balancer will be disabled via
   the AWS API. This will prevent Terraform from deleting the load balancer.
    * _Default value_: `false`
* `enable_cross_zone_load_balancing` - (Optional) If true, cross-zone load balancing of the load balancer will be enabled.
   This is a `network` load balancer feature.
    * _Default value_: `false`
* `enable_http2` - (Optional) Indicates whether HTTP/2 is enabled in `application` load balancers.
    * _Default value_: `true`
* `enable_waf_fail_open` - (Optional) Indicates whether to allow a WAF-enabled load balancer to route requests to targets if it is unable to forward the request to AWS WAF.
    * _Default value_: `false`
* `customer_owned_ipv4_pool` - (Optional) The ID of the customer owned ipv4 pool to use for this load balancer.
* `ip_address_type` - (Optional) The type of IP addresses used by the subnets for your load balancer. The possible values are `ipv4` and `dualstack`
* `desync_mitigation_mode` - (Optional) Determines how the load balancer handles requests that might pose a security risk to an application due to HTTP desync.
    * _Valid values_: `defensive` (default), `monitor`, `strictest`
* `tags` - (Optional) Map of tags to assign to the load balancer. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider level.

Access logs (`access_logs`) support the following:

* `bucket` - (Required) The S3 bucket name to store the logs in.
* `prefix` - (Optional) The S3 bucket prefix. Logs are stored in the root if not configured.
* `enabled` - (Optional) Boolean to enable / disable `access_logs`.
    * _Default value_: `false`, even when `bucket` is specified

Subnet mapping (`subnet_mapping`) blocks support the following:

* `subnet_id` - (Required) The ID of the subnet to attach to the load balancer. You can specify only one subnet per availability zone.
* `allocation_id` - (Optional) The allocation ID of the Elastic IP address.
* `private_ipv4_address` - (Optional) A private ipv4 address within the subnet to assign to the internal-facing load balancer.
* `ipv6_address` - (Optional) An IPv6 address within the subnet to assign to the internet-facing load balancer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the load balancer (matches `arn`).
* `arn` - The Amazon Resource Name (ARN) of the load balancer (matches `id`).
* `arn_suffix` - The Amazon Resource Name (ARN) suffix for use with CloudWatch metrics.
* `dns_name` - The DNS name of the load balancer.
* `tags_all` - Map of tags assigned to the load balancer, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `zone_id` - The canonical hosted zone ID of the load balancer (to be used in a Route 53 Alias record).
* `subnet_mapping.*.outpost_id` - ID of the Outpost containing the load balancer.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10 minutes`) Used for creating LB.
- `update` - (Default `10 minutes`) Used for LB modifications.
- `delete` - (Default `10 minutes`) Used for destroying LB.

## Import

LBs can be imported using their ARNs, e.g.,

```
$ terraform import aws_lb.bar arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/50dc6c495c0c9188
```
