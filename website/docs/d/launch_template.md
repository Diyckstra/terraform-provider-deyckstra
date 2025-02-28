---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_launch_template"
description: |-
  Provides an EC2 launch template data source.
---

[describe-lts]: https://docs.k2.cloud/en/api/ec2/launch_templates/DescribeLaunchTemplates.html

# Data Source: aws_launch_template

Provides information about an EC2 launch template.

## Example Usage

```terraform
data "aws_launch_template" "example" {
  name = "tf-lt"
}
```

### Filter

```terraform
data "aws_launch_template" "example" {
  filter {
    name   = "launch-template-name"
    values = ["some-template"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Configuration block(s) for filtering. Detailed below.
* `id` - (Optional) The ID of the specific launch template to retrieve.
* `name` - (Optional) The name of the launch template.
* `tags` - (Optional) A map of tags, each pair of which must exactly match a pair on the desired launch template.

### filter Configuration Block

The following arguments are supported by the `filter` configuration block:

* `name` - (Required) The name of the filter field.
* `values` - (Required) Set of values that are accepted for the given filter field. Results will be selected if any given value matches.

For more information about filtering, see the [EC2 API documentation][describe-lts].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the launch template.

This resource also exports a full set of attributes corresponding to the arguments of the [`aws_launch_template`](../r/launch_template.md) resource.
