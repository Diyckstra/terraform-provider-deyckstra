---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eips"
description: |-
    Provides a list of Elastic IPs in a region.
---

[describe-addresses]: https://docs.cloud.croc.ru/en/api/ec2/addresses/DescribeAddresses.html

# Data Source: aws_eips

Provides a list of Elastic IPs in a region.

## Example Usage

The following shows all Elastic IPs with the specific tag value.

```terraform
data "aws_eips" "example" {
  tags = {
    Env = "dev"
  }
}

output "allocation_ids" {
  value = data.aws_eips.example.allocation_ids
}

output "public_ips" {
  value = data.aws_eips.example.public_ips
}
```

## Argument Reference

* `filter` – (Optional) One or more name/value pairs to use as filters.
  An Elastic IP will be selected if any of the given values match.
    * _Valid values_: See names and values in [EC2 API documentation][describe-addresses].
* `tags` – (Optional) Map of tags, each pair of which must exactly match a pair on the desired Elastic IPs.

## Attributes Reference

* `id` – The region (e.g., `region-1`).
* `allocation_ids` – A list of all the allocation IDs.
* `public_ips` – A list of all the Elastic IP addresses.
