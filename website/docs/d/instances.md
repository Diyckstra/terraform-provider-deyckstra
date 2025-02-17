---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_instances"
description: |-
  Provides information about EC2 instances.
---

[describe-instances]: https://docs.k2.cloud/en/api/ec2/instances/DescribeInstances.html
[outputs]: https://www.terraform.io/docs/configuration/outputs.html
[remote state]: https://www.terraform.io/docs/state/remote.html
[terraform_remote_state]: https://www.terraform.io/docs/providers/terraform/d/remote_state.html

# Data Source: aws_instances

Provides information about EC2 instances. This data source can be used to get IDs or IPs of EC2 instances to be referenced elsewhere.

-> **Note:** It's a best practice to expose instance details via [outputs], and [remote state],
and **use [`terraform_remote_state`][terraform_remote_state] data source instead** if you manage referenced instances via Terraform.

~> **Note** It's strongly discouraged to use this data source for querying ephemeral
instances (e.g., managed via autoscaling group), as the output may change at any time
and you'd need to re-run `apply` every time an instance comes up or dies.

## Example Usage

```terraform
data "aws_instances" "selected" {
  instance_tags = {
    type = "test"
  }

  filter {
    name   = "instance.group-id"
    values = ["sg-12345678"]
  }

  instance_state_names = ["running", "stopped"]
}

resource "aws_eip" "example" {
  count    = length(data.aws_instances.selected.ids)
  instance = data.aws_instances.selected.ids[count.index]
}
```

## Argument Reference

* `instance_tags` - (Optional) A map of tags, each pair of which must exactly match a pair on desired instances.
* `instance_state_names` - (Optional) A list of instance states that should be applicable to the desired instances.
  Valid values are `pending, running, shutting-down, stopped, stopping, terminated`.
* `filter` - (Optional) One or more name/value pairs to use as filters.

For more information about filtering, see the [EC2 API documentation][describe-instances].

## Attributes Reference

* `id` - The region.
* `ids` - IDs of instances found through the filter.
* `private_ips` - Private IP addresses of instances found through the filter.
* `public_ips` - Public IP addresses of instances found through the filter.
