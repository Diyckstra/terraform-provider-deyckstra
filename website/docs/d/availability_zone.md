---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_availability_zone"
description: |-
    Provides details about a specific availability zone
---

# Data Source: aws_availability_zone

`aws_availability_zone` provides details about a specific availability zone (AZ).

This is different from the [`aws_availability_zones`](availability_zones.md) (plural) data source,
which provides a list of the available zones.

## Example Usage

```terraform
data "aws_availability_zone" "example" {
  name = "ru-msk-vol52"
}

output "availability_zone_to_region" {
  value = data.aws_availability_zone.example.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
availability zones. The given filters must match exactly one availability
zone whose data will be exported as attributes.

* `filter` - (Optional) Configuration block(s) for filtering. Detailed below.
* `name` - (Optional) The full name of the availability zone to select.
* `state` - (Optional) A specific availability zone state to require. May be any of `"available"`, `"information"` or `"impaired"`.

### filter Configuration Block

The following arguments are supported by the `filter` configuration block:

* `name` - (Required) The name of the filter field.
* `values` - (Required) Set of values that are accepted for the given filter field. Results will be selected if any given value matches.

For more information about filtering, see the [EC2 API documentation][describe-azs].

[describe-azs]: https://docs.k2.cloud/en/api/ec2/placements/DescribeAvailabilityZones.html

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `region` - The region where the selected availability zone resides.
* `state` - A specific availability zone state to require. Possible values: `"available"`, `"information"`, `"impaired"`, `"unavailable"`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`all_availability_zones`, `group_name`, `network_border_group`, `opt_in_status`, `parent_zone_id`, `parent_zone_name`, `zone_id`, `zone_type`.
