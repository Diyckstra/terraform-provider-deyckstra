---
subcategory: "Route 53"
layout: "aws"
page_title: "aws_route53_zone"
description: |-
    Provides information about a specific Route 53 hosted zone
---

# Data Source: aws_route53_zone

`Provides information about a specific Route 53 hosted zone.

This data source allows to find a hosted zone ID given hosted zone name and certain search criteria.

## Example Usage

The following example shows how to get a hosted zone from its name and from this data how to create a record set.


```terraform
data "aws_route53_zone" "selected" {
  name         = "test.com."
  private_zone = true
}

resource "aws_route53_record" "www" {
  zone_id = data.aws_route53_zone.selected.zone_id
  name    = "www.${data.aws_route53_zone.selected.name}"
  type    = "A"
  ttl     = "300"
  records = ["10.0.0.1"]
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
hosted zone. You have to use `zone_id` or `name`, not both of them. The given filter must match exactly one
hosted zone. If you use `name` field for private hosted zone, you need to add `private_zone` field to `true`

* `zone_id` - (Optional) The hosted zone id of the desired hosted zone.

* `name` - (Optional) The hosted zone name of the desired hosted zone.
* `private_zone` - (Optional) Used with `name` field to get a private hosted zone.
* `vpc_id` - (Optional) Used with `name` field to get a private hosted zone associated with the vpc_id (in this case, private_zone is not mandatory).
* `tags` - (Optional) Used with `name` field. Map of tags, each pair of which must exactly match a pair on the desired hosted zone.

## Attributes Reference

### Supported attributes

All of the argument attributes are also exported as
result attributes. This data source will complete the data by populating
any fields that are not included in the configuration with the data for
the selected hosted zone.

The following attribute is additionally exported:

* `arn` - The Amazon Resource Name (ARN) of the hosted zone.
* `caller_reference` - Caller Reference of the hosted zone.
* `comment` - The comment field of the hosted zone.
* `name_servers` - The list of DNS name servers for the hosted zone.
* `resource_record_set_count` - The number of record sets in the hosted zone.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`linked_service_principal`, `linked_service_description`.
