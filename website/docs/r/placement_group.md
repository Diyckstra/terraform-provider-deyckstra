---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_placement_group"
description: |-
  Provides an EC2 placement group.
---

[placement-groups]: https://docs.k2.cloud/en/services/compute/placementgroups.html

# Resource: aws_placement_group

Provides an EC2 placement group.

For more information, see the documentation on [Placement groups][placement-groups].

## Example Usage

```terraform
resource "aws_placement_group" "example" {
  name     = "test-pg"
  strategy = "cluster"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the placement group.
* `strategy` - (Required) The placement strategy. Possible values: `"spread"`.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the placement group.
* `id` - The name of the placement group.
* `placement_group_id` - The ID of the placement group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported: `partition_count`.

## Import

Placement groups can be imported using the `name`, e.g.,

```
$ terraform import aws_placement_group.prod_pg production-placement-group
```

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
