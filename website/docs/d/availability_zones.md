---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_availability_zones"
description: |-
    Provides a list of availability zones.
---

# Data Source: aws_availability_zones

Provides a list of availability zones.

This is different from the [`aws_availability_zone`](availability_zone.md) (singular) data source,
which provides some details about a specific availability zone.

## Example Usage

### By State

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}
```


## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Configuration block(s) for filtering. Detailed below.
* `state` - (Optional) Allows to filter list of availability zones based on their
current state. Can be either `"available"`, `"information"`, `"impaired"` or
`"unavailable"`.

### filter Configuration Block

The following arguments are supported by the `filter` configuration block:

* `name` - (Required) The name of the filter field.
* `values` - (Required) Set of values that are accepted for the given filter field. Results will be selected if any given value matches.

For more information about filtering, see the [EC2 API documentation][describe-azs].

[describe-azs]: https://docs.k2.cloud/en/api/ec2/placements/DescribeAvailabilityZones.html

## Attributes Reference

### Supported attributes


In addition to all arguments above, the following attributes are exported:

* `id` - Region of the availability zones.
* `names` - A list of the availability zone names available to the account.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`all_availability_zones`, `exclude_names`, `exclude_zone_ids`, `group_names`, `zone_ids`.
