---
subcategory: "Route 53"
layout: "aws"
page_title: "aws_route53_zone"
description: |-
  Manages a Route53 Hosted Zone
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block

# Resource: aws_route53_zone

Manages a Route53 Hosted Zone.

## Example Usage

### Public Zone

```terraform
resource "aws_route53_zone" "primary" {
  name = "example.com"
}
```

### Public Subdomain Zone

For use in subdomains, note that you need to create a
`aws_route53_record` of type `NS` as well as the subdomain
zone.

```terraform
resource "aws_route53_zone" "main" {
  name = "example.com"
}

resource "aws_route53_zone" "dev" {
  name = "dev.example.com"

  tags = {
    Environment = "dev"
  }
}

resource "aws_route53_record" "dev-ns" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "dev.example.com"
  type    = "NS"
  ttl     = "30"
  records = aws_route53_zone.dev.name_servers
}
```

### Private Zone

~> **Note** Private zones require one VPC association at all times.

```terraform
resource "aws_route53_zone" "private" {
  name = "example.com"

  vpc {
    vpc_id = aws_vpc.example.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This is the name of the hosted zone.
* `comment` - (Optional) A comment for the hosted zone. Defaults to 'Managed by Terraform'.
* `force_destroy` - (Optional) Whether to destroy all records (possibly managed outside of Terraform) in the zone when destroying the zone.
* `tags` - (Optional) A map of tags to assign to the zone. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.
* `vpc` - (Optional) Configuration block(s) specifying VPC to associate with a private hosted zone.

### vpc Argument Reference

* `vpc_id` - (Required) ID of the VPC to associate.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the Hosted Zone.
* `zone_id` - The Hosted Zone ID. This can be referenced by zone records.
* `name_servers` - A list of name servers in associated (or default) delegation set.
  Find more about delegation sets in [AWS docs](https://docs.aws.amazon.com/Route53/latest/APIReference/actions-on-reusable-delegation-sets.html).
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`delegation_set_id`, `vpc.vpc_region`.

## Import

Route53 Zones can be imported using the `zone id`, e.g.,

```
$ terraform import aws_route53_zone.myzone z-xxxxxxxx
```
