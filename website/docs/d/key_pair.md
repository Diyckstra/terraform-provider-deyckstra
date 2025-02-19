---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_key_pair"
description: |-
    Provides information about a specific EC2 key pair.
---

# Data Source: aws_key_pair

Provides information about a specific EC2 key pair.

## Example Usage

The following example shows how to get a EC2 key pair from its name.

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

The arguments of this data source act as filters for querying the available key pairs.
The given filters must match exactly one key pair
whose data will be exported as attributes.

* `key_pair_id` - (Optional) The key pair ID.
* `key_name` - (Optional) The key pair name.
* `filter` -  (Optional) One or more name/value pairs to use as filters.
  Valid names and values can be found in the [EC2 API documentation][describe-key-pairs].

[describe-key-pairs]: https://docs.k2.cloud/en/api/ec2/key_pairs/DescribeKeyPairs.html

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the key pair.
* `arn` - The ARN of the key pair.
* `fingerprint` - The SHA-1 digest of the DER encoded private key.
* `tags` - Map of tags assigned to the key pair.
