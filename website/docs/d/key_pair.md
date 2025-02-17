---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_key_pair"
description: |-
    Provides information about a specific EC2 Key Pair.
---

# Data Source: aws_key_pair

Provides information about a specific EC2 Key Pair.

## Example Usage

The following example shows how to get a EC2 Key Pair from its name.

```terraform
data "aws_key_pair" "example" {
  key_name = "test"
  filter {
    name   = "tag:Component"
    values = ["web"]
  }
}

output "fingerprint" {
  value = data.aws_key_pair.example.fingerprint
}

output "name" {
  value = data.aws_key_pair.example.key_name
}

output "id" {
  value = data.aws_key_pair.example.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
Key Pairs. The given filters must match exactly one Key Pair
whose data will be exported as attributes.

* `key_pair_id` - (Optional) The Key Pair ID.
* `key_name` - (Optional) The Key Pair name.
* `filter` -  (Optional) Custom filter block as described below.

### filter Configuration Block

The following arguments are supported by the `filter` configuration block:

* `name` - (Required) The name of the filter field.
* `values` - (Required) Set of values that are accepted for the given filter field. Results will be selected if any given value matches.

For more information about filtering, see the [EC2 API documentation][describe-key-pairs].

[describe-key-pairs]: https://docs.k2.cloud/en/api/ec2/key_pairs/DescribeKeyPairs.html

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the Key Pair.
* `arn` - The ARN of the Key Pair.
* `fingerprint` - The SHA-1 digest of the DER encoded private key.
* `tags` - Any tags assigned to the Key Pair.
