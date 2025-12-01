---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_ec2_traffic_mirror_filter"
description: |-
  Manages a traffic mirror filter.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[traffic-mirroring]: https://docs.k2.cloud/en/services/interconnect/traffic_mirroring.html

# Resource: aws_ec2_traffic_mirror_filter

Manages a traffic mirror filter. For details about traffic mirroring, see the [user documentation][traffic-mirroring].

## Example Usage

To create a basic traffic mirror filter, use:

```terraform
resource "aws_ec2_traffic_mirror_filter" "foo" {
  description = "traffic mirror filter - terraform example"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, Forces new resource) Description of the filter.
* `tags` - (Optional, Editable) Map of tags to assign to the traffic mirror filter. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the traffic mirror filter.
* `id` - The ID of the traffic mirror filter.
* `tags_all` - Map of tags assigned to the traffic mirror filter, including those inherited from the provider [`default_tags` configuration block][default-tags].

## Import

In Terraform v1.5.0 or later, traffic mirror filter can be imported by `id` using the `import` block.

```terraform
import {
  to = aws_ec2_traffic_mirror_filter.foo
  id = "tmf-12345678"
}
```

In older Terraform versions, the traffic mirror filter can be imported by its `id` using `terraform import`, e.g.:

```console
% terraform import aws_ec2_traffic_mirror_filter.foo tmf-12345678
```
